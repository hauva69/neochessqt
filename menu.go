package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/help"
	"github.com/therecipe/qt/widgets"
)

// Menu struct
type Menu struct {
	widgets.QMenuBar
}

func initMenu(w *widgets.QMainWindow, a *widgets.QApplication, cdbv *ChessDBView) *Menu {

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
	importPGNAction.ConnectTriggered(func(checked bool) { cdbv.loadpgndb(w) })

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
	quitAction.ConnectTriggered(func(checked bool) { SaveExit(a) })

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
	helpContentsAction.ConnectTriggered(func(checked bool) { DisplayHelp(w) })

	// Help / About
	helpMenu.AddSeparator()

	helpAboutAction := helpMenu.AddAction(T("help_about_label"))
	helpAboutAction.SetEnabled(true)
	helpAboutAction.ConnectTriggered(func(checked bool) { DisplayAbout(w) })

	return this
}

// SaveExit application
func SaveExit(a *widgets.QApplication) {
	if err := config.Save(); err != nil {
		log.Fatal("Error saving settings")
	}
	a.Quit()
}

func DisplayHelp(w *widgets.QMainWindow) {
	log.Infof("HelpEngine from: %s", config.HelpFile)
	helpengine := help.NewQHelpEngine(config.HelpFile, w)
	helpengine.SetupData()
	helpdialog := widgets.NewQDialog(w, core.Qt__Dialog)
	hbox := widgets.NewQHBoxLayout()
	tWidget := widgets.NewQTabWidget(nil)
	tWidget.SetMaximumWidth(200)
	tWidget.AddTab(helpengine.ContentWidget(), "Content")
	tWidget.AddTab(helpengine.IndexWidget(), "Index")

	textViewer := widgets.NewQTextBrowser(nil)
	url := core.NewQUrl3("qthelp://net.neodevelop.neochess/doc/index.html", core.QUrl__TolerantMode)
	log.Infof("D: %s", helpengine.FileData(url).Data())
	textViewer.ConnectLoadResource(func(ty int, name *core.QUrl) *core.QVariant {
		if name.Scheme() == "qthelp" {
			return core.NewQVariant15(helpengine.FileData(name))
		}
		return textViewer.LoadResourceDefault(ty, name)
	})
	textViewer.SetSource(url)
	helpengine.ContentWidget().ConnectLinkActivated(func(link *core.QUrl) {
		textViewer.SetSource(link)
	})
	helpengine.IndexWidget().ConnectLinkActivated(func(link *core.QUrl, keyword string) {
		textViewer.SetSource(link)
	})
	horizSplitter := widgets.NewQSplitter2(core.Qt__Horizontal, nil)
	horizSplitter.InsertWidget(0, tWidget)
	horizSplitter.InsertWidget(1, textViewer)
	// horizSplitter.Hide()
	hbox.AddWidget(horizSplitter, 0, core.Qt__AlignTop)
	helpdialog.SetLayout(hbox)
	if helpdialog.Exec() != int(widgets.QDialog__Accepted) {
		log.Info("Canceled option edit")
	} else {
		log.Info("Options editied changes")
	}
}
