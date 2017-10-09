package main

import (
	"os"
	"runtime"

	"github.com/allan-simon/go-singleinstance"
	"github.com/boltdb/bolt"
	"github.com/nicksnyder/go-i18n/i18n"
	log "github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var (
	// VERSION of Application
	VERSION = "0.0.1"

	// REVISION Placeholder for git stamping builds
	REVISION = "na"

	// T translate function
	T i18n.TranslateFunc

	// Application main
	Application *widgets.QApplication

	// MainWindow of Application
	MainWindow *widgets.QMainWindow

	// AppSettings main
	AppSettings *AppConfig

	catdb *bolt.DB
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
	MainWindow := widgets.NewQMainWindow(nil, 0)
	if runtime.GOOS == "darwin" {
		MainWindow.SetUnifiedTitleAndToolBarOnMac(true)
	}
	MainWindow.SetWindowTitle(core.QCoreApplication_ApplicationName())

	log.Info("Starting Application")

	// Load or initialize settings
	AppSettings = initAppConfig(Application, MainWindow)

	var caterr error
	catdb, caterr = bolt.Open(AppSettings.Datadir+"/catalog.db", 0644, nil)
	if caterr != nil {
		log.Error(caterr)
	}

	MainWindow.Resize(core.NewQSize2(AppSettings.GetIntOption("LastWidth"), AppSettings.GetIntOption("LastHeight")))

	statusbar := initStatusBar(MainWindow)
	MainWindow.SetStatusBar(statusbar)
	statusbar.AddItem("Screen W:%d x H:%d", AppSettings.GetIntOption("LastWidth"), AppSettings.GetIntOption("LastHeight"))

	chessdbview := initCDBView(MainWindow)

	toolbar := initToolBar(MainWindow, chessdbview)
	MainWindow.AddToolBar2(toolbar)

	dbdock := initDBDock(MainWindow, chessdbview)
	MainWindow.AddDockWidget(core.Qt__LeftDockWidgetArea, dbdock)

	menu := initMenu(MainWindow, Application, chessdbview)
	MainWindow.SetMenuBar(menu)

	MainWindow.Show()

	widgets.QApplication_Exec()

}
