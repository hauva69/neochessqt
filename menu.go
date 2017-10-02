package main

import (
	log "github.com/sirupsen/logrus"
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

	importPGNAction := fileMenu.AddAction(T("import_pgn_label"))
	importPGNAction.SetEnabled(true)
	importPGNAction.ConnectTriggered(func(checked bool) {})

	// File / Import SCID

	importSCIDAction := fileMenu.AddAction(T("import_scid_label"))
	importSCIDAction.SetEnabled(true)
	importSCIDAction.ConnectTriggered(func(checked bool) {})

	// File / Import Chessbase

	importCBAction := fileMenu.AddAction(T("import_chessbase_label"))
	importCBAction.SetEnabled(true)
	importCBAction.ConnectTriggered(func(checked bool) {})

	fileMenu.AddSeparator()

	// File / Quit
	quitAction := fileMenu.AddAction(T("quit_label"))
	quitAction.SetEnabled(true)
	quitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))
	quitAction.ConnectTriggered(SaveExit)

	// Play
	playMenu := menu.AddMenu2(T("play_label"))
	playMenu.SetEnabled(true)
	playMenu.AddSeparator()

	// Edit
	editMenu := menu.AddMenu2(T("edit_label"))
	editMenu.SetEnabled(true)
	editMenu.AddSeparator()

	// Game
	gameMenu := menu.AddMenu2(T("game_label"))
	gameMenu.SetEnabled(true)
	gameMenu.AddSeparator()

	// Search
	searchMenu := menu.AddMenu2(T("search_label"))
	searchMenu.SetEnabled(true)
	searchMenu.AddSeparator()

	// Tools
	toolsMenu := menu.AddMenu2(T("tools_label"))
	toolsMenu.SetEnabled(true)
	toolsMenu.AddSeparator()

	// Tools / Log
	logAction := toolsMenu.AddAction(T("log_label"))
	logAction.SetEnabled(true)
	logAction.ConnectTriggered(func(checked bool) { DisplayLog(w) })

	// Help
	helpMenu := menu.AddMenu2(T("help_label"))
	helpMenu.SetEnabled(true)
	helpMenu.AddSeparator()

	// Help / Help Contents

	helpContentsAction := helpMenu.AddAction(T("help_contents_label"))
	helpContentsAction.SetEnabled(true)
	helpContentsAction.ConnectTriggered(func(checked bool) {
		/*
			cmd := exec.Command("/opt/Qt5.8.0/5.8/gcc_64/bin/assistant", "-collectionFile", AppSettings.HelpFile, "-enableRemoteControl")
			err := cmd.Start()
			if err != nil {
				log.Error(err)
			}
			log.Info("Waiting for command to finish...")
			err = cmd.Wait()
			log.Info(fmt.Sprintf("Command finished with error: %v", err))
		*/
		// helpengine := help.NewQHelpEngine(AppSettings.HelpFile, w)
		// helpengine.SetupData()
		// helpbrowser := widgets.NewQTextBrowser(w)
		// helpbrowser.LoadResource()
	})

	// Help / About
	helpMenu.AddSeparator()

	helpAboutAction := helpMenu.AddAction(T("help_about_label"))
	helpAboutAction.SetEnabled(true)
	helpAboutAction.ConnectTriggered(func(checked bool) {
		widgets.QMessageBox_About(w, "NeoChess Database", "Version "+VERSION)
	})

	return this
}

// SaveExit application
func SaveExit(checked bool) {
	if err := AppSettings.Save(); err != nil {
		log.Fatal("Error saving settings")
	}
	Application.Quit()
}
