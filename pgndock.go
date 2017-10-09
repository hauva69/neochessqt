package main

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// PGNDock struct
type PGNDock struct {
	widgets.QDockWidget
	editor *widgets.QTextEdit
}

func initPGNDock(w *widgets.QMainWindow) *PGNDock {
	this := NewPGNDock("Game", w, core.Qt__Widget)
	var pgndock = this
	pgndock.editor = widgets.NewQTextEdit(nil)
	figurinefont := gui.NewQFont2("FigurineSymbol T1", 16, 1, false)
	pgndock.editor.SetFontFamily("FigurineSymbol T1")
	pgndock.editor.SetCurrentFont(figurinefont)
	pgndock.editor.SetCurrentFontDefault(figurinefont)
	pgndock.editor.SetFixedWidth(w.Width() / 5)
	pgndock.SetWidget(pgndock.editor)
	return this
}

func (pgndock *PGNDock) SetPGN(game *Game) {
	var gamedata = map[string]string{
		"site":  game.Site,
		"white": game.White,
		"black": game.Black,
	}
	t := template.Must(template.New("App").Parse(`<h3 style="text-align:center;">White: <strong>{{.white}}</strong>  vs  Black: <strong>{{.black}}</strong></h3>
		<h4 style="text-align:center;">Site: <strong>{{.site}}</strong></h4><hr/>`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, gamedata); err != nil {
	}
	content := "<style>" + AppSettings.PGNStyle + "</style>"
	content += tpl.String()
	mn := 0
	cv := len(game.Moves) - 1
	for index, Move := range game.Moves {
		if Move.color() == White {
			mn++
			content += "<span class='movenumber'>" + strconv.Itoa(mn) + ". </span>"
		}
		if index == cv {
			content += "<span class='move current'>" + Move.ToSAN() + " </span>"
		} else {
			content += "<span class='move'>" + Move.ToSAN() + " </span>"
		}
	}
	pgndock.editor.SetHtml(content)
}
