package main

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/allan-simon/go-singleinstance"
	"github.com/nicksnyder/go-i18n/i18n"
	log "github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var (
	// VERSION of Application
	VERSION = "0.0.1"

	// REVISION Placeholder for git stamping builds
	REVISION = "na"

	// T translate function
	T i18n.TranslateFunc

	neocatalog *NeoCatalog

	config *AppConfig
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// For now only one instance running at a time
	_, err := singleinstance.CreateLockFile("neochess.lock")
	if err != nil {
		os.Exit(1)
	}

	// Translation Load
	qfile := core.NewQFile2(":qml/translate/en-us.all.json")
	if !qfile.Open(core.QIODevice__ReadOnly | core.QIODevice__Text) {
		log.Error("Could not open translation file")
		os.Exit(1)
	}
	qstream := core.NewQTextStream2(qfile)
	qtext := qstream.ReadAll()
	i18n.ParseTranslationFileBytes("en-us.all.json", []byte(qtext))
	T, _ = i18n.Tfunc("en-US")

	core.QCoreApplication_SetOrganizationDomain("neodevelop.net")
	core.QCoreApplication_SetOrganizationName("NeoDevelop")
	core.QCoreApplication_SetApplicationName("NeoChess")
	core.QCoreApplication_SetApplicationVersion(VERSION)
	app := widgets.NewQApplication(len(os.Args), os.Args)
	mainw := widgets.NewQMainWindow(nil, 0)
	if runtime.GOOS == "darwin" {
		mainw.SetUnifiedTitleAndToolBarOnMac(true)
	}
	mainw.SetWindowTitle(core.QCoreApplication_ApplicationName())

	log.Info("Starting Application")

	// Load or initialize settings
	config = initAppConfig(app, mainw)

	var caterr error
	neocatalog, caterr = GetCatalog(config.Datadir)
	if caterr != nil {
		log.Error(caterr)
	}

	mainw.Resize(core.NewQSize2(config.GetIntOption("LastWidth"), config.GetIntOption("LastHeight")))

	statusbar := initStatusBar(mainw)
	mainw.SetStatusBar(statusbar)
	statusbar.AddItem("Screen W:%d x H:%d", config.GetIntOption("LastWidth"), config.GetIntOption("LastHeight"))

	chessdbview := initCDBView(mainw)

	toolbar := initToolBar(mainw, chessdbview)
	mainw.AddToolBar2(toolbar)

	menu := initMenu(mainw, app, chessdbview)
	mainw.SetMenuBar(menu)

	mainw.Show()
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
	widgets.QApplication_Exec()
	config.Save()

}
