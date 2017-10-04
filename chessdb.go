package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/widgets"
)

// ChessDataBase for holding reference to a loaded ChessDatabase file
type ChessDataBase struct {
	Displayname string    `json:"displayname"`
	Fullpath    string    `json:"fullpath"`
	Key         string    `json:"key"`
	Basename    string    `json:"basename"`
	Filemod     time.Time `json:"filemod"`
	Filesize    int64     `json:"filesize"`
	Count       int       `json:"count"`
	Gameoffsets []int64   `json:"gameoffsets"`
	Gamelengths []int     `json:"gamelengths"`
	Notes       string    `json:"notes"`
	InitIndex   bool      `json:"initindex"`
	CheckIndex  bool      `json:"checkindex"`
	Kind        string    `json:"kind"`
}

// NewCDB Create instance
func NewCDB() *ChessDataBase {
	return &ChessDataBase{}
}

// LoadDBProperties dockwindow
func (cdb *ChessDataBase) LoadDBProperties() error {
	// dbtree.SetNameProp(cdb.Displayname)
	// dbtree.SetFileProp(cdb.Fullpath)
	// dbtree.SetCountProp(strconv.Itoa(cdb.Count))
	// dbtree.SetNotesProp(cdb.Notes)
	// dbtree.SetDateModifiedProp(cdb.Filemod)
	return nil
}

// OpenFile of Chess Database
func OpenFile(filename string, kind string) (*ChessDataBase, error) {
	var cdb *ChessDataBase
	bucket := kind + "bucket"
	err := catdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			err := errors.New("bucket doesn't exist")
			return err
		}

		k := []byte(filename)
		v := b.Get(k)
		if len(v) > 0 {
			err := json.Unmarshal(v, &cdb)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if cdb == nil {
		fhandle, fileerr := os.Open(filename)
		if fileerr != nil {
			return nil, fileerr
		}
		defer fhandle.Close()

		fileinfo, staterr := fhandle.Stat()
		if staterr != nil {
			return nil, staterr
		}
		cdb = NewCDB()
		cdb.Displayname = fileinfo.Name()
		cdb.Basename = fileinfo.Name()
		cdb.Fullpath = filename
		cdb.Kind = kind
		switch {
		case cdb.Kind == "PGN":
			cdb.Key = filename
		default:
			cdb.Key = filename
		}
		cdb.Filesize = fileinfo.Size()
		cdb.Filemod = fileinfo.ModTime()
		cdb.InitIndex = true
		cdb.Save()
	} else {
		log.Info("Have index.")
		cdb.CheckIndex = true
	}
	return cdb, err
}

// LoadInitialGamesList of ChessDatabase
func (cdb *ChessDataBase) LoadInitialGamesList() error {
	log.Info("Loading initial games list")
	glist := make([][]string, 10)
	for i := 0; i < 10; i++ {
		glist[i] = make([]string, 10)
	}
	for gameindex := 1; gameindex < 10; gameindex++ {
		gamebytes, err := cdb.GetGame(gameindex)
		if err != nil {
			log.Error(err)
		}
		game := ParseGameString(gamebytes, 1, true)
		glist[gameindex][0] = strconv.Itoa(gameindex)
		glist[gameindex][1] = game.GameHeader.Event
		glist[gameindex][2] = game.GameHeader.Site
		glist[gameindex][3] = game.GameHeader.Date
		glist[gameindex][4] = game.GameHeader.Round
		glist[gameindex][5] = game.GameHeader.White
		glist[gameindex][6] = game.GameHeader.Black
		glist[gameindex][7] = game.GameHeader.Result
		glist[gameindex][8] = game.GameHeader.ECO
		glist[gameindex][9] = game.GameHeader.Opening
	}
	//gamelistwidget.SetRows(glist)
	// gamelistwidget.currentdb = cdb
	return nil
}

// Save instance of ChessDatabase
func (cdb *ChessDataBase) Save() error {
	bucket := cdb.Kind + "bucket"
	err := catdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(cdb)
		if err != nil {
			return err
		}
		return b.Put([]byte(cdb.Key), encoded)
	})
	return err
}

// NeedIndex check database for need to refresh
func (cdb *ChessDataBase) NeedIndex() (bool, error) {
	log.Info("Checking for need to re-index")
	fhandle, fileerr := os.Open(cdb.Fullpath)
	if fileerr != nil {
		return false, fileerr
	}
	defer fhandle.Close()

	fileinfo, staterr := fhandle.Stat()
	if staterr != nil {
		return false, staterr
	}
	if cdb.Filemod != fileinfo.ModTime() || cdb.Filesize != fileinfo.Size() {
		return true, nil
	}
	return false, nil
}

// GetGame from database
func (cdb *ChessDataBase) GetGame(gamenumber int) ([]byte, error) {
	if gamenumber < 0 || gamenumber > cdb.Count {
		err := errors.New("game number request not in rage of games in database")
		return nil, err
	}
	fhandle, fileerr := os.Open(cdb.Fullpath)
	if fileerr != nil {
		return nil, fileerr
	}
	defer fhandle.Close()
	offset := int64(cdb.Gameoffsets[gamenumber-1])
	length := cdb.Gamelengths[gamenumber-1]
	buf := make([]byte, length)

	fhandle.Seek(offset, 0)
	_, err := fhandle.Read(buf)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// ReadThroughPGN of Game
func (cdb *ChessDataBase) ReadThroughPGN(input io.ReadSeeker, start int64) (pos int64, gamelength int, err error) {
	ingame := false
	inmoves := false
	if _, err := input.Seek(start, 0); err != nil {
		return 0, 0, err
	}
	reader := bufio.NewReader(input)
	pos = start
	gamelength = 0
	for {
		data, err := reader.ReadBytes('\n')
		// data, _, err := reader.ReadLine()
		gamelength += len(data)
		pos += int64(len(data))
		if err != nil {
			if err != io.EOF {
				fmt.Printf("IO Read Error")
			}
			break
		}
		line := string(data)
		trimline := strings.TrimSpace(string(line))
		if !strings.HasPrefix(trimline, "%") {
			if !ingame {
				if strings.HasPrefix(trimline, "[") {
					ingame = true
					inmoves = false
					//	gamelines = append(gamelines, string(line))
				}
			} else {
				if !inmoves {
					if trimline == "" {
						inmoves = true
					}
					//	gamelines = append(gamelines, string(line))
				} else {
					//	gamelines = append(gamelines, string(line))
					if trimline == "" {
						break
					}
				}
			}
		}
	}
	return pos, gamelength, err
}

// Index ChessDataBase or reindex
func (cdb *ChessDataBase) Index(progress *widgets.QProgressDialog) (int, error) {
	fhandle, fileerr := os.Open(cdb.Fullpath)
	if fileerr != nil {
		return 0, fileerr
	}
	defer fhandle.Close()

	progress.SetRange(0, int(cdb.Filesize/1000))
	var offset int64
	var count int
	for {
		pos, gamelength, err := cdb.ReadThroughPGN(fhandle, offset)
		progress.SetValue(int(offset / 1000))
		cdb.Gameoffsets = append(cdb.Gameoffsets, offset)
		newgamelength := int(pos - offset)
		cdb.Gamelengths = append(cdb.Gamelengths, newgamelength)
		offset = pos
		if err != nil || gamelength == 0 {
			break
		}
		count++
	}
	cdb.Count = count
	err := cdb.Save()
	if err != nil {
		return count, err
	}
	cdb.InitIndex = false
	cdb.CheckIndex = false
	cdb.Save()
	cdb.LoadInitialGamesList()
	return count, nil
}
