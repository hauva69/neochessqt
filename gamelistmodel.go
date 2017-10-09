package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
)

type GameListModel struct {
	core.QAbstractTableModel
	cdb       *ChessDataBase
	gamecount int

	_ func()               `constructor:"init"`
	_ func(int)            `signal:"numberPopulated"`
	_ func(string, string) `slot:"setDatabase"`
}

func (gm *GameListModel) init() {
	gm.ConnectNumberPopulated(gm.numberPopulated)

	gm.ConnectSetDatabase(gm.setDatabase)

	gm.ConnectRowCount(func(p *core.QModelIndex) int {
		return gm.gamecount
	})

	gm.ConnectColumnCount(func(p *core.QModelIndex) int {
		return 10
	})

	gm.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		if !index.IsValid() {
			return core.NewQVariant()
		}
		if index.Row() >= gm.cdb.Count || index.Row() < 0 {
			return core.NewQVariant()
		}
		if role == int(core.Qt__DisplayRole) {
			pgn, _ := gm.cdb.GetGame(index.Row())
			game, _ := ParseGameString([]byte(pgn), index.Row(), false)
			switch index.Column() {
			case 0:
				return core.NewQVariant17(strconv.Itoa(index.Row()))
			case 1:
				return core.NewQVariant17(game.GameHeader.Event)
			case 2:
				return core.NewQVariant17(game.GameHeader.Site)
			case 3:
				return core.NewQVariant17(game.GameHeader.Date)
			case 4:
				return core.NewQVariant17(game.GameHeader.Round)
			case 5:
				return core.NewQVariant17(game.GameHeader.White)
			case 6:
				return core.NewQVariant17(game.GameHeader.Black)
			case 7:
				return core.NewQVariant17(game.GameHeader.Result)
			case 8:
				return core.NewQVariant17(game.GameHeader.ECO)
			case 9:
				return core.NewQVariant17(game.GameHeader.Opening)
			}
		}
		return core.NewQVariant()
	})

	gm.ConnectFetchMore(func(p *core.QModelIndex) {
		remainder := gm.cdb.Count - gm.gamecount
		gamestofetch := 0
		if 100 < remainder {
			gamestofetch = 100
		} else {
			gamestofetch = remainder
		}
		gm.BeginInsertRows(p, gm.gamecount, gm.gamecount+gamestofetch-1)
		gm.gamecount += gamestofetch
		gm.EndInsertRows()
		gm.numberPopulated(gamestofetch)
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
}

func (gm *GameListModel) setDatabase(dbpath string, kind string) {
	log.Infof("Opening %s", dbpath)
	gm.cdb, _ = OpenFile(dbpath, kind)
}

func (gm *GameListModel) numberPopulated(popcount int) {
	log.Infof("Recieved Signal Populated: %d", popcount)
}
