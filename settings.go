package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/therecipe/qt/gui"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// OptionType comment
type OptionType struct {
	Group    string `json:"group"`
	Key      string `json:"key"`
	Label    string `json:"label"`
	Kind     string `json:"kind"`
	Descr    string `json:"descr"`
	Boolval  bool   `json:"boolval"`
	Dirval   string `json:"dirval"`
	Strval   string `json:"strval"`
	Intval   int    `json:"intval"`
	Colorval string `json:"colorval"`
	Modified bool   `json:"modified"`
}

// OptionSetType comment
type OptionSetType struct {
	Version string       `json:"version"`
	Options []OptionType `json:"options"`
}

// AppConfig comment
type AppConfig struct {
	App          *widgets.QApplication
	Window       *widgets.QMainWindow
	SettingsFile string
	Datadir      string
	Programdir   string
	PGNStyle     string
	HelpFile     string
	HDMode       bool
	Tabs         []string
	OptionSet    OptionSetType
}

// IsOption true or false
func (ac *AppConfig) IsOption(key string) bool {
	for _, option := range ac.OptionSet.Options {
		if option.Key == key {
			return option.Boolval
		}
	}
	return false
}

// SetBoolOption to true or false
func (ac *AppConfig) SetBoolOption(key string, val bool) {
	for i, option := range ac.OptionSet.Options {
		if option.Key == key {
			ac.OptionSet.Options[i].Boolval = val
			break
		}
	}
}

// SetIntOption of key
func (ac *AppConfig) SetIntOption(key string, val int) {
	for i, option := range ac.OptionSet.Options {
		if option.Key == key {
			ac.OptionSet.Options[i].Intval = val
			break
		}
	}
}

// SetStrOption of key
func (ac *AppConfig) SetStrOption(key string, val string) {
	for i, option := range ac.OptionSet.Options {
		if option.Key == key {
			ac.OptionSet.Options[i].Strval = val
			break
		}
	}
}

// SetColorOption of key
func (ac *AppConfig) SetColorOption(key string, color *gui.QColor) {
	for i, option := range ac.OptionSet.Options {
		if option.Key == key {
			ac.OptionSet.Options[i].Colorval = color.Name()
			break
		}
	}
}

// GetColorOption of key
func (ac *AppConfig) GetColorOption(key string) *gui.QColor {
	for _, option := range ac.OptionSet.Options {
		if option.Key == key {
			return gui.NewQColor6(option.Colorval)
		}
	}
	return nil
}

// GetBoolOption value of key
func (ac *AppConfig) GetBoolOption(key string) bool {
	for _, option := range ac.OptionSet.Options {
		if option.Key == key {
			return option.Boolval
		}
	}
	return false
}

// GetIntOption value of key
func (ac *AppConfig) GetIntOption(key string) int {
	for _, option := range ac.OptionSet.Options {
		if option.Key == key {
			return option.Intval
		}
	}
	return -1
}

