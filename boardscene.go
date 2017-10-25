package main

import (
	"github.com/rashwell/neochesslib"
	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardScene struct
type BoardScene struct {
	widgets.QGraphicsScene
	view       *BoardView
	pieceitems map[neochesslib.SquareType]*PieceItem
	dragging   bool
}

func initBoardScene(bv *BoardView) *BoardScene {
	this := NewBoardScene(bv)
	this.view = bv
	this.pieceitems = make(map[neochesslib.SquareType]*PieceItem, 64)
	this.dragging = false
	this.AddSquares()
	this.AddPieces()
	return this
}

// RemovePieces from scene
func (bs *BoardScene) RemovePieces() {
	for k, v := range bs.pieceitems {
		bs.RemoveItem(v.qpix)
		delete(bs.pieceitems, k)
	}
}

// AddPieces to Board Scene
func (bs *BoardScene) AddPieces() {
	for rank := 7; rank >= 0; rank-- {
		for file := 7; file >= 0; file-- {
			var square = neochesslib.CoordsToSquare(rank, file)
			var squareIdx = uint(rank*8 + file)
			piece := bs.view.game.PieceOnSquare(squareIdx)
			pi := initPieceItem(bs, piece, square)
			bs.AddItem(pi.qpix)
			bs.pieceitems[square] = pi
		}
	}
}

// RemovePieceItemOnSquare from scene
func (bs *BoardScene) RemovePieceItemOnSquare(sq neochesslib.SquareType) {
	if bs.pieceitems[sq] != nil {
		pitem := bs.pieceitems[sq]
		bs.RemoveItem(pitem.qpix)
		delete(bs.pieceitems, sq)
	}
}

// AddSquares to board scene
func (bs *BoardScene) AddSquares() {
	// Default Borderwitdh for no Labels
	borderwidth := 10
	topborderwidth := 20

	// Board Labels Adjust border
	if config.IsOption("ShowBoardLables") {
		font := gui.NewQFont2("Arial", 11, 2, false)
		fm := gui.NewQFontMetrics(font)
		borderwidth = fm.Height() * 2
	}

	rightborderwidth := borderwidth
	// Side to Move Marker
	if config.IsOption("ShowSideToMoveMarker") {
		indicatorwidth := bs.view.squaresize / 4
		indicatorspace := bs.view.squaresize / 8
		rightborderwidth = indicatorwidth + indicatorspace
	}

	gpen := gui.NewQPen2(core.Qt__NoPen)
	gtransparent := gui.NewQBrush4(core.Qt__transparent, core.Qt__NoBrush)
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
	// bs.view.UpdateBoardLabels()
	// For Nav Butttons
	scalex := bs.view.squaresize / 2
	scaley := bs.view.squaresize / 2
	posy := float64(8*bs.view.squaresize + borderwidth)
	buttonnumber := 1
	posx := float64(8 * bs.view.squaresize)
	bottomborderwidth := scaley + borderwidth

	lastmoveic := gui.NewQIcon5(":qml/assets/toolbar/ic_last_page_white_48dp_2x.png")
	lastmovebutton := widgets.NewQPushButton(nil)
	lastmovebutton.SetFixedSize2(48, 48)
	lastmovebutton.SetIcon(lastmoveic)
	lastmovebutton.SetCheckable(false)
	lastmovebutton.SetFocusPolicy(core.Qt__NoFocus)
	lastmovebutton.SetIconSize(core.NewQSize2(48, 48))
	lastmovebutton.ConnectClicked(func(checked bool) {
		log.Info("Move to last move in game")
		bs.view.cdbv.currentgame.ToLastMove()
		bs.view.cdbv.UpdatePGN()
	})
	lastmoveproxy := bs.AddWidget(lastmovebutton, core.Qt__Widget)
	lastmoveproxy.SetPos2(posx-float64(buttonnumber*scalex), posy)
	lastmoveproxy.SetFlag(widgets.QGraphicsItem__ItemIsMovable, false)

	buttonnumber++
	nextmoveic := gui.NewQIcon5(":qml/assets/toolbar/ic_chevron_right_white_48dp_2x.png")
	nextmovebutton := widgets.NewQPushButton(nil)
	nextmovebutton.SetFixedSize2(48, 48)
	nextmovebutton.SetIcon(nextmoveic)
	nextmovebutton.SetCheckable(false)
	nextmovebutton.SetFocusPolicy(core.Qt__NoFocus)
	nextmovebutton.SetIconSize(core.NewQSize2(48, 48))
	nextmovebutton.ConnectClicked(func(checked bool) {
		log.Info("Move to next move in game")
		bs.view.cdbv.currentgame.ToNextMove()
		bs.view.cdbv.UpdatePGN()
	})
	nextmoveproxy := bs.AddWidget(nextmovebutton, core.Qt__Widget)
	nextmoveproxy.SetPos2(posx-float64(buttonnumber*scalex), posy)
	nextmoveproxy.SetFlag(widgets.QGraphicsItem__ItemIsMovable, false)

	buttonnumber++
	previousmoveic := gui.NewQIcon5(":qml/assets/toolbar/ic_chevron_left_white_48dp_2x.png")
	previousmovebutton := widgets.NewQPushButton(nil)
	previousmovebutton.SetFixedSize2(48, 48)
	previousmovebutton.SetIcon(previousmoveic)
	previousmovebutton.SetCheckable(false)
	previousmovebutton.SetFocusPolicy(core.Qt__NoFocus)
	previousmovebutton.SetIconSize(core.NewQSize2(48, 48))
	previousmovebutton.ConnectClicked(func(checked bool) {
		log.Info("Move to previous move in game")
		bs.view.cdbv.currentgame.ToPreviousMove()
		bs.view.cdbv.UpdatePGN()
	})
	previousmoveproxy := bs.AddWidget(previousmovebutton, core.Qt__Widget)
	previousmoveproxy.SetPos2(posx-float64(buttonnumber*scalex), posy)
	previousmoveproxy.SetFlag(widgets.QGraphicsItem__ItemIsMovable, false)

	buttonnumber++
	firstmoveic := gui.NewQIcon5(":qml/assets/toolbar/ic_first_page_white_48dp_2x.png")
	firstmovebutton := widgets.NewQPushButton(nil)
	firstmovebutton.SetFixedSize2(48, 48)
	firstmovebutton.SetIcon(firstmoveic)
	firstmovebutton.SetCheckable(false)
	firstmovebutton.SetFocusPolicy(core.Qt__NoFocus)
	firstmovebutton.SetIconSize(core.NewQSize2(48, 48))
	firstmovebutton.ConnectClicked(func(checked bool) {
		log.Info("Move to first move in game")
		bs.view.cdbv.currentgame.ToFirstMove()
		bs.view.cdbv.UpdatePGN()
	})
	firstmoveproxy := bs.AddWidget(firstmovebutton, core.Qt__Widget)
	firstmoveproxy.SetPos2(posx-float64(buttonnumber*scalex), posy)
	firstmoveproxy.SetFlag(widgets.QGraphicsItem__ItemIsMovable, false)

	// For Padding around Border
	item := bs.AddRect2(float64(-1*borderwidth), float64(-1*topborderwidth), float64(bs.view.squaresize*8+rightborderwidth+borderwidth), float64(bs.view.squaresize*8+bottomborderwidth+topborderwidth*2), gpen, gtransparent)
	item.SetFlag(widgets.QGraphicsItem__ItemIsSelectable, false)
}

func (bs *BoardScene) SquareFromPos(pos *core.QPointF) neochesslib.SquareType {
	posX := pos.X()
	posY := pos.Y()
	sqRank := int(posY) / bs.view.squaresize
	sqFile := int(posX) / bs.view.squaresize
	return neochesslib.CoordsToSquare(sqRank, sqFile)
}

//func (bs *BoardScene) AddPiece(pieceitem *PieceItem) {
//
//}
