package main

import (
	"fmt"
	"strconv"
	"strings"
)

// EInvalidFEN d
const EInvalidFEN = "e:invalid:fen"

// RuneToPiece mapping
var RuneToPiece = map[rune]PieceType{
	'P': WhitePawn, 'N': WhiteKnight, 'B': WhiteBishop, 'R': WhiteRook, 'Q': WhiteQueen, 'K': WhiteKing,
	'p': BlackPawn, 'n': BlackKnight, 'b': BlackBishop, 'r': BlackRook, 'q': BlackQueen, 'k': BlackKing,
}

// PieceToRune mapping
var PieceToRune = map[PieceType]rune{
	WhitePawn: 'P', WhiteKnight: 'N', WhiteBishop: 'B', WhiteRook: 'R', WhiteQueen: 'Q', WhiteKing: 'K',
	BlackPawn: 'p', BlackKnight: 'n', BlackBishop: 'b', BlackRook: 'r', BlackQueen: 'q', BlackKing: 'k',
}

// BasePieceToRune mapping
var BasePieceToRune = map[PieceBaseType]rune{
	Pawn: 'P', Knight: 'N', Bishop: 'B', Rook: 'R', Queen: 'Q', King: 'K',
}

// RuneToFile mapping
var RuneToFile = map[rune]int{'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7}

// RuneToRank mapping
var RuneToRank = map[rune]int{'1': 7, '2': 6, '3': 5, '4': 4, '5': 3, '6': 2, '7': 1, '8': 0}

// FileToString mapping
var FileToString = map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e", 5: "f", 6: "g", 7: "h"}

// RankToString mapping
var RankToString = map[int]string{0: "8", 1: "7", 2: "6", 3: "5", 4: "4", 5: "3", 6: "2", 7: "1"}

// InitFromFen initialize board from Fen string
func (b *BoardType) InitFromFen(fen string) error {
	b.Reset()
	fen = strings.Trim(fen, " ")
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return fmt.Errorf(EInvalidFEN)
	}
	pieces := parts[0]
	turn := parts[1]
	castling := parts[2]
	enPassant := parts[3]
	halfMoves := parts[4]
	fullMoves := parts[5]

	index := 0
	for _, r := range pieces {
		if index > 63 {
			return fmt.Errorf(EInvalidFEN)
		}
		switch r {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			b.addPiece(RuneToPiece[r], index)
			index++
		case '1', '2', '3', '4', '5', '6', '7', '8':
			EmptySquares, _ := strconv.Atoi(string(r))
			for i := 0; i < EmptySquares; i++ {
				b.addPiece(Empty, index)
				index++
			}
		}
	}

	if turn == "w" {
		b.Turn = White
	} else {
		b.Turn = Black
	}

	for _, r := range castling {
		switch r {
		case 'K':
			b.WhiteCastling |= CastleWhiteKingSide
		case 'Q':
			b.WhiteCastling |= CastleWhiteQueenSide
		case 'k':
			b.BlackCastling |= CastleBlackKingSide
		case 'q':
			b.BlackCastling |= CastleBlackQueenSide
		}
	}

	b.HalfMoves, _ = strconv.Atoi(string(halfMoves))
	b.FullMoves, _ = strconv.Atoi(string(fullMoves))

	if enPassant != "-" {
		enPassant = strings.ToLower(enPassant)
		file := RuneToFile[rune(enPassant[0])]
		rank := RuneToRank[rune(enPassant[1])]
		b.EnPassant = CoordsToSquare(rank, file)
	} else {
		b.EnPassant = 0
	}
	return nil
}

// ToFen return FEN representation of board
func (b *BoardType) ToFen() string {
	var pieces, turn, castling, enPassant string

	pieces = ""
	var i = 0
	for rank := 0; rank < 8; rank++ {
		var EmptyCount = 0
		for file := 0; file < 8; file++ {
			if i < 64 {
				for EmptyCount = 0; file < 8 && b.Squares[i] == Empty; file++ {
					EmptyCount++
					i++
				}
				if EmptyCount > 0 {
					pieces += strconv.Itoa(EmptyCount)
				}
				if i < 64 && b.Squares[i] != Empty && file < 8 {
					pieces += string(PieceToRune[b.Squares[i]])
					i++
				}
			}
		}
		if rank < 7 {
			pieces += "/"
		}
	}

	if b.Turn == White {
		turn = "w"
	} else {
		turn = "b"
	}

	castling = ""
	if (b.WhiteCastling & CastleWhiteKingSide) > 0 {
		castling += "K"
	}
	if (b.WhiteCastling & CastleWhiteQueenSide) > 0 {
		castling += "Q"
	}
	if (b.BlackCastling & CastleBlackKingSide) > 0 {
		castling += "k"
	}
	if (b.BlackCastling & CastleBlackQueenSide) > 0 {
		castling += "q"
	}
	if len(castling) == 0 {
		castling = "-"
	}

	if b.EnPassant == 0 {
		enPassant = "-"
	} else {
		rank, file := squareCoords(b.EnPassant)
		enPassant = FileToString[file] + RankToString[rank]
	}

	return fmt.Sprintf("%s %s %s %s %d %d", pieces, turn, castling, enPassant, b.HalfMoves, b.FullMoves)
}
