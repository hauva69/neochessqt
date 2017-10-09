package main

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// GameListDock comment
type GameListDock struct {
	widgets.QDockWidget
	table *widgets.QTableWidget
}

// NewGameList Initialize
func initGameListDock(w *widgets.QMainWindow, cdbv *ChessDBView) *GameListDock {
	this := NewGameListDock("Game List", w, core.Qt__Widget)
	var gamelist = this
	gamestable := widgets.NewQTableWidget2(1, 10, nil)
	gamestable.SetHorizontalHeaderLabels([]string{"Game #", "Event", "Site", "Date", "Round", "White", "Black", "Result", "ECO", "Opening"})
	gamestable.SetVerticalHeaderLabels([]string{"", "", "", "", "", ""})
	gamestable.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	view := widgets.NewQTableViewFromPointer(widgets.PointerFromQTableWidget(gamestable))
	view.VerticalHeader().Hide()
	view.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
	gamelist.SetAllowedAreas(core.Qt__BottomDockWidgetArea)
	gamelist.SetWidget(view)
	gamelist.table = gamestable
	gamelist.table.ConnectMouseDoubleClickEvent(cdbv.LoadSelectedGame)
	return this
}

// SelectedGame from game list
func (g *GameListDock) SelectedGame() (int, error) {
	items := g.table.SelectedItems()
	if len(items) > 0 {
		item := g.table.Item(items[0].Row(), 0)
		log.Info(fmt.Sprintf("Selected Game with id: %s", item.Text()))
		index, err := strconv.Atoi(item.Text())
		if err != nil {
			return 0, err
		}
		return index, nil
	}
	return 0, errors.New("nothing selected")
}

// AddRow to Game list
func (g *GameListDock) AddRow(contentrow []string) {
	rowcount := g.table.RowCount()
	rowcount = rowcount + 1
	g.table.InsertRow(rowcount)
	for fieldindex, field := range contentrow {
		item := widgets.NewQTableWidgetItem2(field, fieldindex)
		item.SetFlags(item.Flags() ^ core.Qt__ItemIsEditable)
		g.table.SetItem(rowcount, fieldindex, item)
	}
}

// SetRows of Gamelist
func (g *GameListDock) SetRows(content [][]string) {
	for rowindex, row := range content {
		g.table.InsertRow(rowindex - 1)
		for fieldindex, field := range row {
			item := widgets.NewQTableWidgetItem2(field, fieldindex)
			item.SetFlags(item.Flags() ^ core.Qt__ItemIsEditable)
			g.table.SetItem(rowindex-1, fieldindex, item)
		}
	}
}
