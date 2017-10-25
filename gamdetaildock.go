package main

import (
	"github.com/rashwell/neochesslib"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// GameDetailDock comment
type GameDetailDock struct {
	widgets.QDockWidget
	tagitem map[string]*widgets.QLineEdit
}

func initGameDetailDock(w *widgets.QMainWindow) *GameDetailDock {
	this := NewGameDetailDock("Game Detail", w, core.Qt__Widget)
	this.tagitem = make(map[string]*widgets.QLineEdit)
	layout := widgets.NewQFormLayout(nil)
	for _, tag := range neochesslib.TagNames() {
		this.tagitem[tag] = widgets.NewQLineEdit2("", nil)
		layout.AddRow3(tag, this.tagitem[tag])
	}
	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	scroll := widgets.NewQScrollArea(nil)
	scroll.SetWidget(widget)
	this.SetWidget(scroll)
	return this
}

// SetGameTags comment
func (gd *GameDetailDock) SetGameTags(g *neochesslib.Game) {
	for _, tag := range neochesslib.TagNames() {
		tagval := g.GetTag(tag)
		gd.tagitem[tag].SetText(tagval)
	}
}
