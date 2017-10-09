package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type GameDetailDock struct {
	widgets.QDockWidget
	tagitem map[string]*widgets.QLineEdit
}

func initGameDetailDock(w *widgets.QMainWindow) *GameDetailDock {
	this := NewGameDetailDock("Game Detail", w, core.Qt__Widget)
	this.tagitem = make(map[string]*widgets.QLineEdit)
	layout := widgets.NewQFormLayout(nil)
	for k := range supportedtags {
		this.tagitem[k] = widgets.NewQLineEdit2("", nil)
		layout.AddRow3(k, this.tagitem[k])
	}
	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	scroll := widgets.NewQScrollArea(nil)
	scroll.SetWidget(widget)
	this.SetWidget(scroll)
	return this
}

func (gd *GameDetailDock) SetGameTags(g *Game) {
	for k := range supportedtags {
		tagval := g.gettag(k)
		gd.tagitem[k].SetText(tagval)
	}
}
