package main

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

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
	}
	logwatch = core.NewQFileSystemWatcher(nil)
	logwatch.AddPath("neochess.log")
	logwatch.ConnectFileChanged(func(p string) {
		loghtml := ""
		file, err := os.Open("neochess.log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			loghtml = loghtml + "<p>" + scanner.Text() + "</p>"
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		logedit.SetHtml(loghtml)
	})
	logdialog.Show()
}
