package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// GameListDock comment
type GameListDock struct {
	widgets.QDockWidget
	table *widgets.QTableWidget
	// form  *widgets.QWidget
	// currentdb *ChessDataBase
}

// NewGameList Initialize
func initGameListDock(w *widgets.QMainWindow) *GameListDock {
	this := NewGameListDock("Game List", w, core.Qt__Widget)
	var gamelist = this
	gamestable := widgets.NewQTableWidget2(1, 10, nil)
	gamestable.SetHorizontalHeaderLabels([]string{"Game #", "Event", "Site", "Date", "Round", "White", "Black", "Result", "ECO", "Opening"})
	gamestable.SetVerticalHeaderLabels([]string{"", "", "", "", "", ""})
	gamestable.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	gamestableview := widgets.NewQTableViewFromPointer(widgets.PointerFromQTableWidget(gamestable))
	gamestableview.VerticalHeader().Hide()
	gamestableview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
	gamelist.SetAllowedAreas(core.Qt__BottomDockWidgetArea)
	gamelist.SetWidget(gamestableview)
	gamelist.table = gamestable
	// gamelist.table.ConnectMouseDoubleClickEvent(gamelist.LoadSelectedGame)
	return this
}

/*
// LoadSelectedGame from gamelist
func (g *GameList) LoadSelectedGame(event *gui.QMouseEvent) {
	items := g.table.SelectedItems()
	if len(items) > 0 {
		item := g.table.Item(items[0].Row(), 0)
		Log.Info(fmt.Sprintf("Loading game with id: %s", item.Text()))
		index, err := strconv.Atoi(item.Text())
		if err != nil {
			return
		}
		pgn, err := g.currentdb.GetGame(index)
		if err != nil {
			return
		}
		boardview.game = ParseGameString([]byte(pgn), index, false)
		gamedetail.Load(boardview.game)
		boardview.game.LoadMoves(boardview.board)
		for _, v := range boardview.squarelist {
			boardview.gscene.RemoveItem(v.qpixitem)
		}
		boardview.SetInitBoard(boardview.game, boardview.board)
		pgneditorwidget.SetPGN(boardview.game.CurrentPgn)
		Log.Info(string(pgn))
	}
	g.table.MouseDoubleClickEventDefault(event)
}

// SetRows of Gamelist
func (g *GameList) SetRows(content [][]string) {
	for rowindex, row := range content {
		g.table.InsertRow(rowindex - 1)
		for fieldindex, field := range row {
			item := widgets.NewQTableWidgetItem2(field, fieldindex)
			item.SetFlags(item.Flags() ^ core.Qt__ItemIsEditable)
			g.table.SetItem(rowindex-1, fieldindex, item)
		}
	}
}
*/
