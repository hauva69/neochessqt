package main

import (
	"os"
	"runtime"

	"github.com/allan-simon/go-singleinstance"
	"github.com/nicksnyder/go-i18n/i18n"
	log "github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var (
	// VERSION of Application
	VERSION = "0.0.1"

	// T translate function
	T i18n.TranslateFunc

	// Application main
	Application *widgets.QApplication

	// AppSettings main
	AppSettings *AppConfig
)

func main() {
	// For now only one instance running at a time
	_, err := singleinstance.CreateLockFile("neochess.lock")
	if err != nil {
		os.Exit(1)
	}

	// Translation Load
	qfile := core.NewQFile2(":qml/translate/en-us.all.json")
	if !qfile.Open(core.QIODevice__ReadOnly | core.QIODevice__Text) {
		log.Error("Could not open translation file")
		os.Exit(1)
	}
	qstream := core.NewQTextStream2(qfile)
	qtext := qstream.ReadAll()
	i18n.ParseTranslationFileBytes("en-us.all.json", []byte(qtext))
	T, _ = i18n.Tfunc("en-US")

	core.QCoreApplication_SetOrganizationDomain("neodevelop.net")
	core.QCoreApplication_SetOrganizationName("NeoDevelop")
	core.QCoreApplication_SetApplicationName("NeoChess")
	core.QCoreApplication_SetApplicationVersion(VERSION)
	Application = widgets.NewQApplication(len(os.Args), os.Args)
	qwin := widgets.NewQMainWindow(nil, 0)
	if runtime.GOOS == "darwin" {
		qwin.SetUnifiedTitleAndToolBarOnMac(true)
	}
	qwin.SetWindowTitle(core.QCoreApplication_ApplicationName())

	log.Info("Starting Application")

	// Load or initialize settings
	AppSettings = initAppConfig(Application, qwin)

	// Desktop Size
	// desktop := qapp.Desktop()
	// screen := desktop.AvailableGeometry(-1)
	qwin.Resize(core.NewQSize2(AppSettings.GetIntOption("LastWidth"), AppSettings.GetIntOption("LastHeight")))

	statusbar := initStatusBar(qwin)
	qwin.SetStatusBar(statusbar)
	statusbar.AddItem("Screen W:%d x H:%d", AppSettings.GetIntOption("LastWidth"), AppSettings.GetIntOption("LastHeight"))

	menu := initMenu(qwin, Application)
	qwin.SetMenuBar(menu)

	toolbar := initToolBar(qwin)
	qwin.AddToolBar2(toolbar)

	dbdock := initDBDock(qwin)
	qwin.AddDockWidget(core.Qt__LeftDockWidgetArea, dbdock)

	pgndock := initPGNDock(qwin)
	qwin.AddDockWidget(core.Qt__RightDockWidgetArea, pgndock)

	cg := NewGame()
	cb := NewBoard()
	cb.InitFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	cg.LoadMoves(cb)
	boardview := initBoardView(cg, cb, pgndock.editor, qwin)
	qwin.SetCentralWidget(boardview)

	qwin.Show()

	widgets.QApplication_Exec()

}
