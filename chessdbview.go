package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rashwell/neochesslib"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// ChessDBView A view of an Open Chess Database
type ChessDBView struct {
	mainw            *widgets.QMainWindow
	cdb              *neochesslib.ChessDataBase
	model            *GameListModel
	dbdock           *DBDock
	dbdetaildock     *DBDetailDock
	pgndock          *PGNDock
	gamelistdock     *GameListDock
	gamedetaildock   *GameDetailDock
	gameanalysisdock *GameAnalysisDock
	boardview        *BoardView
	currentgame      *neochesslib.Game
}

// initCDBView Create instance
func initCDBView(w *widgets.QMainWindow) *ChessDBView {
	view := &ChessDBView{}
	view.mainw = w
	view.cdb = nil

	view.pgndock = initPGNDock(w, view)
	w.AddDockWidget(core.Qt__RightDockWidgetArea, view.pgndock)

	view.dbdock = initDBDock(w, view)
	w.AddDockWidget(core.Qt__LeftDockWidgetArea, view.dbdock)

	view.dbdetaildock = initDBDetailDock(w)
	w.AddDockWidget(core.Qt__LeftDockWidgetArea, view.dbdetaildock)

	view.gamelistdock = initGameListDock(w, view)
	view.gamelistdock.tableview.ConnectDoubleClicked(view.gameselected)
	w.AddDockWidget(core.Qt__BottomDockWidgetArea, view.gamelistdock)

	view.gamedetaildock = initGameDetailDock(w)
	w.AddDockWidget(core.Qt__RightDockWidgetArea, view.gamedetaildock)
	w.TabifyDockWidget(view.pgndock, view.gamedetaildock)
	w.SetTabPosition(core.Qt__RightDockWidgetArea, widgets.QTabWidget__North)

	view.gameanalysisdock = initGameAnalysisDock(w, view)
	w.AddDockWidget(core.Qt__RightDockWidgetArea, view.gameanalysisdock)
	view.gameanalysisdock.Hide()

	pgntitlewidget := widgets.NewQWidget(w, core.Qt__Widget)
	gamedetailtitlewidget := widgets.NewQWidget(w, core.Qt__Widget)
	view.gamedetaildock.SetTitleBarWidget(gamedetailtitlewidget)
	view.pgndock.SetTitleBarWidget(pgntitlewidget)
	view.pgndock.Raise()

	view.AddGame()

	return view
}

// AddGame add initial game to view
func (cdbv *ChessDBView) AddGame() {
	cdbv.currentgame = neochesslib.NewGame()
	cdbv.boardview = initBoardView(cdbv)
	cdbv.mainw.SetCentralWidget(cdbv.boardview)
	cdbv.gamedetaildock.SetGameTags(cdbv.currentgame)
	cdbv.pgndock.SetPGN(cdbv.currentgame)
	/*
		var gamerow []string
		gamerow = make([]string, 10)
		gamerow[0] = "-1"
		gamerow[1] = cdbv.currentgame.Event
		gamerow[2] = cdbv.currentgame.Site
		gamerow[3] = cdbv.currentgame.Date
		gamerow[4] = cdbv.currentgame.Round
		gamerow[5] = cdbv.currentgame.White
		gamerow[6] = cdbv.currentgame.Black
		gamerow[7] = cdbv.currentgame.Result
		gamerow[8] = cdbv.currentgame.ECO
		gamerow[9] = cdbv.currentgame.Opening
		cdbv.gamelistdock.AddRow(gamerow)
	*/
}

// LoadSelectedGame from gamelist
func (cdbv *ChessDBView) gameselected(index *core.QModelIndex) {
	log.Infof("Game Row Selected: %d", index.Row())
	qvgameid := cdbv.model.data(index, int(core.Qt__UserRole))
	gameid, _ := strconv.Atoi(qvgameid.ToString())
	log.Infof("Loading game id: %d", gameid)
	cdbv.LoadGame(gameid)
	/*	item, err := cdbv.gamelistdock.SelectedGame()
		if err != nil {
			cdbv.gamelistdock.table.MouseDoubleClickEventDefault(event)
		}
		log.Infof("Loading game with id: %d", item)
		cdbv.LoadGame(item)
		cdbv.gamelistdock.table.MouseDoubleClickEventDefault(event)
	*/
}

// LoadGame into view
func (cdbv *ChessDBView) LoadGame(index int) {
	log.Infof("Loading game %d", index)
	pgn, err := cdbv.cdb.GetGame(index)
	if err != nil {
		return
	}
	cdbv.currentgame = ParseGameString([]byte(pgn), index, false)
	// cdbv.currentgame.LoadMoves(cdbv.currentboard)
	cdbv.boardview.SetGame()
	cdbv.gamedetaildock.SetGameTags(cdbv.currentgame)
	cdbv.pgndock.SetPGN(cdbv.currentgame)
}

