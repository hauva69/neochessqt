package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// DBDetailDock display struct
type DBDetailDock struct {
	widgets.QDockWidget
	properties map[string]*widgets.QLineEdit
	notes      *widgets.QTextEdit
}

func initDBDetailDock(w *widgets.QMainWindow) *DBDetailDock {
	this := NewDBDetailDock("Database Detail", w, core.Qt__Widget)
	this.properties = make(map[string]*widgets.QLineEdit)
	layout := widgets.NewQFormLayout(nil)

	this.properties["name"] = widgets.NewQLineEdit2("", nil)
	layout.AddRow3("Name", this.properties["name"])

	this.properties["location"] = widgets.NewQLineEdit2("", nil)
	this.properties["location"].SetReadOnly(true)
	layout.AddRow3("Location", this.properties["location"])

	this.properties["count"] = widgets.NewQLineEdit2("", nil)
	this.properties["count"].SetReadOnly(true)
	layout.AddRow3("# Games", this.properties["count"])

	this.properties["modified"] = widgets.NewQLineEdit2("", nil)
	this.properties["modified"].SetReadOnly(true)
	layout.AddRow3("Modified", this.properties["modified"])

	this.notes = widgets.NewQTextEdit(nil)
	layout.AddRow3("Notes", this.notes)

	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)

	this.SetWidget(widget)
	return this
}

func (dbdetail *DBDetailDock) SetProperty(property string, value string) {
	if property == "notes" {
		dbdetail.notes.SetText(value)
	} else {
		dbdetail.properties[property].SetText(value)
	}
}
