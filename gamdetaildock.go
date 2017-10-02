package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameDetailDock struct {
	widgets.QDockWidget
	game *Game
}

func initGameDetailDock(g *Game, w *widgets.QMainWindow) *GameDetailDock {
	this := NewGameDetailDock("Game Detail", w, core.Qt__Widget)
	var gamedetaildock = this
	gamedetaildock.game = g
	layout := widgets.NewQFormLayout(nil)
	for k := range supportedtags {
		tagval := g.gettag(k)
		item := widgets.NewQLineEdit2(tagval, nil)
		layout.AddRow3(k, item)
	}
	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	scroll := widgets.NewQScrollArea(nil)
	scroll.SetWidget(widget)
	gamedetaildock.SetWidget(scroll)
	return this
}
