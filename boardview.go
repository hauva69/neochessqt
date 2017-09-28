package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardView struct
type BoardView struct {
	widgets.QGraphicsView
	scene *BoardScene
}

func initBoardView(w *widgets.QMainWindow) *BoardView {
	var this = NewBoardView(w)
	var boardview = this
	size := w.FrameSize().Width()
	squaresize := (size / 2) / 8

	boardview.scene = initBoardScene(boardview, squaresize)
	boardview.SetScene(boardview.scene)
	boardview.ConnectResizeEvent(boardview.ResizeBoardView)
	return this
}

// ResizeBoardView event
func (bv *BoardView) ResizeBoardView(event *gui.QResizeEvent) {
	bounding := bv.scene.ItemsBoundingRect()
	if bounding.Width() < 300.0 {
		bounding = core.NewQRectF4(0.0, 0.0, 300.0, 300.0)
	} else {
	}
	bv.FitInView(bounding, core.Qt__KeepAspectRatio)
}
