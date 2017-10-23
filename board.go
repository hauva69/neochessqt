package main

import (
	"fmt"
	"strings"
)

//Castling permissions
const (
	CastleWhiteKingSide  uint = 1
	CastleWhiteQueenSide uint = 2
	CastleBlackKingSide  uint = 1
	CastleBlackQueenSide uint = 2
)

//castling squares
const (
	WhiteKingSideKingCastleSquare  = G1
	WhiteQueenSideKingCastleSquare = C1
	BlackKingSideKingCastleSquare  = G8
	BlackQueenSideKingCastleSquare = C8
	WhiteKingSideRookCastleSquare  = F1
	WhiteQueenSideRookCastleSquare = C1
	BlackKingSideRookCastleSquare  = F8
	BlackQueenSideRookCastleSquare = C8
)

// BoardType struct
type BoardType struct {
	WhitePieces          BitBoardType
	BlackPieces          BitBoardType
	Whites               [7]BitBoardType
	Blacks               [7]BitBoardType
	WhiteCastling        uint
	WhiteCastlingHistory []uint
	BlackCastling        uint
	BlackCastlingHistory []uint
	Occupied             BitBoardType
	EnPassant            SquareType
	EnPassantHistory     []SquareType
	Turn                 ColorType
	Check                bool
	Checkmate            bool
	Stalement            bool
	Squares              [64]PieceType
	HalfMoves            int
	FullMoves            int
	MoveHistory          []MoveType
}

// NewBoard Create new board
func NewBoard() *BoardType {
	b := BoardType{}
	b.Reset()
	return &b
}

// Reset Board
func (b *BoardType) Reset() {
	for _, kind := range [6]PieceType{Pawn, Knight, Bishop, Rook, Queen, King} {
		b.Whites[kind] = BitBoardType(0)
		b.Blacks[kind] = BitBoardType(0)
	}
	b.WhitePieces = BitBoardType(0)
	b.BlackPieces = BitBoardType(0)
	b.Occupied = BitBoardType(0)
	b.EnPassant = 0
	b.Turn = White
	b.Check = false
	b.Checkmate = false
	b.Stalement = false
	b.Squares = [64]PieceType{}
	b.WhiteCastling = CastleWhiteKingSide & CastleWhiteQueenSide
	b.WhiteCastlingHistory = []uint(nil)
	b.BlackCastling = CastleBlackKingSide & CastleBlackQueenSide
	b.BlackCastlingHistory = []uint(nil)
	b.HalfMoves = 0
	b.FullMoves = 0
	b.MoveHistory = []MoveType(nil)
}

//GetColors our color and opponent's
func (b *BoardType) GetColors() (ColorType, ColorType) {
	if b.Turn == White {
		return White, Black
	}
	return Black, White
}

func (b *BoardType) oppositeturn() ColorType {
	if b.Turn == White {
		return Black
	}
	return White
}

func (b *BoardType) addPiece(piece PieceType, index int) {
	b.Squares[index] = piece
	if piece == Empty {
		return
	}
	bit := BitBoardType(0x1 << uint(index))
	kind := piece.Kind()
	switch piece.Color() {
	case White:
		b.Whites[kind] |= bit
		b.WhitePieces |= bit
	case Black:
		b.Blacks[kind] |= bit
		b.BlackPieces |= bit
	}
	b.Occupied |= bit
}

func (b *BoardType) switchTurn() {
	if b.Turn == White {
		b.Turn = Black
		return
	}
	b.Turn = White
}

// For Debugging and repl

// PrintBoard main display board
func (b *BoardType) PrintBoard(title string) {
	mboard := b.sPrintBoard(title)
	for _, line := range mboard {
		fmt.Printf("%s\n", line)
	}
}

// PrintDashboard d
func PrintDashboard(b *BoardType) {
	occupied := sPrintBitboard(b.Occupied, "Occupied")
	whitep := sPrintBitboard(b.WhitePieces, "All White")
	blackp := sPrintBitboard(b.BlackPieces, "All Black")
	mboard := b.sPrintBoard("Main Board")
	for index, line := range occupied {
		fmt.Printf("%s\t%s\t%s\t%s\n", line, whitep[index], blackp[index], mboard[index])
	}
	wp := sPrintBitboard(b.Whites[Pawn], "White Pawns")
	wn := sPrintBitboard(b.Whites[Knight], "White Knights")
	wb := sPrintBitboard(b.Whites[Bishop], "White Bishops")
	wr := sPrintBitboard(b.Whites[Rook], "White Rooks")
	for index, line := range wp {
		fmt.Printf("%s\t%s\t%s\t%s\n", line, wn[index], wb[index], wr[index])
	}
	wq := sPrintBitboard(b.Whites[Queen], "White Queens")
	wk := sPrintBitboard(b.Whites[King], "White King")
	bq := sPrintBitboard(b.Blacks[Queen], "Black Queens")
	bk := sPrintBitboard(b.Blacks[King], "Black King")
	for index, line := range wq {
		fmt.Printf("%s\t%s\t%s\t%s\n", line, wk[index], bq[index], bk[index])
	}
	bp := sPrintBitboard(b.Blacks[Pawn], "Black Pawns")
	bn := sPrintBitboard(b.Blacks[Knight], "Black Knights")
	bb := sPrintBitboard(b.Blacks[Bishop], "Black Bishops")
	br := sPrintBitboard(b.Blacks[Rook], "Black Rooks")
	for index, line := range bp {
		fmt.Printf("%s\t%s\t%s\t%s\n", line, bn[index], bb[index], br[index])
	}
}

// PrintCalcboard d
func PrintCalcboard(b *BoardType, bb1 BitBoardType, bb1title string, bb2 BitBoardType, bb2title string, bb3 BitBoardType, bb3title string) {
	bb1lines := sPrintBitboard(bb1, bb1title)
	bb2lines := sPrintBitboard(bb2, bb2title)
	bb3lines := sPrintBitboard(bb3, bb3title)
	for index, line := range bb1lines {
		fmt.Printf("%s\t%s\t%s\n", line, bb2lines[index], bb3lines[index])
	}
}

func (b *BoardType) sPrintBoard(title string) []string {
	s := []string(nil)
	title = strings.Trim(title, " ")
	if len(title) > 22 {
		title = title[:22]
	}
	ls := (24 - len(title)) / 2
	rs := 24 - len(title) - ls

	title = strings.Repeat("_", ls) + title + strings.Repeat("_", rs)
	s = append(s, fmt.Sprintf("%s", title))
	for rank := 7; rank >= 0; rank-- {
		var line string
		for file := 7; file >= 0; file-- {
			var squareIdx = uint(rank*8 + file)
			if b.Squares[squareIdx] == Empty {
				if squareIdx != 0 && squareIdx == uint(b.EnPassant.rank()*8+b.EnPassant.file()) {
					line += fmt.Sprintf(" z ")
				} else {
					line += fmt.Sprintf(" - ")
				}
			} else {
				line += fmt.Sprintf(" %c ", PieceToRune[b.Squares[squareIdx]])
			}
		}
		s = append(s, line)
	}
	return s
}
