package main

import (
	"bufio"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func parselogitemn(logitem string) (results map[string]string) {
	r, _ := regexp.Compile("time=\"(?P<time>[^\"]*)\" level=(?P<level>[^ ]*) msg=\"(?P<msg>[^\"]*)\"")
	match := r.FindStringSubmatch(logitem)
	results = make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			results[name] = match[i]
		}
	}
	return
}

func init() {
	err := os.Remove("neochess.log")
	file, err := os.OpenFile("neochess.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
		// log.SetOutput(os.Stdout)
	}
	log.SetLevel(log.InfoLevel)
}

var logdialog *widgets.QDialog
var logedit *widgets.QTextEdit
var logwatch *core.QFileSystemWatcher

func updatelog() {
	levelcolors := map[string]string{
		"info":    "style='color:#6699ff;'",
		"warning": "style='color:#cccc00;'",
		"error":   "style='color:#cc2900;'",
	}

	levelstring := map[string]string{
		"info":    "INFO",
		"warning": "WARN",
		"error":   "FATA",
	}
	loghtml := "<pre>"
	file, err := os.Open("neochess.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		results := parselogitemn(scanner.Text())
		layout := "2006-01-02T15:04:05-07:00"
		stamp, _ := time.Parse(layout, results["time"])
		loghtml = loghtml + "<span " + levelcolors[results["level"]] + ">" + levelstring[results["level"]] + "</span> | "
		loghtml = loghtml + "<span>" + stamp.Format("3:04PM") + "</span> | "
		loghtml = loghtml + "<span>" + results["msg"] + "</span><br/>"
	}
	loghtml = loghtml + "</pre>"
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	logedit.SetHtml(loghtml)
	sb := logedit.VerticalScrollBar()
	sb.SetValue(sb.Maximum())
}

// DisplayLog dialog
func DisplayLog(w *widgets.QMainWindow) {
	if logdialog == nil {
		logdialog = widgets.NewQDialog(w, core.Qt__Dialog)
		vbox := widgets.NewQVBoxLayout()

		logedit = widgets.NewQTextEdit(nil)
		logedit.SetFixedWidth(400)
		logedit.SetReadOnly(true)

		vbox.AddWidget(logedit, 0, core.Qt__AlignTop)

		buttonBox := widgets.NewQDialogButtonBox(nil)
		closeButton := widgets.NewQPushButton2("Close", nil)
		buttonBox.AddButton(closeButton, widgets.QDialogButtonBox__ResetRole)

		closeButton.ConnectPressed(func() {
			logdialog.Hide()
		})

		vbox.AddWidget(buttonBox, 0, core.Qt__AlignBottom)
		logdialog.SetLayout(vbox)

		logwatch = core.NewQFileSystemWatcher(nil)
		logwatch.AddPath("neochess.log")
		logwatch.ConnectFileChanged(func(p string) {
			updatelog()
		})
		logdialog.Show()
		updatelog()
	} else {
		logdialog.Show()
	}
}
