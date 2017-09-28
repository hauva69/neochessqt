package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// PieceItem struct for Graphics Piece
type PieceItem struct {
	qpix   *widgets.QGraphicsPixmapItem
	piece  PieceType
	square SquareType
}

func initPieceItem(bs *BoardScene, p PieceType, s SquareType) *PieceItem {
	qpix := gui.NewQPixmap5(p.FileName(), "png", core.Qt__AutoColor)
	scaledqpix := qpix.Scaled2(128, 128, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
	pi := new(PieceItem)
	pi.qpix = widgets.NewQGraphicsPixmapItem2(scaledqpix, nil)
	pi.qpix.SetTransformationMode(core.Qt__SmoothTransformation)

	pi.qpix.SetPos2(float64(s.file()*bs.squaresize), float64(s.rank()*bs.squaresize))
	pi.qpix.SetFlag(widgets.QGraphicsItem__ItemIsMovable, true)
	// pi.qpix.ConnectMouseReleaseEvent(pitem.MouseRelease)
	// pi.qpix.ConnectMouseMoveEvent(pitem.MouseMove)
	// pi.qpix.ConnectMousePressEvent(pitem.MousePress)
	cursor := gui.NewQCursor2(core.Qt__OpenHandCursor)
	pi.qpix.SetCursor(cursor)
	pi.piece = p
	pi.square = s
	return pi
}
