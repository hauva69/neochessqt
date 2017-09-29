package main

import (
	log "github.com/sirupsen/logrus"
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
	scaledqpix := qpix.Scaled2(bs.squaresize, bs.squaresize, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
	pi := new(PieceItem)
	pi.qpix = widgets.NewQGraphicsPixmapItem2(scaledqpix, nil)
	pi.qpix.SetTransformationMode(core.Qt__SmoothTransformation)

	pi.qpix.SetPos2(float64(s.file()*bs.squaresize), float64(s.rank()*bs.squaresize))
	pi.qpix.SetFlag(widgets.QGraphicsItem__ItemIsMovable, true)
	pi.qpix.ConnectMouseReleaseEvent(pi.MouseRelease)
	pi.qpix.ConnectMouseMoveEvent(pi.MouseMove)
	pi.qpix.ConnectMousePressEvent(pi.MousePress)
	cursor := gui.NewQCursor2(core.Qt__OpenHandCursor)
	pi.qpix.SetCursor(cursor)
	pi.piece = p
	pi.square = s
	return pi
}

// MouseRelease Piece is released
func (pi *PieceItem) MouseRelease(event *widgets.QGraphicsSceneMouseEvent) {
	log.Error("Mouse Released")
	pi.qpix.MouseReleaseEventDefault(event)
}

// MousePress Piece is pressed
func (pi *PieceItem) MousePress(event *widgets.QGraphicsSceneMouseEvent) {
	log.Warn("Mouse Pressed")
	pi.qpix.MousePressEventDefault(event)
}

// MouseMove Piece is moved
func (pi *PieceItem) MouseMove(event *widgets.QGraphicsSceneMouseEvent) {
	// Log.Info("Mouse Moved")
	pi.qpix.MouseMoveEventDefault(event)
}
