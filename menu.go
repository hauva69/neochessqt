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
	this := NewMenu(w)
	var menu = this
	fileMenu := menu.AddMenu2("&File")
	fileMenu.SetEnabled(true)
	fileMenu.AddSeparator()

	quitAction := fileMenu.AddAction("&Quit")
	quitAction.SetEnabled(true)
	quitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))
	quitAction.ConnectTriggered(func(checked bool) { a.Quit() })

	return this
}
