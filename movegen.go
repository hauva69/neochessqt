package main

import "fmt"

//GenerateLegalMoves all legal moves
func (b *BoardType) GenerateLegalMoves() []MoveType {
	var ours [7]BitBoardType
	var oursAll BitBoardType
	var theirsAll BitBoardType
	if b.Turn == White {
		ours = b.Whites
		oursAll = b.WhitePieces
		theirsAll = b.BlackPieces
	} else {
		ours = b.Blacks
		oursAll = b.BlackPieces
		theirsAll = b.WhitePieces
	}

	pseudoMoves := []MoveType(nil)
	legalMoves := []MoveType(nil)

	pseudoMoves = append(pseudoMoves, b.genFromMoves(PieceType(b.Turn|King), ours[King], oursAll, KingAttacksFrom[:])...)
	pseudoMoves = append(pseudoMoves, b.genFromMoves(PieceType(b.Turn|Knight), ours[Knight], oursAll, KnightAttacksFrom[:])...)
	pseudoMoves = append(pseudoMoves, b.genRayMoves(PieceType(b.Turn|Bishop), ours[Bishop], oursAll, theirsAll, BishopDirections[:])...)
	pseudoMoves = append(pseudoMoves, b.genRayMoves(PieceType(b.Turn|Rook), ours[Rook], oursAll, theirsAll, RookDirections[:])...)
	pseudoMoves = append(pseudoMoves, b.genRayMoves(PieceType(b.Turn|Queen), ours[Queen], oursAll, theirsAll, BishopDirections[:])...)
	pseudoMoves = append(pseudoMoves, b.genRayMoves(PieceType(b.Turn|Queen), ours[Queen], oursAll, theirsAll, RookDirections[:])...)

	pseudoMoves = append(pseudoMoves, b.genPawnTwoSteps(PieceType(b.Turn|Pawn))...)
	pseudoMoves = append(pseudoMoves, b.genPawnAttacks(b.Turn, PieceType(b.Turn|Pawn), theirsAll)...)
	pseudoMoves = append(pseudoMoves, b.genCastling(PieceType(b.Turn|King))...)
	pseudoMoves = append(pseudoMoves, b.genPawnOneStep(PieceType(b.Turn|Pawn))...)

	// TargetMoveList := NewMoveList()
	//  fmt.Printf("Pseudo Moves: %v\n", pseudoMoves)
	ambigMap := map[string][]int{}
	ambigFileMap := map[string][]int{}
	ambigRankMap := map[string][]int{}
	i := 0
	for _, move := range pseudoMoves {
		canmove, causescheck, causesmate := b.CanMove(move)
		switch {
		case causesmate:
			move.MoveSuffix += "#"
			//			move.Move |= PackedMoveType(0x1 << MoveMateBitLocation)
		case causescheck && !causesmate:
			move.MoveSuffix += "+"
			//			move.Move |= PackedMoveType(0x1 << MoveCheckBitLocation)
		case canmove:
			move.DisAmbiguate = DisAmbiguateNone
			if !move.isPromotionMove() {
				ambikey := ""
				if move.piece().Kind() == Pawn {
					ambikey += "P"
				} else {
					ambikey += move.piece().Kind().ToRune()
				}
				ambikey += move.to().ToRune()
				ambigMap[ambikey] = append(ambigMap[ambikey], i)
				ambigFileMap[ambikey+FileToString[move.from().file()]] = append(ambigFileMap[ambikey+FileToString[move.from().file()]], i)
				ambigRankMap[ambikey+RankToString[move.from().rank()]] = append(ambigRankMap[ambikey+RankToString[move.from().rank()]], i)
			}
		}
		if canmove {
			legalMoves = append(legalMoves, move)
			i++
		}
	}
	for _, sml := range ambigMap {
		if len(sml) > 1 {
			for _, item := range sml {
				mv := legalMoves[item]
				switch {
				case len(ambigFileMap[mv.piece().Kind().ToRune()+mv.to().ToRune()+FileToString[mv.from().file()]]) == 1:
					legalMoves[item].DisAmbiguate = DisAmbiguateByFile
				case len(ambigRankMap[mv.piece().Kind().ToRune()+mv.to().ToRune()+RankToString[mv.from().rank()]]) == 1:

				}
				if len(ambigFileMap[mv.piece().Kind().ToRune()+mv.to().ToRune()+FileToString[mv.from().file()]]) == 1 {
					legalMoves[item].DisAmbiguate = DisAmbiguateByFile
				} else {
					if len(ambigRankMap[mv.piece().Kind().ToRune()+mv.to().ToRune()+RankToString[mv.from().rank()]]) == 1 {
						legalMoves[item].DisAmbiguate = DisAmbiguateByRank
					} else {
						legalMoves[item].DisAmbiguate = DisAmbiguateByBoth
					}
				}
			}
		}
	}
	// fmt.Printf("Legal Moves: %v\n", legalMoves)
	return legalMoves
}

