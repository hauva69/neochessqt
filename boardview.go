package main

import (
	log "github.com/sirupsen/logrus"
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
	editor      *widgets.QTextEdit
	thighlights []*widgets.QGraphicsPathItem
	hhighlight  *widgets.QGraphicsPathItem
	labels      []*widgets.QGraphicsTextItem
}

var (
	highlights = map[string]*gui.QColor{
		"possible": gui.NewQColor3(8, 145, 17, 100),
		"negative": gui.NewQColor3(209, 12, 12, 100),
		"positive": gui.NewQColor3(8, 145, 17, 200),
	}
)

func initBoardView(g *Game, b *BoardType, e *widgets.QTextEdit, w *widgets.QMainWindow) *BoardView {
	var this = NewBoardView(w)
	var boardview = this
	boardview.game = g
	boardview.board = b
	boardview.editor = e
	size := w.FrameSize().Width()
	boardview.squaresize = (size / 2) / 8

	boardview.scene = initBoardScene(boardview)
	boardview.SetScene(boardview.scene)
	boardview.UpdateBoardLabels()
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

// UpdateBoardLabels display
func (bv *BoardView) UpdateBoardLabels() {
	log.Info("Removing board labels")
	for _, label := range bv.labels {
		if label != nil {
			// Log.Info("Removing Label")
			bv.scene.RemoveItem(label)
		}
	}
	bv.labels = nil
	if AppSettings.IsOption("ShowSquareLables") {
		log.Info("Adding board square labels")
		font := gui.NewQFont2("Arial", 11, 2, false)
		font.SetPixelSize(11)
		twhite := gui.NewQColor2(core.Qt__white)
		// tblack := gui.NewQColor2(core.Qt__black)
		for rank := 7; rank >= 0; rank-- {
			for file := 7; file >= 0; file-- {
				index := uint(rank*8 + file)
				sq := SquareType(index)
				label := bv.scene.AddText(sq.ToRune(), font)
				rank := sq.rank()
				file := sq.file()
				label.SetDefaultTextColor(twhite)
				label.SetPos2(float64(file*bv.squaresize+bv.squaresize-26), float64(rank*bv.squaresize+bv.squaresize-24))
				bv.labels = append(bv.labels, label)
			}
		}
	}
	if AppSettings.IsOption("ShowBoardLables") {
		log.Info("Adding board labels")
		font := gui.NewQFont2("Arial", 11, 2, false)
		fm := gui.NewQFontMetrics(font)
		borderwidth := fm.Height() * 2
		twhite := gui.NewQColor2(core.Qt__white)
		filelbls := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
		for index, lbl := range filelbls {
			lblwidth := fm.Width(lbl, len(lbl))
			lblheight := fm.Height()
			label := bv.scene.AddText(lbl, font)
			label.SetDefaultTextColor(twhite)
			label.SetPos2(float64(index*bv.squaresize+(bv.squaresize/2-lblwidth/2)), float64(8*bv.squaresize+lblheight/4))
			bv.labels = append(bv.labels, label)
		}
		ranklbls := []string{"8", "7", "6", "5", "4", "3", "2", "1"}
		for index, lbl := range ranklbls {
			lblwidth := fm.Width(lbl, len(lbl))
			lblheight := fm.Height()
			label := bv.scene.AddText(lbl, font)
			label.SetDefaultTextColor(twhite)
			label.SetPos2(float64(-1*borderwidth/2-lblwidth/2), float64(index*bv.squaresize+(bv.squaresize/2-lblheight/2)))
			bv.labels = append(bv.labels, label)
		}
	}
}
