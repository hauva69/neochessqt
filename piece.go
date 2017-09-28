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