// Counters for movegenerator
func Counters(ml []MoveType) (uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	nodes := len(ml)
	var capcount uint64
	var empassantcount uint64
	var castlecount uint64
	var promotioncount uint64
	var checkcount uint64
	var checkmatecount uint64

	for _, move := range ml {
		if move.isEnPassant() {
			empassantcount++
		}
		if move.isCastlingMove() {
			castlecount++
		}
		if move.isPromotionMove() {
			promotioncount++
		}
		if move.captures() {
			capcount++
		}
		if move.causesCheck() {
			checkcount++
		}
		if move.causesMate() {
			checkcount++
			checkmatecount++
		}
	}
	return uint64(nodes), capcount, empassantcount, castlecount, promotioncount, checkcount, checkmatecount
}

// CanMove Checks whether the given move is possible or not
func (b *BoardType) CanMove(m MoveType) (legal bool, causesCheck bool, causesMate bool) {
	causesCheck = false
	causesMate = false
	legal = false

	if m.isCastlingMove() {
		if b.IsCheck(b.Turn) {
			return
		}
		var testfile int
		if m.to().file() == 2 {
			testfile = 3
		} else {
			testfile = 5
		}
		m1 := NewMove(m.piece(), m.from(), CoordsToSquare(m.from().rank(), testfile), Empty)
		b.MakeMove(m1, false)
		if b.IsCheck(b.Turn) {
			b.UndoMove(false)
			return
		}
		b.UndoMove(false)
	}

	b.MakeMove(m, true)
	if b.IsCheck(b.oppositeturn()) {
		b.UndoMove(true)
		return
	}

	if b.IsCheck(b.Turn) {
		causesCheck = true
		if !b.CanEvadeMate(b.Turn) {
			causesMate = true
		}
	}

	legal = true
	b.UndoMove(true)
	return
}

