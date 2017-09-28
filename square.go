package main

// SquareType for Square
type SquareType uint8

//get rank from square
func (s SquareType) rank() int {
	return int(s / 8)
}

//get file from square
func (s SquareType) file() int {
	return int(s % 8)
}

// ToRune Convert Square coord to runes
func (s SquareType) ToRune() string {
	rank, file := squareCoords(s)
	return FileToString[file] + RankToString[rank]
}

// FileToRune Convert Square coord File to rune
func (s SquareType) FileToRune() string {
	_, file := squareCoords(s)
	return FileToString[file]
}

// RankToRune Convert Square coord File to rune
func (s SquareType) RankToRune() string {
	rank, _ := squareCoords(s)
	return RankToString[rank]
}

// Easy Square Indexes
const (
	A8 = iota
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A1
	B1
	C1
	D1
	E1
	F1
	G1
	H1
)

//IsCoordsOutofBoard whether coords are out of board or not
func IsCoordsOutofBoard(rank, file int) bool {
	return rank < 0 || rank > 7 || file < 0 || file > 7
}

//CoordsToSquare Convert rank, file coordinates to a 64 based square
func CoordsToSquare(rank, file int) SquareType {
	return SquareType(rank*8 + file)
}

//CoordsToIndex Convert rank, file coordinates to a 64 based index
func CoordsToIndex(rank, file int) uint {
	return uint(rank*8 + file)
}

//get rank, file from Square
func squareCoords(sq SquareType) (int, int) {
	return int(sq) / 8, int(sq) % 8
}
