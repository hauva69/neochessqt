package main

// MoveDisAmbiguationType d
type MoveDisAmbiguationType int

// Constants for MoveDisAmbiguation
const (
	DisAmbiguateNone MoveDisAmbiguationType = iota
	DisAmbiguateByFile
	DisAmbiguateByRank
	DisAmbiguateByBoth
)

var disambiguates = [...]string{
	"None",
	"ByFile",
	"ByRank",
	"ByBoth",
}

func (movedis MoveDisAmbiguationType) String() string {
	return disambiguates[movedis]
}

// PackedMoveType for a compact move
type PackedMoveType uint32

// MoveType Structure
type MoveType struct {
	Move         PackedMoveType
	MoveSuffix   string
	DisAmbiguate MoveDisAmbiguationType
}

func (m MoveType) String() string {
	return m.ToSAN()
}

// Move Bit Locations and Masks
// 33222222222211111111110000000000
// 10987654321098765432109876543210
// 00                               - Encoded Normal Move
// 00                        111111 - From
// 00                  111111       - To
// 00                11             - Special EnPass = 01, Check = 10, Mate = 11, Normal = 00
// 00           11111               - Promotion Piece
// 00      11111                    - Capture Piece
// 00 11111                         - Moving Piece
// 001                              - Castling
// 01                      11111111 - NAG
// 10                      		  1 - Variation Increase
// 10                             0 - Variation Decrease
// 11             11111111111111111  - Comment Begin Length uint16 - Stream count bytes after

const (
	MoveFromBitLocation      = 0  // + SquareMask
	MoveToBitLocation        = 6  // + SquareMask
	MoveSpecialBitLocation   = 12 // + SpecialMask
	MoveCastlingBitLocation  = 29
	MovePromotionBitLocation = 14 // + PieceMask
	MoveCaptureBitLocation   = 19 // + PieceMask
	MovePieceBitLocation     = 24 // + MovePieceMask
	NodeBitLocationMark      = 30 // + SpecialMask

	MoveSquareMask  = 0x3F // 6 Bits
	MovePieceMask   = 0x1F // 5 Bits
	NodeSpecialMask = 0x03 // 2 Bits
	NagMask         = 0xFF
	CommentMask     = 0XFFFF
)

func IncrementVariation() MoveType {
	r := MoveType{}
	// r.Move = PackedMoveType((uint32(0x02) << NodeBitLocationMark) && uint32(0x01))
	return r
}

func DecrementVariation() MoveType {
	r := MoveType{}
	r.Move = PackedMoveType((uint32(0x02) << NodeBitLocationMark))
	return r
}

func AddComment(comment string) []MoveType {
	b := []byte(comment)
	l := len(b)
	lb := PackedMoveType(uint32(l) | (uint32(0x03) << NodeBitLocationMark))
	var r []MoveType
	r = append(r, MoveType{Move: lb})
	for i := 0; i < len(b); i += 4 {
		b2 := uint32(0)
		b3 := uint32(0)
		b4 := uint32(0)
		b1 := uint32(b[i]) << 24
		if i+1 < len(b) {
			b2 = uint32(b[i+1] << 16)
		}
		if i+2 < len(b) {
			b3 = uint32(b[i+2] << 8)
		}
		if i+3 < len(b) {
			b4 = uint32(b[i+3])
		}
		bset := PackedMoveType(b1 | b2 | b3 | b4)
		r = append(r, MoveType{Move: bset})
	}
	return r
}

// NewMove Create a new move
func NewMove(piece PieceType, from, to SquareType, captured PieceType) MoveType {
	m := MoveType{}
	m.Move = PackedMoveType((uint32(from) |
		(uint32(to) << MoveToBitLocation) |
		uint32(captured)<<MoveCaptureBitLocation |
		uint32(piece)<<MovePieceBitLocation) &^ (NodeSpecialMask << MoveSpecialBitLocation))
	return m
}