// MakeMove Makes a move
func (b *BoardType) MakeMove(m MoveType, andSwitch bool) {

	// fmt.Printf("Making Move: %s - cap: %s\n",m.ToSAN(),m.captured().String())
	// PrintDashboard(b)

	if andSwitch {
		//   b.PrintBoard("Before MakeMove")
	}
	from := m.from()
	to := m.to()
	captured := m.captured()
	notFromBB := ^BitBoardType(0x1 << from)
	toBB := BitBoardType(0x1 << to)
	movingPiece := m.piece() //b.Squares[from]
	var promotionPiece PieceType
	var promotionPieceKind PieceBaseType
	if m.isPromotionMove() {
		promotionPiece = m.getPromotionTo()
		promotionPieceKind = promotionPiece.Kind()
	}
	movingPieceKind := movingPiece.Kind()

	switch movingPiece.Color() {
	case White:
		b.Whites[movingPieceKind] &= notFromBB
		if m.isPromotionMove() {
			b.Whites[promotionPieceKind] |= toBB
		} else {
			b.Whites[movingPieceKind] |= toBB
		}
		b.WhitePieces &= notFromBB
		b.WhitePieces |= toBB
	case Black:
		b.Blacks[movingPieceKind] &= notFromBB
		if m.isPromotionMove() {
			b.Blacks[promotionPieceKind] |= toBB
		} else {
			b.Blacks[movingPieceKind] |= toBB
		}
		b.BlackPieces &= notFromBB
		b.BlackPieces |= toBB
	}

	b.Occupied &= notFromBB
	b.Occupied |= toBB

	if m.isPromotionMove() {
		b.Squares[m.to()] = promotionPiece
	} else {
		b.Squares[m.to()] = b.Squares[m.from()]
	}
	b.Squares[m.from()] = Empty
	if captured != Empty {
		if m.isEnPassant() {
			capfrom := CoordsToSquare(from.rank(), to.file())
			b.Squares[capfrom] = Empty
			b.capturePiece(capfrom, captured)
		} else {
			b.capturePiece(to, captured)
		}
	}

	b.WhiteCastlingHistory = append(b.WhiteCastlingHistory, b.WhiteCastling)
	b.BlackCastlingHistory = append(b.BlackCastlingHistory, b.BlackCastling)
	if m.piece().Kind() == King {
		if b.Turn == White {
			b.WhiteCastling = 0
			if m.isCastlingMove() {
				b.finishCastle(m.to(), White)
			}
		} else {
			b.BlackCastling = 0
			if m.isCastlingMove() {
				b.finishCastle(m.to(), Black)
			}
		}
	}

	if m.piece().Kind() == Rook {
		switch m.from() {
		case WhiteKingSideRookCastleSquare:
			b.WhiteCastling &= ^CastleWhiteKingSide
		case WhiteQueenSideRookCastleSquare:
			b.WhiteCastling &= ^CastleWhiteQueenSide
		case BlackKingSideRookCastleSquare:
			b.BlackCastling &= ^CastleBlackKingSide
		case BlackQueenSideRookCastleSquare:
			b.BlackCastling &= ^CastleBlackQueenSide
		}
	}

	b.EnPassantHistory = append(b.EnPassantHistory, b.EnPassant)
	b.EnPassant = 0
	if m.piece() == WhitePawn {
		if from.rank() == 6 && to.rank() == 4 {
			b.EnPassant = CoordsToSquare(5, to.file())
		}
	}
	if m.piece() == BlackPawn {
		if from.rank() == 1 && to.rank() == 3 {
			b.EnPassant = CoordsToSquare(2, to.file())
		}
	}

	b.Check = m.causesCheck()
	b.HalfMoves++
	if m.piece().Color() == White {
		b.FullMoves++
	}

	b.MoveHistory = append(b.MoveHistory, m)
	if andSwitch {
		b.switchTurn()
		//    b.PrintBoard("After MakeMove")
	}
}

func (b *BoardType) finishCastle(toSquare SquareType, color ColorType) {
	var rookfrom int
	var rookto int
	var rook PieceType
	switch {
	case toSquare == G1:
		rookfrom = H1
		rookto = F1
		rook = WhiteRook
	case toSquare == C1:
		rookfrom = A1
		rookto = D1
		rook = WhiteRook
	case toSquare == G8:
		rookfrom = H8
		rookto = F8
		rook = BlackRook
	case toSquare == C8:
		rookfrom = A8
		rookto = D8
		rook = BlackRook
	}
	rookNot := ^BitBoardType(0x1 << uint(rookfrom))
	rookAdd := BitBoardType(0x1 << uint(rookto))
	b.Squares[rookfrom] = Empty
	b.Squares[rookto] = rook
	b.Occupied &= rookNot
	b.Occupied |= rookAdd

	if color == White {
		b.Whites[Rook] &= rookNot
		b.Whites[Rook] |= rookAdd
		b.WhitePieces &= rookNot
		b.WhitePieces |= rookAdd
	} else {
		b.Blacks[Rook] &= rookNot
		b.Blacks[Rook] |= rookAdd
		b.BlackPieces &= rookNot
		b.BlackPieces |= rookAdd
	}
}

func (b *BoardType) unFinishCastle(fromSquare SquareType, color ColorType) {
	var rookfrom int
	var rookto int
	var rook PieceType
	switch {
	case fromSquare == G1:
		rookfrom = F1
		rookto = H1
		rook = WhiteRook
	case fromSquare == C1:
		rookfrom = D1
		rookto = A1
		rook = WhiteRook
	case fromSquare == G8:
		rookfrom = F8
		rookto = H8
		rook = BlackRook
	case fromSquare == C8:
		rookfrom = D8
		rookto = A8
		rook = BlackRook
	}
	rookNot := ^BitBoardType(0x1 << uint(rookfrom))
	rookAdd := BitBoardType(0x1 << uint(rookto))
	b.Squares[rookfrom] = Empty
	b.Squares[rookto] = rook
	b.Occupied &= rookNot
	b.Occupied |= rookAdd

	if color == White {
		b.Whites[Rook] &= rookNot
		b.Whites[Rook] |= rookAdd
		b.WhitePieces &= rookNot
		b.WhitePieces |= rookAdd
	} else {
		b.Blacks[Rook] &= rookNot
		b.Blacks[Rook] |= rookAdd
		b.BlackPieces &= rookNot
		b.BlackPieces |= rookAdd
	}
}

