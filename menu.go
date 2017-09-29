package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Menu struct
type Menu struct {
	widgets.QMenuBar
}

func initMenu(w *widgets.QMainWindow, a *widgets.QApplication) *Menu {

	// File
	this := NewMenu(w)
	var menu = this
	fileMenu := menu.AddMenu2(T("file_menu_label"))
	fileMenu.SetEnabled(true)
	fileMenu.AddSeparator()

	// File / New Database

	newDBAction := fileMenu.AddAction(T("new_database_label"))
	newDBAction.SetEnabled(true)
	newDBAction.ConnectTriggered(func(checked bool) {})

	// File / Open Database

	openDBAction := fileMenu.AddAction(T("open_database_label"))
	openDBAction.SetEnabled(true)
	openDBAction.ConnectTriggered(func(checked bool) {})

	// File / Save Database

	saveDBAction := fileMenu.AddAction(T("save_database_label"))
	saveDBAction.SetEnabled(true)
	saveDBAction.ConnectTriggered(func(checked bool) {})

	// File / Close Database

	closeDBAction := fileMenu.AddAction(T("close_database_label"))
	closeDBAction.SetEnabled(true)
	closeDBAction.ConnectTriggered(func(checked bool) {})

	fileMenu.AddSeparator()

	// File / New Game
	newGameAction := fileMenu.AddAction(T("new_game_label"))
	newGameAction.SetEnabled(true)
	newGameAction.ConnectTriggered(func(checked bool) {})

	// File / Open Game

	openGameAction := fileMenu.AddAction(T("open_game_label"))
	openGameAction.SetEnabled(true)
	openGameAction.ConnectTriggered(func(checked bool) {})

	// File / Save Game

	saveGameAction := fileMenu.AddAction(T("save_game_label"))
	saveGameAction.SetEnabled(true)
	saveGameAction.ConnectTriggered(func(checked bool) {})

	// File / Close Game

	closeGameAction := fileMenu.AddAction(T("close_game_label"))
	closeGameAction.SetEnabled(true)
	closeGameAction.ConnectTriggered(func(checked bool) {})

	fileMenu.AddSeparator()

	// File / Import PGN

	importPGNAction := fileMenu.AddAction("Import PGN")
	importPGNAction.SetEnabled(true)
	importPGNAction.ConnectTriggered(func(checked bool) {})

	// File / Import SCID

	importSCIDAction := fileMenu.AddAction("Import SCID")
	importSCIDAction.SetEnabled(true)
	importSCIDAction.ConnectTriggered(func(checked bool) {})

	// File / Import Chessbase

	importCBAction := fileMenu.AddAction("Import ChessBase")
	importCBAction.SetEnabled(true)
	importCBAction.ConnectTriggered(func(checked bool) {})

	fileMenu.AddSeparator()

	// File / Quit
	quitAction := fileMenu.AddAction("&Quit")
	quitAction.SetEnabled(true)
	quitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))
	quitAction.ConnectTriggered(func(checked bool) { a.Quit() })

	// Play
	playMenu := menu.AddMenu2("Play")
	playMenu.SetEnabled(true)
	playMenu.AddSeparator()

	// Edit
	editMenu := menu.AddMenu2("Edit")
	editMenu.SetEnabled(true)
	editMenu.AddSeparator()

	// Game
	gameMenu := menu.AddMenu2("Game")
	gameMenu.SetEnabled(true)
	gameMenu.AddSeparator()

	// Search
	searchMenu := menu.AddMenu2("Search")
	searchMenu.SetEnabled(true)
	searchMenu.AddSeparator()

	// Tools
	toolsMenu := menu.AddMenu2("Tools")
	toolsMenu.SetEnabled(true)
	toolsMenu.AddSeparator()

	// Tools / Log
	logAction := toolsMenu.AddAction("Log")
	logAction.SetEnabled(true)
	logAction.ConnectTriggered(func(checked bool) { DisplayLog(w) })

	// Help
	helpMenu := menu.AddMenu2("Help")
	helpMenu.SetEnabled(true)
	helpMenu.AddSeparator()

	return this
}