// GetStrOption value of key
func (ac *AppConfig) GetStrOption(key string) string {
	for _, option := range ac.OptionSet.Options {
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
	appconfig.Programdir = core.QCoreApplication_ApplicationDirPath()
	appconfig.SettingsFile = appconfig.Datadir + "/settings.json"
	appconfig.Tabs = []string{"Profile", "General", "Board", "PGN", "Engines"}
	if err := os.MkdirAll(appconfig.Datadir, os.ModePerm); err != nil {
		log.Fatal("Error creating application data directory")
	}
	if !appconfig.Load() {
		log.Info("Initializing Setting Storage")
		desktop := qapp.Desktop()
		screenrect := desktop.AvailableGeometry(-1)
		dpi := qapp.Screens()[0].LogicalDotsPerInch()
		log.Infof("DPI :%f", dpi)
		if dpi <= 96.0 {
			appconfig.HDMode = false
		} else {
			appconfig.HDMode = true
		}

		appconfig.OptionSet.Version = "1.0.0"
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "PlayerName", "Player Name", "string", "", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "Country", "Country", "dropdown", "United States of America", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "Title", "Title", "dropdown", "None;GM;WGM;FM;M", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "FIDEID", "FIDE ID", "string", "", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "ELO", "ELO Rating", "string", "", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "USCF", "USCF Rating", "string", "", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Profile", "OTHER", "Other Rating", "string", "", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"General", "LastWidth", "Last Application Width", "int", "Last width of application", false, "", "", int(screenrect.Width() * 80 / 100), "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"General", "LastHeight", "Last Application Height", "int", "Last height of application", false, "", "", int(screenrect.Height() * 90 / 100), "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"General", "HelpFile", "Help File", "file", "Help File (*.qch)", false, appconfig.Datadir, appconfig.Datadir + "/neochess_US.qhc", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"General", "StyleFile", "Style File", "file", "CSS Files (*.css)", false, appconfig.Datadir, appconfig.Datadir + "/basedark.css", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Board", "ShowBoardLables", "Show Board Labels", "bool", "Show Algebraic labels on the edges of the board.", true, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Board", "ShowSquareLables", "Show Square Labels", "bool", "Show Algebraic labels on each square of board.", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Board", "ShowSideToMoveMarker", "Show Side to move", "bool", "Display indicator to side of the board.", true, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Board", "PossibleMove", "Possible Move Color", "color", "Possible Move Color", false, "", "", 0, "#089111", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"PGN", "PGNStyleFile", "PGN Style File", "file", "CSS Files (*.css)", false, appconfig.Datadir, appconfig.Datadir + "/pgntextstyle.css", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"PGN", "PGNPieceCountryDisplay", "PGN Piece Display", "dropdown", "Figurine;English;Dutch", false, "", "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Engines", "Engine1", "Engine #1", "file", "", false, appconfig.Programdir, "", 0, "", false})
		appconfig.OptionSet.Options = append(appconfig.OptionSet.Options, OptionType{"Engines", "Engine2", "Engine #2", "file", "", false, appconfig.Programdir, "", 0, "", false})
	}

	var fontfile = core.NewQFile2(":qml/assets/FIG-TB-1.TTF")
	if fontfile.Open(core.QIODevice__ReadOnly) {
		fontdata := fontfile.ReadAll()
		fntdb := gui.NewQFontDatabase()
		fntid := fntdb.AddApplicationFontFromData(fontdata)
		log.Infof("Font ID: %d", fntid)
		fontfamily := fntdb.ApplicationFontFamilies(fntid)
		log.Infof("Font Family: %v", fontfamily)
	}

	helpfile := appconfig.GetStrOption("HelpFile")
	helpfile2 := strings.Replace(helpfile, ".qhc", ".qch", -1)
	if _, err := os.Stat(helpfile); err == nil {
		appconfig.HelpFile = helpfile
	} else {
		var file = core.NewQFile2(":qml/help/neochess_US.qhc")
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
	if _, err := os.Stat(helpfile2); err != nil {
		var file = core.NewQFile2(":qml/help/neochess_US.qch")
		if file.Open(core.QIODevice__ReadOnly) {
			qdata := file.ReadAll()
			datastr := qdata.ConstData()
			err = ioutil.WriteFile(helpfile2, []byte(datastr), 0644)
			if err != nil {
				log.Fatalf("Error writing file: %v", err)
			}
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
	log.Infof("Loading Config from: %s", ac.SettingsFile)
	if _, err := os.Stat(ac.SettingsFile); err != nil {
		log.Info("Settings file not found")
		return false
	}
	log.Info("Decoding Settings")
	optionsdata, _ := ioutil.ReadFile(ac.SettingsFile)
	err := json.Unmarshal(optionsdata, &ac.OptionSet)
	if err != nil {
		return false
	}
	log.Info("Reading Settings")
	return true
}

// Save comment
func (ac *AppConfig) Save() error {
	log.Info("Saving Config")
	ac.SetIntOption("LastWidth", int(ac.Window.Width()))
	ac.SetIntOption("LastHeight", int(ac.Window.Height()))
	optionsJSON, _ := json.MarshalIndent(ac.OptionSet, "", "\t")
	err := ioutil.WriteFile(ac.SettingsFile, optionsJSON, 0644)
	return err
}

// EditConfig Dialog
func (ac *AppConfig) EditConfig() {
	configdialog := widgets.NewQDialog(ac.Window, core.Qt__Dialog)
	configdialog.SetMinimumSize2(600, 400)
	vbox := widgets.NewQVBoxLayout()
	configdialog.SetLayout(vbox)

	maintabwidget := widgets.NewQTabWidget(nil)
	vbox.AddWidget(maintabwidget, 0, core.Qt__AlignTop)

	var tabs map[string]*widgets.QWidget
	tabs = make(map[string]*widgets.QWidget)

	var forms map[string]*widgets.QFormLayout
	forms = make(map[string]*widgets.QFormLayout)

	for i := 0; i < len(ac.Tabs); i++ {
		tabs[ac.Tabs[i]] = widgets.NewQWidget(nil, core.Qt__Widget)
		forms[ac.Tabs[i]] = widgets.NewQFormLayout(nil)
		tabs[ac.Tabs[i]].SetLayout(forms[ac.Tabs[i]])
		maintabwidget.AddTab(tabs[ac.Tabs[i]], ac.Tabs[i])
	}

	for _, option := range ac.OptionSet.Options {
		key := option.Key
		switch option.Kind {
		case "bool":
			item := widgets.NewQCheckBox(nil)
			if ac.GetBoolOption(key) {
				item.SetCheckState(core.Qt__Checked)
			}
			forms[option.Group].AddRow3(option.Label, item)
			item.ConnectStateChanged(func(state int) {
				if state == int(core.Qt__Unchecked) {
					ac.SetBoolOption(key, false)
				} else {
					ac.SetBoolOption(key, true)
				}
			})
		case "dropdown":
			item := widgets.NewQComboBox(nil)
			vals := strings.Split(option.Descr, ";")
			item.AddItems(vals)
			item.SetCurrentText(ac.GetStrOption(key))
			forms[option.Group].AddRow3(option.Label, item)
			item.ConnectCurrentIndexChanged2(func(val string) {
				ac.SetStrOption(key, val)
			})
		case "string":
			item := widgets.NewQLineEdit2(ac.GetStrOption(key), nil)
			item.Home(true)
			forms[option.Group].AddRow3(option.Label, item)
			item.ConnectTextEdited(func(text string) {
				ac.SetStrOption(key, text)
			})
		case "file":
			layout := widgets.NewQHBoxLayout()
			item := widgets.NewQLineEdit2(ac.GetStrOption(key), nil)
			item.Home(true)
			button := widgets.NewQPushButton2("...", nil)
			button.ConnectClicked(func(checked bool) {
				fileDialog := widgets.NewQFileDialog2(ac.Window, option.Label, option.Dirval, option.Descr)
				fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
				fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
				if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
					return
				}
				filename := fileDialog.SelectedFiles()[0]
				log.Infof("Picked: %s", filename)
				item.SetText(filename)
				ac.SetStrOption(key, filename)
			})
			layout.AddWidget(item, 0, core.Qt__AlignTop)
			layout.AddWidget(button, 0, core.Qt__AlignTop)
			widget := widgets.NewQWidget(nil, core.Qt__Widget)
			widget.SetLayout(layout)
			forms[option.Group].AddRow3(option.Label, widget)
			item.ConnectTextEdited(func(text string) {
				log.Infof("For key: %s file changed to: %s", key, text)
				ac.SetStrOption(key, text)
			})
		case "color":
			button := widgets.NewQPushButton2(option.Colorval, nil)
			button.SetStyleSheet("QPushButton {background-color: " + option.Colorval + ";}")
			button.ConnectClicked(func(checked bool) {
				log.Info("Picked Color Picker")
				colorDialog := widgets.NewQColorDialog2(ac.GetColorOption(key), nil)
				color := colorDialog.GetColor(ac.GetColorOption(key), ac.Window, option.Label, widgets.QColorDialog__ShowAlphaChannel)
				if color.IsValid() {
					log.Infof("Picked Color: %s", color.Name())
					button.SetStyleSheet("QPushButton {background-color: " + color.Name() + ";}")
					button.SetText(color.Name())
					ac.SetColorOption(key, color)
				}
			})
			forms[option.Group].AddRow3(option.Label, button)
		case "int":
			i := ac.GetIntOption(key)
			st := strconv.Itoa(i)
			item := widgets.NewQLineEdit2(st, nil)
			validator := gui.NewQIntValidator2(0, 10000, nil)
			item.SetValidator(validator)
			forms[option.Group].AddRow3(option.Label, item)
			item.ConnectTextChanged(func(val string) {
				if i, err := strconv.Atoi(val); err == nil {
					ac.SetIntOption(key, i)
				}
			})
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

	vbox.AddWidget(buttonBox, 0, core.Qt__AlignBottom)

	if configdialog.Exec() != int(widgets.QDialog__Accepted) {
		log.Info("Canceled option edit")
		undo := ac.Load()
		if !undo {
			log.Error("Error undoing settings changes")
		}
	} else {
		log.Info("Options edited need save")
		if err := ac.Save(); err != nil {
			log.Error("Error saving settings changes")
		}
	}
}