//UndoMove undo a move
func (b *BoardType) UndoMove(andSwitch bool) {
	if len(b.MoveHistory) == 0 {
		fmt.Printf("Say What??")
		return
	}
	var m MoveType
	m, b.MoveHistory = b.MoveHistory[len(b.MoveHistory)-1], b.MoveHistory[:len(b.MoveHistory)-1]
	b.WhiteCastling, b.WhiteCastlingHistory = b.WhiteCastlingHistory[len(b.WhiteCastlingHistory)-1], b.WhiteCastlingHistory[:len(b.WhiteCastlingHistory)-1]
	b.BlackCastling, b.BlackCastlingHistory = b.BlackCastlingHistory[len(b.BlackCastlingHistory)-1], b.BlackCastlingHistory[:len(b.BlackCastlingHistory)-1]
	b.EnPassant, b.EnPassantHistory = b.EnPassantHistory[len(b.EnPassantHistory)-1], b.EnPassantHistory[:len(b.EnPassantHistory)-1]

	// fmt.Printf(" Undoing %s\n",m.ToSAN())

	from := m.from()
	to := m.to()
	captured := m.captured()
	fromBB := BitBoardType(0x1 << from)
	notToBB := ^BitBoardType(0x1 << to)
	movingPiece := m.piece() //  b.Squares[to]
	movingPieceKind := movingPiece.Kind()
	var promotionPiece PieceType
	var promotionPieceKind PieceBaseType
	if m.isPromotionMove() {
		promotionPiece = m.getPromotionTo()
		promotionPieceKind = promotionPiece.Kind()
	}
	switch movingPiece.Color() {
	case White:
		b.Whites[movingPieceKind] |= fromBB
		if m.isPromotionMove() {
			b.Whites[promotionPieceKind] &= notToBB
		} else {
			b.Whites[movingPieceKind] &= notToBB
		}
		b.WhitePieces |= fromBB
		b.WhitePieces &= notToBB
	case Black:
		b.Blacks[movingPieceKind] |= fromBB
		if m.isPromotionMove() {
			b.Blacks[promotionPieceKind] &= notToBB
		} else {
			b.Blacks[movingPieceKind] &= notToBB
		}
		b.BlackPieces |= fromBB
		b.BlackPieces &= notToBB
	}

	if m.isPromotionMove() {
		b.Squares[from] = PieceType(movingPiece.Color() | Pawn)
	} else {
		b.Squares[from] = b.Squares[to]
	}
	b.Squares[to] = captured
	if m.isEnPassant() {
		b.Squares[to] = Empty
		b.Squares[CoordsToSquare(from.rank(), to.file())] = captured
	}

	b.Occupied |= fromBB
	if captured == Empty {
		b.Occupied &= notToBB
	} else {
		if m.isEnPassant() {
			b.uncapturePiece(CoordsToSquare(from.rank(), to.file()), captured)
		} else {
			b.uncapturePiece(to, captured)
		}
	}

	if movingPieceKind == King {
		if m.isCastlingMove() {
			b.unFinishCastle(m.to(), movingPiece.Color())
		}
	}

	if andSwitch {
		b.switchTurn()
	}
}

//Remove captured piece from opponent's pieces
func (b *BoardType) capturePiece(sq SquareType, captured PieceType) {
	if captured == Empty {
		return
	}
	sqBB := BitBoardType(0x1 << sq)
	kind := captured.Kind()
	switch captured.Color() {
	case White:
		b.Whites[kind] &= ^sqBB
		b.WhitePieces &= ^sqBB
	case Black:
		b.Blacks[kind] &= ^sqBB
		b.BlackPieces &= ^sqBB
	}
}

//Undo capture piece
func (b *BoardType) uncapturePiece(sq SquareType, captured PieceType) {
	if captured == Empty {
		return
	}
	sqBB := BitBoardType(0x1 << sq)
	kind := captured.Kind()
	switch captured.Color() {
	case White:
		b.Whites[kind] |= sqBB
		b.WhitePieces |= sqBB
	case Black:
		b.Blacks[kind] |= sqBB
		b.BlackPieces |= sqBB
	}
}

