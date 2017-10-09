package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// PieceItem struct for Graphics Piece
type PieceItem struct {
	qpix            *widgets.QGraphicsPixmapItem
	scene           *BoardScene
	piece           PieceType
	square          SquareType
	nowmoving       bool
	targethighlight *widgets.QGraphicsPathItem
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
	pi.nowmoving = false
	return pi
}

// MouseRelease Piece is released
func (pi *PieceItem) MouseRelease(event *widgets.QGraphicsSceneMouseEvent) {
	log.Info("Mouse Released")
	pi.scene.dragging = false
	pi.nowmoving = false
	if pi.targethighlight != nil {
		pi.scene.RemoveItem(pi.targethighlight)
	}
	pi.scene.view.RemoveHighlightMoves()
	pos := event.ScenePos()
	toSq := pi.scene.SquareFromPos(pos)
	fromSq := pi.square
	piecescale := pi.scene.view.squaresize
	move := fromSq.ToRune() + toSq.ToRune()
	legal, mindex := pi.scene.view.game.IsMoveInFromMoves(move)
	if legal {
		log.Info("Making Move: " + move)
		gamemove := pi.scene.view.game.PotentialMoves[mindex]
		pi.scene.view.board.MakeMove(gamemove, true)
		pi.scene.view.game.Moves = append(pi.scene.view.game.Moves, gamemove)
		pi.scene.view.game.LoadMoves(pi.scene.view.board)
		if gamemove.captures() {
			pi.scene.RemovePieceItemOnSquare(toSq)
		}
		delete(pi.scene.pieceitems, fromSq)
		pi.square = toSq
		pi.qpix.SetPos2(float64(toSq.file()*piecescale), float64(toSq.rank()*piecescale))
		pi.scene.pieceitems[toSq] = pi
		pi.scene.view.cdbv.UpdatePGN()
		//		pi.bv.UpdateSideToMoveIndicator()
	} else {
		log.Info("Illegal Move returning piece.")
		pi.qpix.SetPos2(float64(fromSq.file()*piecescale), float64(fromSq.rank()*piecescale))
	}
	pi.qpix.SetZValue(0.0)
	cursor := gui.NewQCursor2(core.Qt__OpenHandCursor)
	pi.qpix.SetCursor(cursor)
	pi.qpix.MouseReleaseEventDefault(event)
}

// MousePress Piece is pressed
func (pi *PieceItem) MousePress(event *widgets.QGraphicsSceneMouseEvent) {
	if !pi.scene.dragging {
		pi.scene.dragging = true
		log.Info("Mouse Pressed")
		pi.qpix.SetZValue(1.0)
		pi.scene.view.HighlightMovesFrom(pi.square)
		midpiece := float64(pi.scene.view.squaresize / 2)
		pi.qpix.SetPos2(event.ScenePos().X()-midpiece, event.ScenePos().Y()-midpiece)
		cursor := gui.NewQCursor2(core.Qt__ClosedHandCursor)
		pi.qpix.SetCursor(cursor)
		pi.qpix.MousePressEventDefault(event)
	}
}

// MouseMove Piece is moved
func (pi *PieceItem) MouseMove(event *widgets.QGraphicsSceneMouseEvent) {
	if !pi.nowmoving {
		log.Info("Piece " + pi.piece.String() + " dragging.")
		pi.nowmoving = true
	}
	pi.qpix.MouseMoveEventDefault(event)
	squaresize := pi.scene.view.squaresize
	midpiece := float64(squaresize / 2)
	pos := event.ScenePos()
	fromSq := pi.square
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
	sqRank := int(newY+midpiece) / pi.scene.view.squaresize
	sqFile := int(newX+midpiece) / pi.scene.view.squaresize
	toSq := CoordsToSquare(sqRank, sqFile)
	if fromSq != toSq {
		gpen := gui.NewQPen2(core.Qt__SolidLine)
		targetcolor := gui.NewQColor3(209, 12, 12, 100)
		if pi.scene.view.game.IsTarget(fromSq, toSq) {
			targetcolor = gui.NewQColor3(8, 145, 17, 200)
		}
		highlightborderwidth := squaresize / 10
		gpen.SetWidth(highlightborderwidth)
		offset := float64(squaresize / 20)
		highlightwidth := float64(squaresize) - offset*2
		gpen.SetColor(targetcolor)
		gtransparent := gui.NewQBrush4(core.Qt__transparent, core.Qt__NoBrush)
		path := gui.NewQPainterPath()
		path.AddRoundedRect2(float64(toSq.file()*squaresize)+offset, float64(toSq.rank()*squaresize)+offset, highlightwidth, highlightwidth, float64(highlightborderwidth), float64(highlightborderwidth), core.Qt__AbsoluteSize)
		if pi.targethighlight != nil {
			pi.scene.RemoveItem(pi.targethighlight)
		}
		pi.targethighlight = pi.scene.AddPath(path, gpen, gtransparent)
		pi.targethighlight.SetFlag(widgets.QGraphicsItem__ItemIsSelectable, false)
	} else {
		if pi.targethighlight != nil {
			pi.scene.RemoveItem(pi.targethighlight)
		}
	}
}
