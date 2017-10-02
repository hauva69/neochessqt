package main

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/therecipe/qt/gui"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// OptionType comment
type OptionType struct {
	Version  string
	Key      string
	Label    string
	Kind     string
	Descr    string
	Boolval  bool
	Strval   string
	Intval   int
	Colorval *gui.QColor
}

// AppConfig comment
type AppConfig struct {
	App          *widgets.QApplication
	Window       *widgets.QMainWindow
	SettingsFile string
	Datadir      string
	PGNStyle     string
	HelpFile     string
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

// SetIntOption of key
func (ac *AppConfig) SetIntOption(key string, val int) {
	for _, option := range ac.Options {
		if option.Key == key {
			option.Intval = val
			break
		}
	}
}

// SetStrOption of key
func (ac *AppConfig) SetStrOption(key string, val string) {
	for _, option := range ac.Options {
		if option.Key == key {
			option.Strval = val
			break
		}
	}
}

// GetBoolOption value of key
func (ac *AppConfig) GetBoolOption(key string) bool {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Boolval
		}
	}
	return false
}

// GetIntOption value of key
func (ac *AppConfig) GetIntOption(key string) int {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Intval
		}
	}
	return -1
}

// GetStrOption value of key
func (ac *AppConfig) GetStrOption(key string) string {
	for _, option := range ac.Options {
		if option.Key == key {
			return option.Strval
		}
	}
	return ""
}

// initAppConfig initialize
func initAppConfig(qapp *widgets.QApplication, qwin *widgets.QMainWindow) *AppConfig {
	appconfig := new(AppConfig)
	appconfig.App = qapp
	appconfig.Window = qwin
	appconfig.Datadir = core.QStandardPaths_StandardLocations(core.QStandardPaths__AppDataLocation)[0]
	appconfig.SettingsFile = appconfig.Datadir + "/settings.gob"
	if err := os.MkdirAll(appconfig.Datadir, os.ModePerm); err != nil {
		log.Fatal("Error creating application data directory")
	}
	if !appconfig.Load() {
		desktop := qapp.Desktop()
		screenrect := desktop.AvailableGeometry(-1)
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "LastWidth", "Last Application Width", "int", "Last width of application", false, "", int(screenrect.Width() * 80 / 100), nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "LastHeight", "Last Application Height", "int", "Last height of application", false, "", int(screenrect.Height() * 90 / 100), nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowBoardLables", "Show Board Labels", "bool", "Show Algebraic labels on the edges of the board.", true, "", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowSquareLables", "Show Square Labels", "bool", "Show Algebraic labels on each square of board.", false, "", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "ShowSideToMoveMarker", "Show Side to move", "bool", "Display indicator to side of the board.", true, "", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "StyleFile", "Style File", "file", "CSS Files (*.css)", false, appconfig.Datadir + "/basedark.css", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "PGNStyleFile", "PGN Style File", "file", "CSS Files (*.css)", false, appconfig.Datadir + "/pgntextstyle.css", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "HelpFile", "Help File", "file", "Help File (*.qch)", false, appconfig.Datadir + "/neochess_US.qch", 0, nil})
		appconfig.Options = append(appconfig.Options, OptionType{"1.0.0", "PossibleMove", "Possible Move Color", "color", "Possible Move Color", false, "", 0, gui.NewQColor3(8, 145, 17, 100)})
	}

	helpfile := appconfig.GetStrOption("HelpFile")
	if _, err := os.Stat(helpfile); err == nil {
		appconfig.HelpFile = helpfile
	} else {
		var file = core.NewQFile2(":qml/help/neochess_US.qch")
		if file.Open(core.QIODevice__ReadOnly) {
			qdata := file.ReadAll()
			datastr := qdata.ConstData()
			err = ioutil.WriteFile(helpfile, []byte(datastr), 0644)
			if err != nil {
				log.Fatalf("Error writing file: %v", err)
			}
			appconfig.HelpFile = helpfile
		}
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

	pgnstylefile := appconfig.GetStrOption("PGNStyleFile")
	if _, err := os.Stat(pgnstylefile); err == nil {
		pgnstylebytes, err := ioutil.ReadFile(pgnstylefile)
		if err == nil {
			appconfig.PGNStyle = string(pgnstylebytes)
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
			appconfig.PGNStyle = datastr
		}
	}

	return appconfig
}