//Generates King & Knight pseudo-legal moves
func (b *BoardType) genFromMoves(piece PieceType, pieces, ours BitBoardType, attackFrom []BitBoardType) []MoveType {
	moves := []MoveType(nil)
	for pieces > 0 {
		from := pieces.popLSB()
		targets := attackFrom[from] & ^ours
		for targets > 0 {
			to := targets.popLSB()
			m := NewMove(piece, SquareType(from), SquareType(to), b.Squares[to])
			moves = append(moves, m)
		}
	}
	return moves
}

//Generate sliding-piece's pseudo-legal moves
func (b *BoardType) genRayMoves(piece PieceType, pieces, ours BitBoardType, theirs BitBoardType, directions []Direction) []MoveType {
	moves := []MoveType(nil)
	for pieces > 0 {
		from := pieces.popLSB()
		var allTargets, targets BitBoardType
		for _, direction := range directions {
			targets = RayMasks[direction][from]
			blockers := targets & b.Occupied
			// PrintCalcboard(b,targets,"Targets",blockers,"Blockers",ours,"Ours")
			if blockers > 0 {
				if DirectionLSBMSP[direction] == LSB {
					targets ^= RayMasks[direction][blockers.lsbIndex()]
				} else {
					targets ^= RayMasks[direction][blockers.msbIndex()]
				}
			}
			allTargets |= targets & (^ours)
		}
		for allTargets > 0 {
			to := allTargets.popLSB()
			moves = append(moves, NewMove(piece, SquareType(from), SquareType(to), b.Squares[to]))
		}
	}
	return moves
}

//Generate castling pseudo-legal moves
func (b *BoardType) genCastling(piece PieceType) []MoveType {
	moves := []MoveType(nil)
	switch b.Turn {
	case White:
		if (b.WhiteCastling&CastleWhiteKingSide > 0) && (b.Occupied&(0x6<<60) == 0) {
			from := SquareType(b.Whites[King].lsbIndex())
			to := SquareType(WhiteKingSideKingCastleSquare)
			moves = append(moves, NewCastlingMove(piece, from, to))
		}
		if (b.WhiteCastling&CastleWhiteQueenSide > 0) && (b.Occupied&(0x70<<53) == 0) {
			from := SquareType(b.Whites[King].lsbIndex())
			to := SquareType(WhiteQueenSideKingCastleSquare)
			moves = append(moves, NewCastlingMove(piece, from, to))
		}
	case Black:
		if (b.BlackCastling&CastleBlackKingSide > 0) && (b.Occupied&0x60 == 0) {
			from := SquareType(b.Blacks[King].lsbIndex())
			to := SquareType(BlackKingSideKingCastleSquare)
			moves = append(moves, NewCastlingMove(piece, from, to))
		}
		if (b.BlackCastling&CastleBlackQueenSide > 0) && (b.Occupied&0x6 == 0) {
			from := SquareType(b.Blacks[King].lsbIndex())
			to := SquareType(BlackQueenSideKingCastleSquare)
			moves = append(moves, NewCastlingMove(piece, from, to))
		}
	}
	return moves
}

//Generate Pawn-one-step-forward pseudo-legal moves
func (b *BoardType) genPawnOneStep(piece PieceType) []MoveType {
	moves := []MoveType(nil)
	var targets BitBoardType
	var shift = 8
	if b.Turn == White {
		targets = (b.Whites[Pawn] >> 8) & ^b.Occupied
	} else {
		targets = (b.Blacks[Pawn] << 8) & ^b.Occupied
		shift = -8
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		if (b.Turn == White && (BitBoardType(0x1<<to)&Rank8BB > 0)) || (b.Turn == Black && (BitBoardType(0x1<<to)&Rank1BB > 0)) {
			// fmt.Printf("Promoting via genPawnOneStep for color: %s from: %s to: %s \n",b.Turn,SquareType(from).ToRune(), SquareType(to).ToRune())
			moves = append(moves, NewPromotionMove(piece, SquareType(from), SquareType(to), b.Squares[to], PieceType(b.Turn|Queen)))
			moves = append(moves, NewPromotionMove(piece, SquareType(from), SquareType(to), b.Squares[to], PieceType(b.Turn|Rook)))
			moves = append(moves, NewPromotionMove(piece, SquareType(from), SquareType(to), b.Squares[to], PieceType(b.Turn|Knight)))
			moves = append(moves, NewPromotionMove(piece, SquareType(from), SquareType(to), b.Squares[to], PieceType(b.Turn|Bishop)))
		} else {
			moves = append(moves, NewMove(piece, SquareType(from), SquareType(to), b.Squares[to]))
		}
	}
	return moves
}

