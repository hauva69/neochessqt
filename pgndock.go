package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// PGNDock struct
type PGNDock struct {
	widgets.QDockWidget
	editor *widgets.QTextEdit
}

func initPGNDock(w *widgets.QMainWindow) *PGNDock {
	this := NewPGNDock("Game", w, core.Qt__Widget)
	var pgndock = this
	pgndock.editor = widgets.NewQTextEdit(nil)
	pgndock.editor.SetHtml("<p>1. </p>")
	pgndock.editor.SetFixedWidth(w.Width() / 4)
	pgndock.SetWidget(pgndock.editor)
	return this
}
