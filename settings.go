package main

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// OptionType comment
type OptionType struct {
	Version string
	Key     string
	Label   string
	Kind    string
	Descr   string
	Boolval bool
	Strval  string
	Intval  int
}

// AppConfig comment
type AppConfig struct {
	SettingsFile string
	Datadir      string
	Options      []OptionType
}

// IsOption true or false
func (ac *AppConfig) IsOption(key string) bool {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Boolval
		}
	}
	return false
}

// SetBoolOption to true or false
func (ac *AppConfig) SetBoolOption(key string, val bool) {
	for _, option := range ac.Options {
		if option.Key == key {
			option.Boolval = val
			break
		}
	}
}

// GetOption value of key
func (ac *AppConfig) GetBoolOption(key string) bool {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Boolval
		}
	}
	return false
}

func (ac *AppConfig) GetIntOption(key string) int {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Intval
		}
	}
	return -1
}

func (ac *AppConfig) GetStrOption(key string) string {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Strval
		}
	}
	return ""
}

// initAppConfig initialize
func initAppConfig(qapp *widgets.QApplication) *AppConfig {
	appconfig := new(AppConfig)
	appconfig.Datadir = core.QStandardPaths_StandardLocations(core.QStandardPaths__AppDataLocation)[0]
	appconfig.SettingsFile = appconfig.Datadir + "/settings.gob"
	if err := os.MkdirAll(appconfig.Datadir, os.ModePerm); err != nil {
		log.Fatal("Error creating application data directory")
	}
	appconfig.Options = LoadConfig(appconfig.SettingsFile)
	if appconfig.Options == nil {
		desktop := qapp.Desktop()
		screenrect := desktop.AvailableGeometry(-1)
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "LastWidth", "Last Application Width", "int", "Last width of application", false, "", int(screenrect.Width() * 80 / 100)})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "LastHeight", "Last Application Height", "int", "Last height of application", false, "", int(screenrect.Height() * 90 / 100)})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowBoardLables", "Show Board Labels", "bool", "Show Algebraic labels on the edges of the board.", true, "", 0})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowSquareLables", "Show Square Labels", "bool", "Show Algebraic labels on each square of board.", false, "", 0})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowSideToMoveMarker", "Show Side to move", "bool", "Display indicator to side of the board.", true, "", 0})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "StyleFile", "Style File", "string", "File for applications style", false, appconfig.Datadir + "/basedark.css", 0})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "PGNStyleFile", "PGN Style File", "string", "File for PGN Text Editor", false, appconfig.Datadir + "/pgntextstyle.css", 0})
	}

	stylefile := appconfig.GetStrOption("StyleFile")
	if _, err := os.Stat(stylefile); err == nil {
		appstylebytes, err := ioutil.ReadFile(stylefile)
		if err == nil {
			appstyle := string(appstylebytes)
			qapp.SetStyleSheet(appstyle)
		}
	} else {
		var file = core.NewQFile2(":qml/assets/basedark.css")
		if file.Open(core.QIODevice__ReadOnly) {
			qdata := file.ReadAll()
			datastr := qdata.ConstData()
			err = ioutil.WriteFile(stylefile, []byte(datastr), 0644)
			if err != nil {
				log.Fatalf("Error writing file: %v", err)
			}
			qapp.SetStyleSheet(datastr)
		}
	}

	/*
		pgnstylefile := appconfig.GetOption("PGNStyleFile")
		if _, err := os.Stat(pgnstylefile); err == nil {
			pgnstylebytes, err := ioutil.ReadFile(pgnstylefile)
			if err == nil {
				// pgnstyle := string(pgnstylebytes)
			}
		} else {
			var file = core.NewQFile2(":qml/assets/pgntextstyle.css")
			if file.Open(core.QIODevice__ReadOnly) {
				qdata := file.ReadAll()
				datastr := qdata.ConstData()
				err = ioutil.WriteFile(pgnstylefile, []byte(datastr), 0644)
				if err != nil {
					log.Fatalf("Error writing file: %v", err)
				}
				// pgnstyle := datastr
			}
		}
	*/

	return appconfig
}

// LoadConfig comment
func LoadConfig(sfile string) []OptionType {
	_, err := os.Stat(sfile)
	if err != nil {
		return nil
	}
	var options []OptionType
	file, err := os.Open(sfile)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(options)
	}
	file.Close()
	if err != nil {
		return nil
	}
	return options
}

// SaveConfig comment
func (ac *AppConfig) SaveConfig() error {
	file, err := os.Create(ac.SettingsFile)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(ac.Options)
	}
	file.Close()
	return err
}

