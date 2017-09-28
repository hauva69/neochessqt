package main

import (
	"os"
	"runtime"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// VERSION of Application
var VERSION = "0.0.1"

func main() {
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

	pgndock := initPGNDock(qwin)
	qwin.AddDockWidget(core.Qt__RightDockWidgetArea, pgndock)

	boardview := initBoardView(qwin)
	qwin.SetCentralWidget(boardview)

	qwin.Show()
	widgets.QApplication_Exec()

}
