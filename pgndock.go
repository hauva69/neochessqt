package main

import (
	"bytes"
	"text/template"

	"github.com/rashwell/neochesslib"
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

func (pgndock *PGNDock) SetPGN(game *neochesslib.Game) {
	var gamedata = map[string]string{
		"site":  game.Site,
		"white": game.White,
		"black": game.Black,
	}
	t := template.Must(template.New("App").Parse(`<h3 style="text-align:center;">White: {{.white}} vs  Black: {{.black}}</h3>
		<h4 style="text-align:center;">Site: {{.site}}</h4><hr/>`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, gamedata); err != nil {
	}
	content := "<style>" + config.PGNStyle + "</style>"
	content += tpl.String()
	content += game.GetPGNMarkup()
	pgndock.editor.SetHtml(content)
}