//Generate Pawn-two-step-forward pseudo-legal moves
func (b *BoardType) genPawnTwoSteps(piece PieceType) []MoveType {
	moves := []MoveType(nil)
	var targets BitBoardType
	var shift int
	if b.Turn == White {
		targets = ((b.Whites[Pawn] & Rank2BB) >> 16) & (^(Rank3BB & b.Occupied) >> 8) & ^b.Occupied
		shift = 16
	} else {
		targets = ((b.Blacks[Pawn] & Rank7BB) << 16) & (^(Rank6BB & b.Occupied) << 8) & ^b.Occupied
		shift = -16
	}
	for targets > 0 {
		to := targets.popLSB()
		from := int(to) + shift
		moves = append(moves, NewMove(piece, SquareType(from), SquareType(to), b.Squares[to]))
	}
	return moves
}

//Generate pawns left and right attacks
func (b *BoardType) genPawnAttacks(side ColorType, piece PieceType, theirsAll BitBoardType) []MoveType {
	moves := []MoveType(nil)
	var ours BitBoardType
	if side == White {
		ours = b.Whites[Pawn]
	} else {
		ours = b.Blacks[Pawn]
	}
	var targets BitBoardType
	enPassant := BitBoardType(0x1<<uint(b.EnPassant)) & (Rank3BB | Rank6BB)
	for _, shift := range [2]int{7, 9} {
		fromShift := shift
		if b.Turn == White {
			if shift == 9 {
				targets = BitBoardType((ours&^FileABB)>>uint(shift)) & (theirsAll | enPassant)
			} else {
				targets = BitBoardType((ours&^FileHBB)>>uint(shift)) & (theirsAll | enPassant)
			}
		} else {
			if shift == 9 {
				targets = BitBoardType((ours&^FileHBB)<<uint(shift)) & (theirsAll | enPassant)
			} else {
				targets = BitBoardType((ours&^FileABB)<<uint(shift)) & (theirsAll | enPassant)
			}
			fromShift *= -1
		}
		for targets > 0 {
			to := SquareType(targets.popLSB())
			from := SquareType(int(to) + fromShift)
			if b.EnPassant > 0 && to == b.EnPassant {
				if b.Turn == White {
					moves = append(moves, NewEnPassantMove(piece, from, to, BlackPawn))
				} else {
					moves = append(moves, NewEnPassantMove(piece, from, to, WhitePawn))
				}
			} else {
				if (b.Turn == White && (BitBoardType(0x1<<to)&Rank8BB > 0)) || (b.Turn == Black && (BitBoardType(0x1<<to)&Rank1BB > 0)) {
					//           fmt.Printf("Promoting via genPawnAttacks for color: %s from: %s to: %s \n",b.Turn,from.ToRune(), to.ToRune())
					moves = append(moves, NewPromotionMove(piece, from, to, b.Squares[to], PieceType(b.Turn|Queen)))
					moves = append(moves, NewPromotionMove(piece, from, to, b.Squares[to], PieceType(b.Turn|Rook)))
					moves = append(moves, NewPromotionMove(piece, from, to, b.Squares[to], PieceType(b.Turn|Knight)))
					moves = append(moves, NewPromotionMove(piece, from, to, b.Squares[to], PieceType(b.Turn|Bishop)))
				} else {
					moves = append(moves, NewMove(piece, from, to, b.Squares[to]))
				}
			}

		}
	}
	return moves
}

