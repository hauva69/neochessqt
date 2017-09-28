package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// DBDock struct
type DBDock struct {
	widgets.QDockWidget
	tree  *widgets.QTreeView
	model *gui.QStandardItemModel
}

func initDBDock(w *widgets.QMainWindow) *DBDock {
	var this = NewDBDock("Databases", w, core.Qt__Widget)
	var dbdock = this
	dbdock.SetAllowedAreas(core.Qt__LeftDockWidgetArea | core.Qt__RightDockWidgetArea)
	dbdock.tree = widgets.NewQTreeView(this)
	dbdock.model = gui.NewQStandardItemModel(nil)
	dbdock.model.SetHorizontalHeaderLabels([]string{"Recent Databases"})
	dbdock.tree.SetModel(dbdock.model)
	dbdock.tree.SetUniformRowHeights(true)
	dbdock.tree.SetAlternatingRowColors(true)

	categories := []string{"NeoChess Databases", "PGN DataBases", "SCID Databases"}
	// collectionbuckets := [][]byte{[]byte("NEObucket"), []byte("PGNbucket"), []byte("SCIDbucket")}
	for index := 0; index < 3; index++ {
		cat := categories[index]
		// bucket := collectionbuckets[index]
		catnode := gui.NewQStandardItem2(cat)
		switch {
		case cat == "PGN DataBases":
		case cat == "NeoChess Databases":
		case cat == "SCID Databases":
		}
		dbdock.model.AppendRow2(catnode)
	}
	dbdock.SetWidget(dbdock.tree)
	dbdock.tree.ExpandAll()

	return this
}