// SetPosition of View
func (cdbv *ChessDBView) SetPosition(movepos int) {
	/*
		cdbv.currentboard.Reset()
		if cdbv.currentgame.FEN != "" {
			cdbv.currentboard.InitFromFen(cdbv.currentgame.FEN)
		} else {
			cdbv.currentboard.InitFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		}
		// g.LoadMoves(cb)
		for mn := 0; mn < movepos; mn++ {
			cdbv.currentboard.MakeMove(cdbv.currentgame.Moves[mn], true)
		}
		cdbv.boardview.SetGame()
		cdbv.pgndock.SetPGN(cdbv.currentgame)
	*/
}

func (cdbv *ChessDBView) loadpgndb(w *widgets.QMainWindow) {
	docdir := core.QStandardPaths_StandardLocations(core.QStandardPaths__DocumentsLocation)[0]
	filter := "PGN Files (*.pgn);;SCID Database (*.sg4);;NeoChess Database (*.ncb)"
	fileDialog := widgets.NewQFileDialog2(w, "Open Database", docdir, filter)
	fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
	fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
	if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
		return
	}
	filename := fileDialog.SelectedFiles()[0]
	log.Info(filename)
	extension := filepath.Ext(filename)
	dbkind := ""
	if strings.ToUpper(extension) == ".PGN" {
		dbkind = "PGN"
	}
	var err error
	cdbv.cdb, err = OpenFile(filename, dbkind)
	if err != nil {
		log.Error(err)
	}
	if cdbv.cdb.InitIndex {
		dialog := widgets.NewQProgressDialog2("Indexing PGN File", "Cancel", 0, 100, nil, core.Qt__Dialog)
		dialog.SetWindowModality(core.Qt__WindowModal)
		_, err := cdbv.cdb.Index(dialog)
		if err != nil {
			log.Error(err)
		}
	} else {
		if cdbv.cdb.CheckIndex {
			needsupdate, err := cdbv.cdb.NeedIndex()
			if err == nil && needsupdate {
				dialog := widgets.NewQProgressDialog2("Re-indexing PGN File", "Cancel", 0, 100, nil, core.Qt__Dialog)
				dialog.SetWindowModality(core.Qt__WindowModal)
				_, err := cdbv.cdb.Index(dialog)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	/*
		liniterror := cdbv.LoadInitialGamesList()
		if liniterror != nil {
			log.Error(liniterror)
		}
	*/

	ldbproperror := cdbv.LoadDBProperties()
	if ldbproperror != nil {
		log.Error(ldbproperror)
	}
}

// LoadDBProperties dockwindow
func (cdbv *ChessDBView) LoadDBProperties() error {
	cdbv.dbdock.InsertPGNDB(cdbv.cdb.Displayname, cdbv.cdb.Fullpath, cdbv.cdb.Count)
	// dbtree.SetNameProp(cdb.Displayname)
	// dbtree.SetFileProp(cdb.Fullpath)
	// dbtree.SetCountProp(strconv.Itoa(cdb.Count))
	// dbtree.SetNotesProp(cdb.Notes)
	// dbtree.SetDateModifiedProp(cdb.Filemod)
	return nil
}

// UpdatePGN display
func (cdbv *ChessDBView) UpdatePGN() {
	cdbv.pgndock.SetPGN(cdbv.currentgame)
}

func (cdbv *ChessDBView) setDatabase(dbpath string, kind string) {
	log.Info(fmt.Sprintf("Opening Database with type: %s with key: %s ", dbpath, kind))
	cdbv.cdb, _ = OpenFile(dbpath, kind)
	if cdbv.cdb.CheckIndex {
		needsupdate, err := cdbv.cdb.NeedIndex()
		if err == nil && needsupdate {
			log.Info("Re-indexing database")
			dialog := widgets.NewQProgressDialog2("Re-indexing PGN File", "Cancel", 0, 100, nil, core.Qt__Dialog)
			dialog.SetWindowModality(core.Qt__WindowModal)
			_, err := cdbv.cdb.Index(dialog)
			if err != nil {
				log.Error(err)
			}
		}
	}
	cdbv.model = NewGameListModel(nil)
	cdbv.model.cdb = cdbv.cdb
	cdbv.gamelistdock.tableview.SetModel(cdbv.model)
	cdbv.dbdetaildock.SetProperty("name", cdbv.cdb.Displayname)
	cdbv.dbdetaildock.SetProperty("location", cdbv.cdb.Fullpath)
	cdbv.dbdetaildock.SetProperty("count", strconv.Itoa(cdbv.cdb.Count))
	cdbv.dbdetaildock.SetProperty("modified", cdbv.cdb.Filemod.Format("01/02/2006 03:04:05pm"))

}