// CanEvadeMate can get out of checkmate
func (b *BoardType) CanEvadeMate(side ColorType) bool {
	var ours [7]BitBoardType
	var oursAll BitBoardType
	var theirsAll BitBoardType
	if side == White {
		ours = b.Whites
		oursAll = b.WhitePieces
		theirsAll = b.BlackPieces
	} else {
		ours = b.Blacks
		oursAll = b.BlackPieces
		theirsAll = b.WhitePieces
	}

	defendMoves := []MoveType(nil)

	defendMoves = append(defendMoves, b.genFromMoves(PieceType(side|King), ours[King], oursAll, KingAttacksFrom[:])...)
	defendMoves = append(defendMoves, b.genFromMoves(PieceType(side|Knight), ours[Knight], oursAll, KnightAttacksFrom[:])...)
	defendMoves = append(defendMoves, b.genRayMoves(PieceType(side|Bishop), ours[Bishop], oursAll, theirsAll, BishopDirections[:])...)
	defendMoves = append(defendMoves, b.genRayMoves(PieceType(side|Rook), ours[Rook], oursAll, theirsAll, RookDirections[:])...)
	defendMoves = append(defendMoves, b.genRayMoves(PieceType(side|Queen), ours[Queen], oursAll, theirsAll, BishopDirections[:])...)
	defendMoves = append(defendMoves, b.genRayMoves(PieceType(side|Queen), ours[Queen], oursAll, theirsAll, RookDirections[:])...)

	defendMoves = append(defendMoves, b.genPawnTwoSteps(PieceType(side|Pawn))...)
	defendMoves = append(defendMoves, b.genPawnAttacks(side, PieceType(side|Pawn), theirsAll)...)
	defendMoves = append(defendMoves, b.genPawnOneStep(PieceType(side|Pawn))...)

	escapeCheck := false
	// fmt.Printf("Defensive MoveList: %v\n", defendMoves)
	for _, move := range defendMoves {
		b.MakeMove(move, true)
		if !b.IsCheck(b.oppositeturn()) {
			escapeCheck = true
			//     fmt.Printf("Escaped Mate with: %s\n", move.ToSAN())
		}
		b.UndoMove(true)
	}

	return escapeCheck
}

// IsCheck Checks whether our king is in check or not
func (b *BoardType) IsCheck(side ColorType) bool {
	var kingBB, theirsAll, attackers, targets BitBoardType
	var theirs []BitBoardType
	if side == White {
		kingBB, theirs, theirsAll = b.Whites[King], b.Blacks[:], b.BlackPieces
	} else {
		kingBB, theirs, theirsAll = b.Blacks[King], b.Whites[:], b.WhitePieces
	}
	kingIdx := kingBB.lsbIndex()
	possibleAttackers := theirsAll & AttacksTo[kingIdx]

	attackers = (theirs[Rook] | theirs[Queen]) & possibleAttackers
	if attackers > 0 && b.isCheckedFromRay(kingBB, attackers, RookDirections[:]) {
		return true
	}

	attackers = (theirs[Bishop] | theirs[Queen]) & possibleAttackers
	if attackers > 0 && b.isCheckedFromRay(kingBB, attackers, BishopDirections[:]) {
		return true
	}

	attackers = theirs[Knight] & possibleAttackers
	for attackers > 0 {
		from := attackers.popLSB()
		if KnightAttacksFrom[from]&kingBB > 0 {
			return true
		}
	}

	enPassant := BitBoardType(0x1<<uint(b.EnPassant)) & (Rank6BB | Rank3BB)
	if side == White {
		targets = BitBoardType((b.Blacks[Pawn]&^FileABB)<<uint(7)) & (b.WhitePieces | enPassant)
		targets |= BitBoardType((b.Blacks[Pawn]&^FileHBB)<<uint(9)) & (b.WhitePieces | enPassant)
	} else {
		targets = BitBoardType((b.Whites[Pawn]&^FileHBB)>>uint(7)) & (b.BlackPieces | enPassant)
		targets |= BitBoardType((b.Whites[Pawn]&^FileABB)>>uint(9)) & (b.BlackPieces | enPassant)
	}
	if targets&kingBB > 0 {
		return true
	}

	attackers = theirs[King] & possibleAttackers
	if attackers > 0 {
		from := attackers.popLSB()
		if KingAttacksFrom[from]&kingBB > 0 {
			return true
		}
	}
	return false
}

//checks whether target is attacked by one of the "attackers"
func (b *BoardType) isCheckedFromRay(target, attackers BitBoardType, directions []Direction) bool {
	var targets BitBoardType
	var from uint
	for attackers > 0 {
		from = attackers.popLSB()
		for _, direction := range directions {
			targets = RayMasks[direction][from]
			blockers := targets & b.Occupied
			if blockers > 0 {
				if DirectionLSBMSP[direction] == LSB {
					targets ^= RayMasks[direction][blockers.lsbIndex()]
				} else {
					targets ^= RayMasks[direction][blockers.msbIndex()]
				}
			}
			if targets&target > 0 {
				return true
			}
		}
	}
	return false
}
