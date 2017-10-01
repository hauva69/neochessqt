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
)

func main() {
	_, err := singleinstance.CreateLockFile("neochess.lock")
	if err != nil {
		os.Exit(1)
	}

	qfile := core.NewQFile2(":qml/translate/en-us.all.json")
	if !qfile.Open(core.QIODevice__ReadOnly | core.QIODevice__Text) {
		log.Error("Could not open translation file")
		os.Exit(1)
	}
	qstream := core.NewQTextStream2(qfile)
	qtext := qstream.ReadAll()

	i18n.ParseTranslationFileBytes("en-us.all.json", []byte(qtext))
	// i18n.MustLoadTranslationFile("en-us.all.json")
	T, _ = i18n.Tfunc("en-US")

	core.QCoreApplication_SetOrganizationDomain("neodevelop.net")
	core.QCoreApplication_SetOrganizationName("NeoDevelop")
	core.QCoreApplication_SetApplicationName("NeoChess")
	core.QCoreApplication_SetApplicationVersion(VERSION)
	qapp := widgets.NewQApplication(len(os.Args), os.Args)
	qwin := widgets.NewQMainWindow(nil, 0)
	if runtime.GOOS == "darwin" {
		qwin.SetUnifiedTitleAndToolBarOnMac(true)
	}
	qwin.SetWindowTitle(core.QCoreApplication_ApplicationName())
	desktop := qapp.Desktop()
	screen := desktop.AvailableGeometry(-1)
	qwin.Resize(core.NewQSize2(screen.Width()*80/100, screen.Height()*90/100))

	var stylefile = core.NewQFile2(":qml/assets/basedark.css")
	if stylefile.Open(core.QIODevice__ReadOnly) {
		styledata := stylefile.ReadAll()
		stylestr := styledata.ConstData()
		qapp.SetStyleSheet(stylestr)
	}

	statusbar := initStatusBar(qwin)
	qwin.SetStatusBar(statusbar)
	statusbar.AddItem("Screen W:%d x H:%d", screen.Width()*80/100, screen.Height()*90/100)

	menu := initMenu(qwin, qapp)
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
	// DisplayLog(qwin)
	log.Info("Starting Application")

	widgets.QApplication_Exec()

}