/* EditConfig Dialog
func (a *App) EditConfig() {
	boolchanges = make(map[string]bool)
	strchanges = make(map[string]string)

	configdialog := widgets.NewQDialog(window, core.Qt__Dialog)
	vbox := widgets.NewQVBoxLayout()
	optionstable := widgets.NewQTableWidget2(0, 4, nil)
	optionstable.SetHorizontalHeaderLabels([]string{"Key", "Option", "Value", "Description"})

	propcolor := gui.NewQColor3(52, 57, 61, 255)
	propbrush := gui.NewQBrush3(propcolor, core.Qt__SolidPattern)
	for index, option := range a.Options {
		optionstable.InsertRow(index)
		key := widgets.NewQTableWidgetItem(0)
		key.SetText(option.Key)
		key.SetFlags(key.Flags() ^ core.Qt__ItemIsEditable)
		optionstable.SetItem(index, 0, key)
		label := widgets.NewQTableWidgetItem(0)
		label.SetText(option.Label)
		label.SetFlags(label.Flags() ^ core.Qt__ItemIsEditable)
		optionstable.SetItem(index, 1, label)
		optionstable.Item(index, 1).SetBackground(propbrush)
		if option.Kind == "bool" {
			item := widgets.NewQTableWidgetItem(1)
			item.Data(int(core.Qt__CheckStateRole))
			if option.Boolval {
				item.SetCheckState(core.Qt__Checked)
			} else {
				item.SetCheckState(core.Qt__Unchecked)
			}
			optionstable.SetItem(index, 2, item)
		}
		if option.Kind == "string" {
			item := widgets.NewQTableWidgetItem(2)
			item.SetText(option.Strval)
			optionstable.SetItem(index, 2, item)
		}
		itemdescr := widgets.NewQTableWidgetItem2(option.Descr, 0)
		itemdescr.SetFlags(itemdescr.Flags() ^ core.Qt__ItemIsEditable)
		optionstable.SetItem(index, 3, itemdescr)
	}
	minHeight := 300
	minWidth := 500
	if mainapp.Width/2 > minWidth {
		minWidth = mainapp.Width / 2
	}
	if mainapp.Height/2 > minHeight {
		minHeight = mainapp.Height / 2
	}
	optionstable.SetMinimumHeight(minHeight)
	optionstable.SetMinimumWidth(minWidth)
	optionstable.SetColumnHidden(0, true)
	optionstable.ResizeRowsToContents()
	optionsview := widgets.NewQTableViewFromPointer(widgets.PointerFromQTableWidget(optionstable))
	optionsview.VerticalHeader().Hide()
	optionsview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
	optionsview.HorizontalHeader().SetStretchLastSection(true)
	// optionsview.SetStyleSheet("QHeaderView::section { background-color:#90CAF9; }")

	optionstable.ConnectItemChanged(func(item *widgets.QTableWidgetItem) {
		row := item.Row()
		currentitem := optionstable.Item(row, 0)
		key := currentitem.Text()
		if item.Type() == 2 {
			strchanges[key] = item.Text()
		}
		if item.Type() == 1 {
			if item.CheckState() == core.Qt__Checked {
				boolchanges[key] = true
			} else {
				boolchanges[key] = false
			}
		}
	})

	vbox.AddWidget(optionstable, 0, core.Qt__AlignTop)

	buttonBox := widgets.NewQDialogButtonBox(nil)
	acceptButton := widgets.NewQPushButton2("Apply", nil)
	cancelButton := widgets.NewQPushButton2("Cancel", nil)
	buttonBox.AddButton(acceptButton, widgets.QDialogButtonBox__AcceptRole)
	buttonBox.AddButton(cancelButton, widgets.QDialogButtonBox__RejectRole)

	acceptButton.ConnectPressed(func() {
		configdialog.Done(int(widgets.QDialog__Accepted))
	})

	cancelButton.ConnectPressed(func() {
		configdialog.Done(int(widgets.QDialog__Rejected))
	})

	vbox.AddWidget(buttonBox, 0, core.Qt__AlignBottom)

	configdialog.SetLayout(vbox)
	//		optionwidget := widgets.NewQWidget(nil, core.Qt__Widget)
	//		optionwidget.SetLayout(vbox)
	if configdialog.Exec() != int(widgets.QDialog__Accepted) {
		Log.Info("Canceled option edit")
	} else {
		Log.Info("Options editied changes")
		for k, v := range boolchanges {
			for optindex, option := range mainapp.Options {
				if option.Key == k {
					mainapp.Options[optindex].Boolval = v
					break
				}
			}
		}
		for k, v := range strchanges {
			for optindex, option := range mainapp.Options {
				if option.Key == k {
					mainapp.Options[optindex].Strval = v
					break
				}
			}
		}
		boardview.UpdateBoardLabels()
		boardview.UpdateSideToMoveIndicator()
		err := mainapp.SaveConfig()
		if err != nil {
			Log.Info("Error saving configuration")
		}
	}
}
*/
