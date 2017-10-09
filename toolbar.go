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
	cdbv *ChessDBView
}

func initToolBar(w *widgets.QMainWindow, cdbv *ChessDBView) *ToolBar {
	this := NewToolBar("MainToolBar", w)
	var toolbar = this
	toolbar.cdbv = cdbv
	hdiconstr := "_1x"
	if AppSettings.HDMode {
		toolbar.SetIconSize(core.NewQSize2(48, 48))
		hdiconstr = "_2x"
	} else {
		toolbar.SetIconSize(core.NewQSize2(24, 24))
	}

	toolbar.SetToolButtonStyle(core.Qt__ToolButtonIconOnly) // core.Qt__ToolButtonTextUnderIcon
	toolbar.addbutton(T("new_database_label"), T("new_database_label"), ":/qml/assets/toolbar/ic_create_new_folder_white_48dp"+hdiconstr+".png", defaultbutton)
	toolbar.addbutton(T("open_database_label"), T("open_database_label"), ":/qml/assets/toolbar/ic_folder_white_48dp"+hdiconstr+".png", defaultbutton)
	toolbar.addbutton("Import DB", "Import PGN Database", ":/qml/assets/toolbar/ic_file_upload_white_48dp"+hdiconstr+".png", func(checked bool) { cdbv.loadpgndb(w) })
	toolbar.addbutton("Export DB", "Export Database", ":/qml/assets/toolbar/ic_file_download_white_48dp"+hdiconstr+".png", defaultbutton)
	// toolbar.addbutton("Import PGN", "Import PGN Database", ":/qml/assets/toolbar/ic_file_upload_white_48dp" + hdiconstr+".png", func(checked bool) { loadpgndb(w) })
	//toolbar.addbutton("Export Game", "Export Game", ":/qml/assets/blue-document-export.png", defaultbutton)
	toolbar.addbutton("Search DB", "Search Database", ":/qml/assets/toolbar/ic_search_white_48dp"+hdiconstr+".png", defaultbutton)
	toolbar.addbutton("Options", "Options", ":/qml/assets/toolbar/ic_settings_white_48dp"+hdiconstr+".png", openoptions)
	spacer := widgets.NewQWidget(nil, core.Qt__Widget)
	sp := widgets.NewQSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__ToolButton)
	spacer.SetSizePolicy(sp)
	toolbar.AddWidget(spacer)
	toolbar.addbutton("Help", "Help", ":/qml/assets/toolbar/ic_help_white_48dp"+hdiconstr+".png", func(checked bool) { DisplayHelp(w) })

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

func openoptions(checked bool) {
	log.Info("Options Button Clicked")
	AppSettings.EditConfig()
}
