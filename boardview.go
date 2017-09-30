package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardView struct
type BoardView struct {
	widgets.QGraphicsView
	squaresize  int
	game        *Game
	board       *BoardType
	scene       *BoardScene
	thighlights []*widgets.QGraphicsPathItem
	hhighlight  *widgets.QGraphicsPathItem
}

var (
	highlights = map[string]*gui.QColor{
		"possible": gui.NewQColor3(8, 145, 17, 100),
		"negative": gui.NewQColor3(209, 12, 12, 100),
		"positive": gui.NewQColor3(8, 145, 17, 200),
	}
)

func initBoardView(g *Game, b *BoardType, w *widgets.QMainWindow) *BoardView {
	var this = NewBoardView(w)
	var boardview = this
	boardview.game = g
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

func (bv *BoardView) HighLightSquare(sq SquareType, highlighttype string) {
	gpen := gui.NewQPen2(core.Qt__SolidLine)
	targetcolor := highlights[highlighttype]
	highlightborderwidth := bv.squaresize / 10
	gpen.SetWidth(highlightborderwidth)
	offset := float64(bv.squaresize / 20)
	highlightwidth := float64(bv.squaresize) - offset*2
	gpen.SetColor(targetcolor)
	gtransparent := gui.NewQBrush4(core.Qt__transparent, core.Qt__NoBrush)
	path := gui.NewQPainterPath()
	path.AddRoundedRect2(float64(sq.file()*bv.squaresize)+offset, float64(sq.rank()*bv.squaresize)+offset, highlightwidth, highlightwidth, float64(highlightborderwidth), float64(highlightborderwidth), core.Qt__AbsoluteSize)
	pothighlight := bv.scene.AddPath(path, gpen, gtransparent)
	pothighlight.SetFlag(widgets.QGraphicsItem__ItemIsSelectable, false)
	bv.thighlights = append(bv.thighlights, pothighlight)
}

func (bv *BoardView) RemoveHighlightMoves() {
	for _, highlight := range bv.thighlights {
		if highlight != nil {
			bv.scene.RemoveItem(highlight)
		}
	}
	bv.thighlights = nil
}

func (bv *BoardView) HighlightMovesFrom(from SquareType) {
	for _, sq := range bv.game.GetTargetSquares(from) {
		bv.HighLightSquare(sq, "possible")
	}
}
