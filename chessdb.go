package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rashwell/neochesslib"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/widgets"
)

// NeedIndex check database for need to refresh
func (cdb *neochesslib.ChessDataBase) NeedIndex() (bool, error) {
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
func (cdb *neochesslib.ChessDataBase) GetGame(gamenumber int) ([]byte, error) {
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
func (cdb *neochesslib.ChessDataBase) ReadThroughPGN(input io.ReadSeeker, start int64) (pos int64, gamelength int, err error) {
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
func (cdb *neochesslib.ChessDataBase) Index(progress *widgets.QProgressDialog) (int, error) {
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
	return count, nil
}
