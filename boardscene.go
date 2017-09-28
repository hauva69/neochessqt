package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// BoardScene struct
type BoardScene struct {
	widgets.QGraphicsScene
}

func initBoardScene(bv *BoardView, squaresize int) *BoardScene {
	var this = NewBoardScene(bv)
	this.AddSquares(squaresize)
	return this
}

// AddSquares to board scene
func (bs *BoardScene) AddSquares(size int) {
	gpen := gui.NewQPen2(core.Qt__NoPen)
	// gtransparent := gui.NewQBrush4(core.Qt__transparent, core.Qt__NoBrush)
	darkpixmap := gui.NewQPixmap5(":qml/assets/darktexture.png", "png", core.Qt__AutoColor)
	gbrushdarktex := gui.NewQBrush6(core.Qt__white, darkpixmap)
	lightpixmap := gui.NewQPixmap5(":qml/assets/lighttexture.png", "png", core.Qt__AutoColor)
	gbrushlighttex := gui.NewQBrush6(core.Qt__white, lightpixmap)

	// Place Board Square Textures
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			var item *widgets.QGraphicsRectItem
			if i%2 == j%2 {
				item = bs.AddRect2(float64(j*size), float64(i*size), float64(size), float64(size), gpen, gbrushlighttex)
			} else {
				item = bs.AddRect2(float64(j*size), float64(i*size), float64(size), float64(size), gpen, gbrushdarktex)
			}
			item.SetFlag(widgets.QGraphicsItem__ItemIsSelectable, false)
			item.SetAcceptDrops(true)
			//      item.ConnectDragEnterEvent(boardview.DragEnter)
		}
	}
}
