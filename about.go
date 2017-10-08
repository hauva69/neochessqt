package main

import (
	"bytes"
	"text/template"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var aboutdialog *widgets.QDialog

// DisplayAbout Dialog Box
func DisplayAbout(w *widgets.QMainWindow) {
	if aboutdialog == nil {
		aboutdialog = widgets.NewQDialog(w, core.Qt__Dialog)

		vbox := widgets.NewQVBoxLayout()

		gridbox := widgets.NewQGridLayout(nil)

		logo := gui.NewQPixmap5(":qml/assets/NeoLogo.png", "png", core.Qt__AutoColor)
		logolabel := widgets.NewQLabel(nil, core.Qt__Widget)
		logolabel.SetPixmap(logo)

		gridbox.AddWidget(logolabel, 0, 0, core.Qt__AlignCenter)

		var appdata = map[string]string{
			"appname":    core.QCoreApplication_ApplicationName(),
			"appversion": core.QCoreApplication_ApplicationVersion(),
			"orglinkurl": "https://" + core.QCoreApplication_OrganizationDomain(),
			"orgname":    core.QCoreApplication_OrganizationName(),
		}

		//applabel := widgets.NewQLabel2(appname, nil, core.Qt__Widget)
		//applabel.SetFont(titlefont)
		//versionlabel := widgets.NewQLabel2(appversion, nil, core.Qt__Widget)
		//versionlabel.SetFont(versionfont)

		//orglink := widgets.NewQLabel2("<a style=\"color:lightblue;\" href=\""+orglinkurl+"\">"+orgname+"</a>", nil, core.Qt__Widget)
		//orglink.SetFont(versionfont)
		//orglink.SetTextInteractionFlags(core.Qt__TextBrowserInteraction)
		//orglink.ConnectLinkActivated(func(link string) {
		//	qurl := core.NewQUrl3(orglinkurl, core.QUrl__StrictMode)
		//	gui.QDesktopServices_OpenUrl(qurl)
		//})

		appinfo := widgets.NewQTextBrowser(nil)
		appinfo.SetReadOnly(true)

		t := template.Must(template.New("App").Parse(`<h2 style="text-align:center;">{{.appname}} Database</h2>
			<h3 style="text-align:center;">Version: <strong>{{.appversion}}</strong></h3>
			<h3 style="text-align:center;"><a style="color:lightblue;" href="{{.orglinkurl}}">{{.orgname}}</h3>`))
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, appdata); err != nil {
		}
		appinfo.SetHtml(tpl.String())
		appinfo.SetFixedHeight(128)
		appinfo.SetOpenExternalLinks(true)

		gridbox.AddWidget(appinfo, 0, 1, core.Qt__AlignCenter)
		gridgroupbox := widgets.NewQGroupBox(nil)
		gridgroupbox.SetLayout(gridbox)

		vbox.AddWidget(gridgroupbox, 0, core.Qt__AlignTop)

		notes := widgets.NewQTextEdit(nil)
		notes.SetReadOnly(true)
		notes.SetHtml(`
			<h2>Credits</h2>
			<ul>
			<li>First Credit</li>
			<li>Second Credit</li>
			</ul>
			<h2>Release Notes</h2>
			<ul>
			<li>Release 1</li>
			<li>Release 2</li>
			</ul>
		`)

		vbox.AddWidget(notes, 0, core.Qt__AlignBottom)

		buttonBox := widgets.NewQDialogButtonBox(nil)
		closeButton := widgets.NewQPushButton2("Close", nil)
		buttonBox.AddButton(closeButton, widgets.QDialogButtonBox__ResetRole)

		closeButton.ConnectPressed(func() {
			aboutdialog.Hide()
		})

		vbox.AddWidget(buttonBox, 0, core.Qt__AlignRight)

		aboutdialog.SetLayout(vbox)
		aboutdialog.Show()
	} else {
		aboutdialog.Show()
	}
}
