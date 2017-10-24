package main

import (
	"bufio"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/rashwell/neochesslib"
)

var (
	tagregex *regexp.Regexp
)

func init() {
	tagregex, _ = regexp.Compile(`\[\s*(?P<tagName>\w+)\s*"(?P<tagValue>[^"]*)"\s*\]\s*`)
}

func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func scanmoves(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	peekr, _ := utf8.DecodeRune(data[start:])
	if peekr == '(' || peekr == ')' {
		return start + 1, data[start:], nil
	}
	// Scan until special token or end of move space
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if peekr == '{' {
			if r == '}' {
				return i + width, data[start:i], nil
			}
		} else {
			if isDigit(peekr) {
				if isSpace(r) || r == '.' {
					return i + width, data[start:i], nil
				}
			} else {
				if isSpace(r) {
					return i + width, data[start:i], nil
				}
			}
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

// ParseGameString into Game
func ParseGameString(gbytes []byte, gameid int, debug bool) *neochesslib.Game {
	gs := string(gbytes)
	g := neochesslib.NewGame()
	if debug {
		//	fmt.Printf("Text size: %d\n", len(gs))
		//	fmt.Printf("Game String:\n%s\n", gs)
	}
	g.ID = uint32(gameid)
	scanner := bufio.NewScanner(strings.NewReader(gs))
	var textsize int
	var err error
	for scanner.Scan() {
		line := scanner.Text()
		if debug {
			//	fmt.Printf("Line: %s\n", line)
		}
		textsize += len(line)
		g, err = ParseLine(line, g)
		if err != nil {
			log.Fatal(err)
		}
	}
	if debug {
		//	fmt.Printf("Game Scanned\n")
		//	fmt.Printf("Parsing Moves\n")
	}
	movescanner := bufio.NewScanner(strings.NewReader(g.Importtext))
	movescanner.Split(scanmoves)
	for movescanner.Scan() {
		token := movescanner.Text()
		if debug {
			// fmt.Printf("Move Token: %s\n", token)
		}
		r, _ := utf8.DecodeRuneInString(token)
		if !isDigit(r) {
			if move, exists := g.IsTokenInMoves(token); exists {
				g.AppendMove(move)
			}
		}
	}
	return g
}

// ParseLine directs line based on status of game parsing
func ParseLine(line string, g *neochesslib.Game) (*neochesslib.Game, error) {
	trimline := strings.Trim(line, " ")
	if trimline == "" {
		return g, nil
	}
	if trimline[0] == '%' {
		return g, nil
	}
	if trimline[0] == '[' {
		g, err := ParseTag(trimline, g)
		if err != nil {
			return g, err
		}
		return g, nil
	}
	g.Importtext += " " + line
	return g, nil
}

// ParseTag parses a chess tag into game
func ParseTag(line string, g *neochesslib.Game) (*neochesslib.Game, error) {
	results := tagregex.FindStringSubmatch(line)
	if len(results) == 3 {
		tagname := results[1]
		tagvalue := results[2]
		if !neochesslib.IsTagSupported(tagname) {
			return g, nil
		}
		if !(len(tagvalue) > 0) {
			return g, nil
		}
		g.SetTag(tagname, tagvalue)
		return g, nil
	}
	return g, nil
}
