package main

import "strings"

// Color Masks for piece color detect
const (
	WhiteMask = 0x8
	BlackMask = 0x10
)

// PieceBaseType enum
type PieceBaseType int

// Constants for PieceBaseType
const (
	Empty = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

var basepieces = [...]string{
	"Empty",
	"Pawn",
	"Knight",
	"Bishop",
	"Rook",
	"Queen",
	"King",
}

func (p PieceBaseType) String() string {
	return basepieces[p]
}

// ToRune Base Piece to Rune
func (p PieceBaseType) ToRune() string {
	return string(BasePieceToRune[p])
}

// PieceType enum
type PieceType int

// Constants for PieceType
const (
	NoPiece     PieceType = 0
	WhitePawn   PieceType = PieceType(White | Pawn)
	WhiteKnight PieceType = PieceType(White | Knight)
	WhiteBishop PieceType = PieceType(White | Bishop)
	WhiteRook   PieceType = PieceType(White | Rook)
	WhiteQueen  PieceType = PieceType(White | Queen)
	WhiteKing   PieceType = PieceType(White | King)
	BlackPawn   PieceType = PieceType(Black | Pawn)
	BlackKnight PieceType = PieceType(Black | Knight)
	BlackBishop PieceType = PieceType(Black | Bishop)
	BlackRook   PieceType = PieceType(Black | Rook)
	BlackQueen  PieceType = PieceType(Black | Queen)
	BlackKing   PieceType = PieceType(Black | King)
)

// Kind get the base Piece of a colored piece
func (p PieceType) Kind() PieceBaseType {
	return PieceBaseType(p & (0x7))
}

// Color get the base Color of Piece
func (p PieceType) Color() ColorType {
	if p&WhiteMask > 0 {
		return White
	}
	if p&BlackMask > 0 {
		return Black
	}
	return NoColor
}

// ToRune d
func (p PieceType) ToRune() string {
	return string(PieceToRune[p])
}

func (p PieceType) FileName() string {
	strpcolor := strings.ToLower(p.Color().String())
	strpkind := strings.ToLower(p.Kind().String())
	return ":qml/assets/" + strpcolor + strpkind + ".png"
}

func (p PieceType) String() string {
	return p.Color().String() + " " + p.Kind().String()
}

func (pb PieceBaseType) ToRuneCountry(country string) string {
	return string(PieceToRuneCountry[pb][country])
}

var PieceToRuneCountry = map[PieceBaseType]map[string]rune{
	Pawn: {
		"Czech":      'P',
		"Danish":     'B',
		"Dutch":      'O',
		"English":    'P',
		"Estonian":   'P',
		"Finnish":    'P',
		"French":     'P',
		"German":     'B',
		"Hungarian":  'G',
		"Icelandic":  'P',
		"Italian":    'P',
		"Norwegian":  'B',
		"Polish":     'P',
		"Portuguese": 'P',
		"Romanian":   'P',
		"Spanish":    'P',
		"Swedish":    'B',
	},
	Knight: {
		"Czech":      'J',
		"Danish":     'S',
		"Dutch":      'P',
		"English":    'N',
		"Estonian":   'R',
		"Finnish":    'R',
		"French":     'C',
		"German":     'S',
		"Hungarian":  'H',
		"Icelandic":  'R',
		"Italian":    'C',
		"Norwegian":  'S',
		"Polish":     'S',
		"Portuguese": 'C',
		"Romanian":   'C',
		"Spanish":    'C',
		"Swedish":    'S',
	},
	Bishop: {
		"Czech":      'S',
		"Danish":     'L',
		"Dutch":      'L',
		"English":    'B',
		"Estonian":   'O',
		"Finnish":    'L',
		"French":     'F',
		"German":     'L',
		"Hungarian":  'F',
		"Icelandic":  'B',
		"Italian":    'A',
		"Norwegian":  'L',
		"Polish":     'G',
		"Portuguese": 'B',
		"Romanian":   'N',
		"Spanish":    'A',
		"Swedish":    'L',
	},
	Rook: {
		"Czech":      'V',
		"Danish":     'T',
		"Dutch":      'T',
		"English":    'R',
		"Estonian":   'V',
		"Finnish":    'T',
		"French":     'T',
		"German":     'T',
		"Hungarian":  'B',
		"Icelandic":  'H',
		"Italian":    'T',
		"Norwegian":  'T',
		"Polish":     'W',
		"Portuguese": 'T',
		"Romanian":   'T',
		"Spanish":    'T',
		"Swedish":    'T',
	},
	Queen: {
		"Czech":      'D',
		"Danish":     'D',
		"Dutch":      'D',
		"English":    'Q',
		"Estonian":   'L',
		"Finnish":    'D',
		"French":     'D',
		"German":     'D',
		"Hungarian":  'V',
		"Icelandic":  'D',
		"Italian":    'D',
		"Norwegian":  'D',
		"Polish":     'H',
		"Portuguese": 'D',
		"Romanian":   'D',
		"Spanish":    'D',
		"Swedish":    'D',
	},
	King: {
		"Czech":      'K',
		"Danish":     'K',
		"Dutch":      'K',
		"English":    'K',
		"Estonian":   'K',
		"Finnish":    'K',
		"French":     'R',
		"German":     'K',
		"Hungarian":  'K',
		"Icelandic":  'K',
		"Italian":    'R',
		"Norwegian":  'K',
		"Polish":     'K',
		"Portuguese": 'R',
		"Romanian":   'R',
		"Spanish":    'R',
		"Swedish":    'K',
	},
}
