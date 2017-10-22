package main

import (
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// DBDock struct
type DBDock struct {
	widgets.QDockWidget
	tree         *widgets.QTreeView
	model        *gui.QStandardItemModel
	current      string
	currenttype  string
	currentfile  string
	currentcount string
	currentview  *ChessDBView
	pgncat       *gui.QStandardItem
}

func initDBDock(w *widgets.QMainWindow, cdbview *ChessDBView) *DBDock {
	var this = NewDBDock("Databases", w, core.Qt__Widget)
	var dbdock = this
	dbdock.currentview = cdbview
	dbdock.SetAllowedAreas(core.Qt__LeftDockWidgetArea | core.Qt__RightDockWidgetArea)
	dbdock.tree = widgets.NewQTreeView(this)
	dbdock.model = gui.NewQStandardItemModel(nil)
	dbdock.model.SetHorizontalHeaderLabels([]string{"Recent Databases"})
	dbdock.tree.SetModel(dbdock.model)
	dbdock.tree.SetUniformRowHeights(true)
	dbdock.tree.SetAlternatingRowColors(true)
	dbdock.tree.ConnectMouseDoubleClickEvent(dbdock.LoadCatalogDatabase)
	dbdock.tree.ConnectSelectionCommand(dbdock.SelectionCommand)
	categories := []string{"NeoChess Databases", "PGN DataBases", "SCID Databases"}
	collectionbuckets := [][]byte{[]byte("NEObucket"), []byte("PGNbucket"), []byte("SCIDbucket")}
	for index := 0; index < 3; index++ {
		cat := categories[index]
		bucket := collectionbuckets[index]
		catnode := gui.NewQStandardItem2(cat)
		switch {
		case cat == "PGN DataBases":
			dbdock.pgncat = catnode
			catdb.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(bucket)
				if b != nil {
					c := b.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						var cdb *ChessDataBase
						err := json.Unmarshal(v, &cdb)
						if err != nil {
							log.Error(err)
						}
						child := gui.NewQStandardItem()
						child.SetEditable(false)
						displaystr := core.NewQVariant17(cdb.Displayname)
						child.SetData(displaystr, int(core.Qt__DisplayRole))
						keystr := core.NewQVariant17(cdb.Fullpath)
						child.SetData(keystr, int(core.Qt__UserRole))
						typestr := core.NewQVariant17(cdb.Kind)
						child.SetData(typestr, int(core.Qt__UserRole)+1)
						countstr := core.NewQVariant17(strconv.Itoa(cdb.Count))
						child.SetData(countstr, int(core.Qt__UserRole)+2)
						catnode.AppendRow2(child)
					}
				}
				return nil
			})
		case cat == "NeoChess Databases":
		case cat == "SCID Databases":
		}
		dbdock.model.AppendRow2(catnode)
	}
	dbdock.SetWidget(dbdock.tree)
	dbdock.tree.ExpandAll()

	return this
}

func (dbd *DBDock) InsertPGNDB(displayname string, fp string, count int) {
	child := gui.NewQStandardItem()
	child.SetEditable(false)
	displaystr := core.NewQVariant17(displayname)
	child.SetData(displaystr, int(core.Qt__DisplayRole))
	keystr := core.NewQVariant17(fp)
	child.SetData(keystr, int(core.Qt__UserRole))
	typestr := core.NewQVariant17("PGN")
	child.SetData(typestr, int(core.Qt__UserRole)+1)
	countstr := core.NewQVariant17(strconv.Itoa(count))
	child.SetData(countstr, int(core.Qt__UserRole)+2)
	dbd.pgncat.AppendRow2(child)
}

func (dbd *DBDock) SelectionCommand(index *core.QModelIndex, event *core.QEvent) core.QItemSelectionModel__SelectionFlag {
	dbd.current = dbd.model.Data(index, int(core.Qt__DisplayRole)).ToString()
	dbd.currentfile = dbd.model.Data(index, int(core.Qt__UserRole)).ToString()
	dbd.currenttype = dbd.model.Data(index, int(core.Qt__UserRole)+1).ToString()
	dbd.currentcount = dbd.model.Data(index, int(core.Qt__UserRole)+2).ToString()
	// Log.Info(fmt.Sprintf("Current Type: %s with Index: %s", dbt.currenttype, dbt.current))
	return dbd.tree.SelectionCommandDefault(index, event)
}

// LoadCatalogDatabase comment
func (dbd *DBDock) LoadCatalogDatabase(event *gui.QMouseEvent) {
	dbd.tree.MouseDoubleClickEventDefault(event)
	dbd.currentview.setDatabase(dbd.currentfile, dbd.currenttype)
	/*
		liniterror := dbd.currentview.LoadInitialGamesList()
		if liniterror != nil {
			log.Error(err)
		}
	*/
	/*
		ldbproperror := dbd.currentview.LoadDBProperties()
		if ldbproperror != nil {
			log.Error(err)
		}
	*/
}
