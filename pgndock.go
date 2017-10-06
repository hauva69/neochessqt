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
	// figurinefont := gui.NewQFont2("ChessAlpha2", 12, 1, false)
	// pgndock.editor.SetFontFamily("ChessAlpha2")
	// pgndock.editor.SetCurrentFont(figurinefont)
	// pgndock.editor.SetCurrentFontDefault(figurinefont)
	pgndock.editor.SetHtml("<p>1. </p>")
	pgndock.editor.SetFixedWidth(w.Width() / 5)
	pgndock.SetWidget(pgndock.editor)
	return this
}
