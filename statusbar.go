package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// StatusBar struct
type StatusBar struct {
	widgets.QStatusBar
}

func initStatusBar(w *widgets.QMainWindow) *StatusBar {
	var this = NewStatusBar(w)
	return this
}

// AddItem to statusbar
func (sb *StatusBar) AddItem(format string, args ...interface{}) {
	item := fmt.Sprintf(format, args...)
	label := widgets.NewQLabel2(item, sb, core.Qt__Widget)
	sb.AddWidget(label, 0)
}
