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
	scene  *BoardScene
	piece  PieceType
	square SquareType
}

func initPieceItem(bs *BoardScene, p PieceType, s SquareType) *PieceItem {
	squaresize := bs.view.squaresize
	qpix := gui.NewQPixmap5(p.FileName(), "png", core.Qt__AutoColor)
	scaledqpix := qpix.Scaled2(squaresize, squaresize, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
	pi := new(PieceItem)
	pi.scene = bs
	pi.qpix = widgets.NewQGraphicsPixmapItem2(scaledqpix, nil)
	pi.qpix.SetTransformationMode(core.Qt__SmoothTransformation)
	pi.qpix.SetPos2(float64(s.file()*squaresize), float64(s.rank()*squaresize))
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
	pi.qpix.SetZValue(1.0)
	midpiece := float64(pi.scene.view.squaresize / 2)
	pi.qpix.SetPos2(event.ScenePos().X()-midpiece, event.ScenePos().Y()-midpiece)
	cursor := gui.NewQCursor2(core.Qt__ClosedHandCursor)
	pi.qpix.SetCursor(cursor)
	pi.qpix.MousePressEventDefault(event)
}

// MouseMove Piece is moved
func (pi *PieceItem) MouseMove(event *widgets.QGraphicsSceneMouseEvent) {
	// Log.Info("Mouse Moved")
	pi.qpix.MouseMoveEventDefault(event)
	squaresize := pi.scene.view.squaresize
	midpiece := float64(squaresize / 2)
	pos := event.ScenePos()
	// fromSq := pi.square
	posX := pos.X()
	posY := pos.Y()
	minX := 0.0
	minY := 0.0
	maxX := float64(squaresize * 8)
	maxY := float64(squaresize * 8)
	newX := posX - midpiece
	newY := posY - midpiece
	if posX < minX+midpiece {
		newX = minX
	}
	if posY < minY+midpiece {
		newY = minY
	}
	if posX > maxX-midpiece {
		newX = maxX - midpiece*2
	}
	if posY > maxY-midpiece {
		newY = maxY - midpiece*2
	}
	if posX < minX+midpiece || posY < minY+midpiece || posX > maxX-midpiece || posY > maxY-midpiece {
		pi.qpix.SetPos2(float64(newX), float64(newY))
	}
	// sqRank := int(newY+midpiece) / pi.size
	// sqFile := int(newX+midpiece) / pi.size
	// toSq := CoordsToSquare(sqRank, sqFile)
}
