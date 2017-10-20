package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// GameAnalysisDock comment
type GameAnalysisDock struct {
	widgets.QDockWidget
	analysis      *widgets.QTextEdit
	enginecomm    io.WriteCloser
	analysisgame  *Game
	analysisboard *BoardType
	running       bool
	visible       bool
}

func initGameAnalysisDock(w *widgets.QMainWindow) *GameAnalysisDock {
	this := NewGameAnalysisDock("Game Analysis", w, core.Qt__Widget)
	this.analysis = widgets.NewQTextEdit(nil)
	this.analysis.SetReadOnly(true)
	this.running = false
	this.visible = false
	layout := widgets.NewQVBoxLayout()
	this.analysis.SetLayout(layout)
	this.SetAllowedAreas(core.Qt__RightDockWidgetArea)

	this.SetWidget(this.analysis)
	return this
}

func (ga *GameAnalysisDock) enginerun(engine string) (<-chan string, io.WriteCloser) {
	cmd := exec.Command(engine)
	cmdReader, _ := cmd.StdoutPipe()
	cmdWriter, _ := cmd.StdinPipe()

	scanner := bufio.NewScanner(cmdReader)
	stream := make(chan string)

	go func() {
		defer close(stream)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "info") && !strings.Contains(line, "currmove") {
				parts := strings.Split(line, " pv ")
				if len(parts) >= 2 {
					stream <- parts[1]
				}
			}
		}
	}()

	go func() {
		if err := cmd.Run(); err != nil {
			log.Error("Error Running Engine")
		}
	}()

	return stream, cmdWriter
}

// ToggleEngine comment
func (ga *GameAnalysisDock) ToggleEngine(engine int, fen string) {
	if !ga.running {
		log.Info("Attempting to Start Chess Engine")
		enginekey := fmt.Sprintf("Engine%d", engine)
		if config.GetStrOption(enginekey) == "" {
			ga.analysis.SetText("Engine not configured")
			ga.running = false
			if !ga.visible {
				ga.Show()
				ga.visible = true
			}
		} else {
			if !ga.visible {
				ga.Show()
				ga.visible = true
			}
			ga.analysis.Clear()
			ga.running = true
			stream, enginecomm := ga.enginerun(config.GetStrOption(enginekey))
			ga.enginecomm = enginecomm
			go func() {
				for line := range stream {
					ga.analysis.Clear()
					enginemoves := strings.Split(line, " ")
					board := NewBoard()
					board.InitFromFen(fen)
					boardmoves := board.GenerateLegalMoves()
					pvline := ""
					if board.Turn == Black {
						pvline = "... "
					}
					mn := board.FullMoves
					for _, enginemove := range enginemoves {
						for _, move := range boardmoves {
							if move.ToRune() == enginemove {
								if board.Turn == White {
									pvline += "<span class='movenumber'>" + strconv.Itoa(mn) + ". </span>"
									mn++
								}
								pvline += "<span class='move'>" + move.ToSAN() + "</span> "
								board.MakeMove(move, true)
								break
							}
						}
						boardmoves = board.GenerateLegalMoves()
					}
					ga.analysis.SetHtml("<style>" + config.PGNStyle + "</style>" + pvline)
				}
			}()
			ga.EngineCommand("uci")
			ga.EngineCommand("ucinewgame")
			ga.EngineCommand("position fen " + fen)
			ga.EngineCommand("go infinite")
		}
	} else {
		log.Info("Stopping Engine")
		ga.EngineCommand("quit")
		ga.analysis.Clear()
		ga.running = false
		ga.Hide()
		ga.visible = false
	}
}

// EngineCommand comment
func (ga *GameAnalysisDock) EngineCommand(cmd string) {
	if ga.running {
		io.WriteString(ga.enginecomm, cmd+"\n")
	}
}