// NewEnPassantMove a enpassant move
func NewEnPassantMove(piece PieceType, from, to SquareType, captured PieceType) MoveType {
	m := NewMove(piece, from, to, captured)
	m.Move |= PackedMoveType(uint32(0x1) << MoveSpecialBitLocation)
	return m
}

// NewPromotionMove a promotion move
func NewPromotionMove(piece PieceType, from, to SquareType, captured PieceType, promotionTo PieceType) MoveType {
	m := NewMove(piece, from, to, captured)
	m.Move |= PackedMoveType(uint32(promotionTo) << MovePromotionBitLocation)
	return m
}

// NewCastlingMove a castle move
func NewCastlingMove(piece PieceType, from, to SquareType) MoveType {
	m := NewMove(piece, from, to, Empty)
	m.Move |= PackedMoveType(0x1 << MoveCastlingBitLocation)
	return m
}

func (m MoveType) from() SquareType {
	return SquareType(m.Move & MoveSquareMask)
}

func (m MoveType) to() SquareType {
	return SquareType((m.Move >> MoveToBitLocation) & MoveSquareMask)
}

// ToRune will be from to
func (m MoveType) ToRune() string {
	return m.from().ToRune() + m.to().ToRune()
}

// ToSAN Create SAN notation from Move
func (m MoveType) ToSAN() string {
	san := ""
	country := "English"
	if !m.isCastlingMove() {
		if m.piece().Kind() != Pawn {
			san += m.piece().Kind().ToRuneCountry(country)
		} else {

		}
		if m.DisAmbiguate == DisAmbiguateByFile {
			san += m.from().FileToRune()
		}
		if m.DisAmbiguate == DisAmbiguateByRank {
			san += m.from().RankToRune()
		}
		if m.DisAmbiguate == DisAmbiguateByBoth {
			san += m.from().ToRune()
		}
		if m.captured().Kind() != Empty {
			if m.piece().Kind() == Pawn {
				if m.DisAmbiguate != DisAmbiguateByFile {
					san += m.from().FileToRune()
				}
			}
			san += "x"
		}
		san += m.to().ToRune()
		if m.isPromotionMove() {
			san += "="
			san += m.getPromotionTo().ToRune()
		}
		san += m.MoveSuffix
		return san
	}
	if m.to() == SquareType(WhiteKingSideKingCastleSquare) ||
		m.to() == SquareType(BlackKingSideKingCastleSquare) {
		san = "O-O"
		return san
	}
	if m.to() == SquareType(WhiteQueenSideKingCastleSquare) ||
		m.to() == SquareType(BlackQueenSideKingCastleSquare) {
		san = "O-O-O"
		return san
	}
	return san
}

func (m MoveType) piece() PieceType {
	return PieceType((m.Move >> MovePieceBitLocation) & MovePieceMask)
}

func (m MoveType) color() ColorType {
	return PieceType((m.Move >> MovePieceBitLocation) & MovePieceMask).Color()
}

func (m MoveType) captured() PieceType {
	return PieceType((m.Move >> MoveCaptureBitLocation) & MovePieceMask)
}

func (m MoveType) captures() bool {
	return PieceType((m.Move>>MoveCaptureBitLocation)&MovePieceMask) > 0
}

func (m MoveType) isEnPassant() bool {
	return (uint32(m.Move)>>MoveSpecialBitLocation)&0x01 > 0
}

func (m MoveType) causesCheck() bool {
	return (uint32(m.Move)>>MoveSpecialBitLocation)&0x10 == 2
}

func (m MoveType) causesMate() bool {
	return false // (uint32(m.Move)>>MoveMateBitLocation)&0x11 == 3
}

func (m MoveType) isPromotionMove() bool {
	return (uint32(m.Move)>>MovePromotionBitLocation)&MovePieceMask > 0
}

func (m MoveType) getPromotionTo() PieceType {
	return PieceType((m.Move >> MovePromotionBitLocation) & MovePieceMask)
}

func (m MoveType) isCastlingMove() bool {
	return (uint32(m.Move)>>MoveCastlingBitLocation)&0x1 > 0
}
