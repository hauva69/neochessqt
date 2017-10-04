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
	collectionbuckets := [][]byte{[]byte("NEObucket"), []byte("PGNbucket"), []byte("SCIDbucket")}
	for index := 0; index < 3; index++ {
		cat := categories[index]
		bucket := collectionbuckets[index]
		catnode := gui.NewQStandardItem2(cat)
		switch {
		case cat == "PGN DataBases":
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

// LoadCatalogDatabase comment
func (dbd *DBDock) LoadCatalogDatabase(event *gui.QMouseEvent) {
	dbd.tree.MouseDoubleClickEventDefault(event)
	// log.Info(fmt.Sprintf("Loadinging type: %s with key: %s ", dbd.tree.currenttype, dbd.tree.currentfile))
	cdb, err := OpenFile(d dbd.tree.currentfile, dbd.tree.currenttype)
	if cdb.CheckIndex {
		needsupdate, err := cdb.NeedIndex()
		if err == nil && needsupdate {
			log.Info("Re-indexing database")
			dialog := widgets.NewQProgressDialog2("Re-indexing PGN File", "Cancel", 0, 100, nil, core.Qt__Dialog)
			dialog.SetWindowModality(core.Qt__WindowModal)
			_, err := cdb.Index(dialog)
			if err != nil {
				log.Error(err)
			}
		}
	}
	liniterror := cdb.LoadInitialGamesList()
	if liniterror != nil {
		log.Error(err)
	}
	ldbproperror := cdb.LoadDBProperties()
	if ldbproperror != nil {
		log.Error(err)
	}
}
