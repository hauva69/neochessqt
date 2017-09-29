package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// ToolBar struct
type ToolBar struct {
	widgets.QToolBar
}

func initToolBar(w *widgets.QMainWindow) *ToolBar {
	this := NewToolBar("MainToolBar", w)
	var toolbar = this
	toolbar.SetIconSize(core.NewQSize2(24, 24))
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextUnderIcon)

	toolbar.addbutton("Create DB", "Create Database", ":/qml/assets/address-book-blue.png", defaultbutton)
	toolbar.addbutton("Open DB", "Open PGN Database", ":/qml/assets/blue-folder-open.png", defaultbutton)
	toolbar.addbutton("Import DB", "Import Database", ":/qml/assets/blue-folder-import.png", defaultbutton)
	toolbar.addbutton("Export DB", "Open Database", ":/qml/assets/blue-folder-export.png", defaultbutton)
	toolbar.addbutton("Import Game", "Import Game", ":/qml/assets/blue-document-import.png", defaultbutton)
	toolbar.addbutton("Export Game", "Export Game", ":/qml/assets/blue-document-export.png", defaultbutton)
	toolbar.addbutton("Search DB", "Search Database", ":/qml/assets/magnifier.png", defaultbutton)
	toolbar.addbutton("Options", "Options", ":/qml/assets/system-monitor.png", defaultbutton)
	spacer := widgets.NewQWidget(nil, core.Qt__Widget)
	sp := widgets.NewQSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__ToolButton)
	spacer.SetSizePolicy(sp)
	toolbar.AddWidget(spacer)
	toolbar.addbutton("Help", "Help", ":/qml/assets/question.png", defaultbutton)

	return this
}

func (tb *ToolBar) addbutton(action string, tip string, iconname string, f func(bool)) {
	ba := tb.AddAction2(gui.NewQIcon5(iconname), action)
	ba.SetStatusTip(tip)
	ba.ConnectTriggered(f)
}

func defaultbutton(checked bool) {
	log.Info("Button Clicked")
}
