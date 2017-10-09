package main

import (
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// ChessDBView A view of an Open Chess Database
type ChessDBView struct {
	mainw          *widgets.QMainWindow
	cdb            *ChessDataBase
	pgndock        *PGNDock
	gamelistdock   *GameListDock
	gamedetaildock *GameDetailDock
	boardview      *BoardView
	currentgame    *Game
	currentboard   *BoardType
}

// initCDBView Create instance
func initCDBView(w *widgets.QMainWindow) *ChessDBView {
	view := &ChessDBView{}
	view.mainw = w
	view.cdb = nil

	view.pgndock = initPGNDock(w)
	w.AddDockWidget(core.Qt__RightDockWidgetArea, view.pgndock)

	view.gamelistdock = initGameListDock(w, view)
	w.AddDockWidget(core.Qt__BottomDockWidgetArea, view.gamelistdock)

	view.gamedetaildock = initGameDetailDock(w)
	w.AddDockWidget(core.Qt__RightDockWidgetArea, view.gamedetaildock)
	w.TabifyDockWidget(view.pgndock, view.gamedetaildock)
	w.SetTabPosition(core.Qt__RightDockWidgetArea, widgets.QTabWidget__North)

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
	cdbv.currentgame = NewGame()
	cdbv.currentboard = NewBoard()
	cdbv.currentboard.InitFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	cdbv.currentgame.LoadMoves(cdbv.currentboard)
	cdbv.boardview = initBoardView(cdbv)
	cdbv.mainw.SetCentralWidget(cdbv.boardview)
	cdbv.gamedetaildock.SetGameTags(cdbv.currentgame)
	cdbv.pgndock.SetPGN(cdbv.currentgame)
	var gamerow []string
	gamerow = make([]string, 10)
	gamerow[0] = "-1"
	gamerow[1] = cdbv.currentgame.GameHeader.Event
	gamerow[2] = cdbv.currentgame.GameHeader.Site
	gamerow[3] = cdbv.currentgame.GameHeader.Date
	gamerow[4] = cdbv.currentgame.GameHeader.Round
	gamerow[5] = cdbv.currentgame.GameHeader.White
	gamerow[6] = cdbv.currentgame.GameHeader.Black
	gamerow[7] = cdbv.currentgame.GameHeader.Result
	gamerow[8] = cdbv.currentgame.GameHeader.ECO
	gamerow[9] = cdbv.currentgame.GameHeader.Opening
	cdbv.gamelistdock.AddRow(gamerow)
}

// LoadSelectedGame from gamelist
func (cdbv *ChessDBView) LoadSelectedGame(event *gui.QMouseEvent) {
	item, err := cdbv.gamelistdock.SelectedGame()
	if err != nil {
		cdbv.gamelistdock.table.MouseDoubleClickEventDefault(event)
	}
	log.Infof("Loading game with id: %d", item)
	cdbv.LoadGame(item)
	cdbv.gamelistdock.table.MouseDoubleClickEventDefault(event)
}

// LoadGame into view
func (cdbv *ChessDBView) LoadGame(index int) {
	pgn, err := cdbv.cdb.GetGame(index)
	if err != nil {
		return
	}
	cdbv.currentgame, cdbv.currentboard = ParseGameString([]byte(pgn), index, false)
	cdbv.currentgame.LoadMoves(cdbv.currentboard)
	cdbv.boardview.SetGame()
	cdbv.gamedetaildock.SetGameTags(cdbv.currentgame)
	cdbv.pgndock.SetPGN(cdbv.currentgame)
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
	liniterror := cdbv.LoadInitialGamesList()
	if liniterror != nil {
		log.Error(liniterror)
	}
	ldbproperror := cdbv.LoadDBProperties()
	if ldbproperror != nil {
		log.Error(ldbproperror)
	}
}

// LoadDBProperties dockwindow
func (cdbv *ChessDBView) LoadDBProperties() error {
	// dbtree.SetNameProp(cdb.Displayname)
	// dbtree.SetFileProp(cdb.Fullpath)
	// dbtree.SetCountProp(strconv.Itoa(cdb.Count))
	// dbtree.SetNotesProp(cdb.Notes)
	// dbtree.SetDateModifiedProp(cdb.Filemod)
	return nil
}

func (cdbv *ChessDBView) UpdatePGN() {
	cdbv.pgndock.SetPGN(cdbv.currentgame)
}

// LoadInitialGamesList on ChessDatabase Open
func (cdbv *ChessDBView) LoadInitialGamesList() error {
	log.Info("Loading initial games list")
	glist := make([][]string, 10)
	for i := 0; i < 10; i++ {
		glist[i] = make([]string, 10)
	}
	for gameindex := 1; gameindex < 10; gameindex++ {
		gamebytes, err := cdbv.cdb.GetGame(gameindex)
		if err != nil {
			log.Error(err)
		}
		game, _ := ParseGameString(gamebytes, 1, true)
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
	cdbv.gamelistdock.SetRows(glist)
	return nil
}
