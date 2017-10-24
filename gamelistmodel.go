package main

import (
	"strconv"

	"github.com/rashwell/neochesslib"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
)

type GameListModel struct {
	core.QAbstractTableModel
	cdb       *ChessDataBase
	glHeaders []string
	gamecount int
	buffgames map[int]*neochesslib.Game

	_ func()    `constructor:"init"`
	_ func(int) `slot:"numberPopulated"`
}

func (gm *GameListModel) init() {
	gm.glHeaders = []string{"Game #", "Event", "Site", "Date", "Round #", "White", "Black", "Result", "ECO", "Opening"}
	gm.buffgames = make(map[int]*neochesslib.Game)
	gm.ConnectHeaderData(gm.headerdata)
	gm.ConnectNumberPopulated(gm.populated)
	gm.ConnectRowCount(func(p *core.QModelIndex) int {
		return gm.gamecount
	})
	gm.ConnectColumnCount(func(p *core.QModelIndex) int {
		return 10
	})
	gm.ConnectData(gm.data)

	gm.ConnectFetchMore(func(p *core.QModelIndex) {
		remainder := gm.cdb.Count - gm.gamecount
		gamestofetch := 0
		if 100 < remainder {
			gamestofetch = 10
		} else {
			gamestofetch = remainder
		}
		gm.BeginInsertRows(p, gm.gamecount, gm.gamecount+gamestofetch-1)
		gm.gamecount += gamestofetch
		gm.EndInsertRows()
		gm.populated(gamestofetch)
	})
	gm.ConnectCanFetchMore(func(p *core.QModelIndex) bool {
		if gm.cdb == nil {
			return false
		}
		if gm.gamecount < gm.cdb.Count {
			return true
		}
		return false
	})
	gm.ConnectNumberPopulated(gm.populated)
}

func (gm *GameListModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}
	if gm.cdb == nil {
		return core.NewQVariant()
	}
	if index.Row() >= gm.cdb.Count || index.Row() < 0 {
		return core.NewQVariant()
	}
	if role == int(core.Qt__UserRole) {
		_, exists := gm.buffgames[index.Row()+1]
		if !exists {
			log.Infof("Model grabbing game: %d", index.Row()+1)
			pgn, err := gm.cdb.GetGame(index.Row() + 1)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Game size: %d", len(pgn))
			gm.buffgames[index.Row()+1] = ParseGameString([]byte(pgn), index.Row()+1, false)
		}
		game := gm.buffgames[index.Row()+1]
		return core.NewQVariant17(strconv.FormatUint(uint64(game.ID), 10))
	}
	if role == int(core.Qt__DisplayRole) {
		_, exists := gm.buffgames[index.Row()+1]
		if !exists {
			log.Infof("Model grabbing game: %d", index.Row()+1)
			pgn, err := gm.cdb.GetGame(index.Row() + 1)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Game size: %d", len(pgn))
			gm.buffgames[index.Row()+1] = ParseGameString([]byte(pgn), index.Row()+1, false)
		}
		game := gm.buffgames[index.Row()+1]
		switch index.Column() {
		case 0:
			return core.NewQVariant17(strconv.FormatUint(uint64(game.ID), 10))
		case 1:
			return core.NewQVariant17(game.Event)
		case 2:
			return core.NewQVariant17(game.Site)
		case 3:
			return core.NewQVariant17(game.Date)
		case 4:
			return core.NewQVariant17(game.Round)
		case 5:
			return core.NewQVariant17(game.White)
		case 6:
			return core.NewQVariant17(game.Black)
		case 7:
			return core.NewQVariant17(game.Result)
		case 8:
			return core.NewQVariant17(game.ECO)
		case 9:
			return core.NewQVariant17(game.Opening)
		}
	}
	return core.NewQVariant()
}

func (gm *GameListModel) headerdata(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if orientation == core.Qt__Horizontal && role == int(core.Qt__DisplayRole) {
		if section >= 0 && section <= 9 {
			return core.NewQVariant17(gm.glHeaders[section])
		}
	}
	if orientation == core.Qt__Vertical && role == int(core.Qt__DisplayRole) {
		if section >= 0 && section <= 9 {
			return core.NewQVariant17(gm.glHeaders[section])
		}
	}
	return core.NewQVariant()
}

func (gm *GameListModel) populated(popcount int) {
	log.Infof("Recieved Signal Populated: %d", popcount)
}
