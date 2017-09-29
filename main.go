package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// VERSION of Application
var VERSION = "0.0.1"

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("neochess.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
		// log.SetOutput(os.Stdout)
	}
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.WithFields(log.Fields{
		"version": VERSION,
	}).Info("Starting Application")
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

	dbdock := initDBDock(qwin)
	qwin.AddDockWidget(core.Qt__LeftDockWidgetArea, dbdock)

	cb := NewBoard()
	cb.InitFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	pgndock := initPGNDock(qwin)
	qwin.AddDockWidget(core.Qt__RightDockWidgetArea, pgndock)

	boardview := initBoardView(cb, qwin)
	qwin.SetCentralWidget(boardview)

	qwin.Show()
	widgets.QApplication_Exec()

}
