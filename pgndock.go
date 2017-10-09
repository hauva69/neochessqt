package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
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
	figurinefont := gui.NewQFont2("FigurineSymbol T1", 16, 1, false)
	pgndock.editor.SetFontFamily("FigurineSymbol T1")
	pgndock.editor.SetCurrentFont(figurinefont)
	pgndock.editor.SetCurrentFontDefault(figurinefont)
	pgndock.editor.SetFixedWidth(w.Width() / 5)
	pgndock.SetWidget(pgndock.editor)
	return this
}

func (pgndock *PGNDock) SetPGN(game *Game) {
	pgndock.editor.SetHtml(game.CurrentPgn)
}