// Load Config settings
func (ac *AppConfig) Load() bool {
	log.Info("Loading Config")
	if _, err := os.Stat(ac.SettingsFile); err != nil {
		return false
	}
	gobdata, err := os.Open(ac.SettingsFile)
	defer gobdata.Close()
	if err != nil {
		return false
	}
	decoder := gob.NewDecoder(gobdata)
	if err := decoder.Decode(ac.Options); err != nil {
		return false
	}
	return true
}

// Save comment
func (ac *AppConfig) Save() error {
	log.Info("Saving Config")
	ac.SetIntOption("LastWidth", int(ac.Window.Width()))
	ac.SetIntOption("LastHeight", int(ac.Window.Height()))
	file, err := os.Create(ac.SettingsFile)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(ac.Options)
	}
	file.Close()
	return err
}

// EditConfig Dialog
func (ac *AppConfig) EditConfig() {
	configdialog := widgets.NewQDialog(ac.Window, core.Qt__Dialog)
	fbox := widgets.NewQFormLayout(nil)
	for _, option := range ac.Options {
		switch option.Kind {
		case "bool":
			item := widgets.NewQCheckBox(nil)
			if option.Boolval {
				item.SetCheckState(core.Qt__Checked)
			}
			fbox.AddRow3(option.Label, item)
		case "string":
			item := widgets.NewQLineEdit2(option.Strval, nil)
			item.Home(true)
			fbox.AddRow3(option.Label, item)
		case "file":
			layout := widgets.NewQHBoxLayout()
			item := widgets.NewQLineEdit2(option.Strval, nil)
			item.Home(true)
			button := widgets.NewQPushButton2("...", nil)
			button.ConnectClicked(func(checked bool) {
				fileDialog := widgets.NewQFileDialog2(ac.Window, option.Label, ac.Datadir, option.Descr)
				fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
				fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
				if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
					return
				}
				filename := fileDialog.SelectedFiles()[0]
				log.Infof("Picked: %s", filename)
				item.SetText(filename)
			})
			layout.AddWidget(item, 0, core.Qt__AlignLeft)
			layout.AddWidget(button, 0, core.Qt__AlignRight)
			widget := widgets.NewQWidget(nil, core.Qt__Widget)
			widget.SetLayout(layout)
			fbox.AddRow3(option.Label, widget)
		case "color":
			button := widgets.NewQPushButton2(option.Colorval.Name(), nil)
			button.SetStyleSheet("QPushButton {background-color: " + option.Colorval.Name() + ";}")
			button.ConnectClicked(func(checked bool) {
				log.Info("Picked Color Picker")
				colorDialog := widgets.NewQColorDialog2(option.Colorval, nil)
				color := colorDialog.GetColor(option.Colorval, ac.Window, option.Label, widgets.QColorDialog__ShowAlphaChannel)
				if color.IsValid() {
					log.Infof("Picked Color: %s", color.Name())
					button.SetStyleSheet("QPushButton {background-color: " + color.Name() + ";}")
					button.SetText(color.Name())
				}
			})
			fbox.AddRow3(option.Label, button)
		case "int":
			strint := strconv.Itoa(option.Intval)
			item := widgets.NewQLineEdit2(strint, nil)
			validator := gui.NewQIntValidator2(0, 10000, nil)
			item.SetValidator(validator)
			fbox.AddRow3(option.Label, item)
		}
	}

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

	fbox.AddRow5(buttonBox)

	configdialog.SetLayout(fbox)

	if configdialog.Exec() != int(widgets.QDialog__Accepted) {
		log.Info("Canceled option edit")
	} else {
		log.Info("Options editied changes")
		/*
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
		*/
	}
}
