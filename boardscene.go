package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardScene struct
type BoardScene struct {
	widgets.QGraphicsScene
	view       *BoardView
	pieceitems map[SquareType]*PieceItem
	dragging   bool
}

func initBoardScene(bv *BoardView) *BoardScene {
	this := NewBoardScene(bv)
	this.view = bv
	this.pieceitems = make(map[SquareType]*PieceItem, 64)
	this.dragging = false
	this.AddSquares()
	this.AddPieces()
	return this
}

func (bs *BoardScene) AddPieces() {
	for rank := 7; rank >= 0; rank-- {
		for file := 7; file >= 0; file-- {
			var square = CoordsToSquare(rank, file)
			var squareIdx = uint(rank*8 + file)
			piece := bs.view.board.Squares[squareIdx]
			pi := initPieceItem(bs, piece, square)
			bs.AddItem(pi.qpix)
			bs.pieceitems[square] = pi
		}
	}
}

func (bs *BoardScene) RemovePieceItemOnSquare(sq SquareType) {
	if bs.pieceitems[sq] != nil {
		pitem := bs.pieceitems[sq]
		bs.RemoveItem(pitem.qpix)
		delete(bs.pieceitems, sq)
	}
}

// AddSquares to board scene
func (bs *BoardScene) AddSquares() {
	gpen := gui.NewQPen2(core.Qt__NoPen)
	// gtransparent := gui.NewQBrush4(core.Qt__transparent, core.Qt__NoBrush)
	darkpixmap := gui.NewQPixmap5(":qml/assets/darktexture.png", "png", core.Qt__AutoColor)
	gbrushdarktex := gui.NewQBrush6(core.Qt__white, darkpixmap)
	lightpixmap := gui.NewQPixmap5(":qml/assets/lighttexture.png", "png", core.Qt__AutoColor)
	gbrushlighttex := gui.NewQBrush6(core.Qt__white, lightpixmap)
	squaresize := bs.view.squaresize
	// Place Board Square Textures
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			var item *widgets.QGraphicsRectItem
			if i%2 == j%2 {
				item = bs.AddRect2(float64(j*squaresize), float64(i*squaresize), float64(squaresize), float64(squaresize), gpen, gbrushlighttex)
			} else {
				item = bs.AddRect2(float64(j*squaresize), float64(i*squaresize), float64(squaresize), float64(squaresize), gpen, gbrushdarktex)
			}
			item.SetFlag(widgets.QGraphicsItem__ItemIsSelectable, false)
			item.SetAcceptDrops(true)
			//      item.ConnectDragEnterEvent(boardview.DragEnter)
		}
	}
}

func (bs *BoardScene) SquareFromPos(pos *core.QPointF) SquareType {
	posX := pos.X()
	posY := pos.Y()
	sqRank := int(posY) / bs.view.squaresize
	sqFile := int(posX) / bs.view.squaresize
	return CoordsToSquare(sqRank, sqFile)
}

//func (bs *BoardScene) AddPiece(pieceitem *PieceItem) {
//
//}
