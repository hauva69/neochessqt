package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardView struct
type BoardView struct {
	widgets.QGraphicsView
	squaresize int
	board      *BoardType
	scene      *BoardScene
}

func initBoardView(b *BoardType, w *widgets.QMainWindow) *BoardView {
	var this = NewBoardView(w)
	var boardview = this
	boardview.board = b
	size := w.FrameSize().Width()
	boardview.squaresize = (size / 2) / 8

	boardview.scene = initBoardScene(boardview)
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
