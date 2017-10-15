package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameAnalysisDock struct {
	widgets.QDockWidget
	analysis *widgets.QTextEdit
}

func initGameAnalysisDock(w *widgets.QMainWindow) *GameAnalysisDock {
	this := NewGameAnalysisDock("Game Analysis", w, core.Qt__Widget)
	this.analysis = widgets.NewQTextEdit(nil)

	layout := widgets.NewQVBoxLayout()
	this.analysis.SetLayout(layout)
	this.SetAllowedAreas(core.Qt__RightDockWidgetArea)

	this.SetWidget(this.analysis)
	return this
}
