package main

import (
	// "bytes"
	"encoding/binary"
)

// GameHeader Structure
type GameHeader struct {
	Event        string
	Site         string
	Date         string
	Round        string
	White        string
	Black        string
	Result       string
	WhiteTitle   string
	BlackTitle   string
	WhiteElo     string
	BlackElo     string
	WhiteUSCF    string
	BlackUSCF    string
	WhiteNA      string
	BlackNA      string
	WhiteType    string
	BlackType    string
	WhiteFideID  string
	BlackFideID  string
	EventType    string
	EventDate    string
	EventSponsor string
	Section      string
	Stage        string
	Board        string
	Opening      string
	Variation    string
	SubVariation string
	ECO          string
	NIC          string
	Time         string
	UTCTime      string
	TimeControl  string
	SetUp        string
	FEN          string
	Termination  string
	Annotator    string
	Mode         string
	PlyCount     string
}

// Game Structure
type Game struct {
	ID uint32
	GameHeader
	OriginalMoveText    string
	Moves               []MoveType
	PositionTree        []PositionNode
	CurrentPositionNode int
}

// Game Constants
const (
	MaxGameNumber uint32 = 4294967295
)

// GamePacked struct
type GamePacked struct {
	neomagic   []byte // {'\n','N','e','o'}  // In case Indexes need rebuilding
	gamenumber []byte // uint32 covers 0 - 4294967295 games
	tags       []byte // byte 0 uint8 Number of tags
	// byte 1 uint8 Tag mark type
	// byte 2 uint8 Tag length
	// byte 3 - byte[lenth] tag content
	compresstree []byte // Move Array See Move Type for details
	// uint32 4 bite chunks until comment
	//     Coment length in mask uint16 - Max comment length 65535
	//     Comment in variation uint16  - Max comment length 65535
	neoclosemagic []byte // {'o','e','N'}
}

// TagMark binary Type
type TagMark uint8

var supportedtags = map[string]uint8{
	"Event":        1,
	"Site":         2,
	"Date":         3,
	"Round":        4,
	"White":        5,
	"Black":        6,
	"Result":       7,
	"WhiteTitle":   8,
	"BlackTitle":   9,
	"WhiteElo":     10,
	"BlackElo":     11,
	"WhiteUSCF":    12,
	"BlackUSCF":    13,
	"WhiteNA":      14,
	"BlackNA":      15,
	"WhiteType":    16,
	"BlackType":    17,
	"WhiteFideID":  18,
	"BlackFideID":  19,
	"EventType":    20,
	"EventDate":    21,
	"EventSponsor": 22,
	"Section":      23,
	"Stage":        24,
	"Board":        25,
	"Opening":      26,
	"Variation":    27,
	"SubVariation": 28,
	"ECO":          29,
	"NIC":          30,
	"Time":         31,
	"UTCTime":      32,
	"TimeControl":  33,
	"SetUp":        34,
	"FEN":          35,
	"Termination":  36,
	"Annotator":    37,
	"Mode":         38,
	"PlyCount":     39,
}
var reversesupportedtags map[uint8]string

func init() {
	reversesupportedtags = reverseMap(supportedtags)
}

func reverseMap(m map[string]uint8) map[uint8]string {
	n := make(map[uint8]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

// PackGame create
func (g *Game) PackGame() *GamePacked {
	gp := new(GamePacked)
	gp.neomagic = make([]byte, 4)
	gp.neomagic = []byte{'\n', 'N', 'e', 'o'}
	gp.neoclosemagic = make([]byte, 4)
	gp.neoclosemagic = []byte{'o', 'e', 'N'}
	gp.gamenumber = make([]byte, 4)
	binary.BigEndian.PutUint32(gp.gamenumber, g.ID)
	gp.compresstree = make([]byte, 4*len(g.Moves))
	bp := 0
	for _, Move := range g.Moves {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(Move.Move))
		bp += copy(gp.compresstree[bp:], b)
	}
	totaltags := uint8(0)
	totalbytes := 1 // byte 0 is for count of tags
	for k := range supportedtags {
		tag := g.gettag(k)
		if tag != "" {
			totaltags++
			totalbytes += len(tag) + 2 // 1 byte for mark, 1 for length, rest for content
		}
	}
	bp = 0
	gp.tags = make([]byte, totalbytes)
	gp.tags[bp] = uint8(totaltags) // Total Tag Count
	bp++
	for k, tagmark := range supportedtags {
		tag := g.gettag(k)
		if tag != "" {
			gp.tags[bp] = uint8(tagmark) // Tag Marker
			bp++
			cpcount := copy(gp.tags[bp+1:], tag)
			gp.tags[bp] = uint8(cpcount) // Tag Lengh
			bp += cpcount + 1
		}
	}
	return gp
}

// LoadPackedGame from buffer
func LoadPackedGame(buffer []byte) (*GamePacked, *Game) {
	g := NewGame()
	gp := new(GamePacked)
	bp := 0
	gp.neomagic = make([]byte, 4)
	bp += copy(gp.neomagic, buffer[bp:])
	g.ID = binary.BigEndian.Uint32(buffer[bp:])
	gp.gamenumber = make([]byte, 4)
	copy(gp.gamenumber, buffer[bp:])
	bp += 4
	totaltags := int(uint8(buffer[bp]))
	bp++
	for i := 0; i < totaltags; i++ {
		tagmark := uint8(buffer[bp])
		bp++
		taglength := uint8(buffer[bp])
		bp++
		tagvaluebytes := make([]byte, taglength)
		copy(tagvaluebytes, buffer[bp:])
		temptagbytes := make([]byte, taglength+2)
		temptagbytes[0] = tagmark
		temptagbytes[1] = taglength
		copy(temptagbytes[2:], tagvaluebytes)
		if len(gp.tags) == 0 {
			gp.tags = make([]byte, taglength+3)
			gp.tags[0] = uint8(totaltags)
			copy(gp.tags[1:], temptagbytes)
		} else {
			currenttags := make([]byte, len(gp.tags))
			copy(currenttags, gp.tags)
			gp.tags = make([]byte, len(currenttags)+int(taglength)+2)
			copy(gp.tags, currenttags)
			copy(gp.tags[len(currenttags)+1:], temptagbytes)
		}
		tagvalue := string(tagvaluebytes)
		tagname := reversesupportedtags[tagmark]
		g.settag(tagname, tagvalue)
		bp += int(taglength)
	}
	remaining := len(buffer) - bp - 3
	for cp := bp; cp < remaining; cp += 4 {
		move := binary.BigEndian.Uint32(buffer[cp:])
		Move := MoveType{}
		Move.Move = PackedMoveType(uint32(move))
		g.Moves = append(g.Moves, Move)
	}
	gp.compresstree = make([]byte, remaining)
	n := copy(gp.compresstree, buffer[bp:])
	gp.neoclosemagic = make([]byte, 3)
	copy(gp.neoclosemagic, buffer[bp+n:])
	return gp, g
}

// Bytes Return Bytes of PackedGame
func (gp *GamePacked) Bytes() []byte {
	bytecount := 0
	bytecount += len(gp.neomagic)
	bytecount += len(gp.gamenumber)
	bytecount += len(gp.tags)
	bytecount += len(gp.compresstree)
	bytecount += len(gp.neoclosemagic)
	bytes := make([]byte, bytecount)
	bp := 0
	bp = copy(bytes[bp:], gp.neomagic)
	bp += copy(bytes[bp:], gp.gamenumber)
	bp += copy(bytes[bp:], gp.tags)
	bp += copy(bytes[bp:], gp.compresstree)
	bp += copy(bytes[bp:], gp.neoclosemagic)
	return bytes
}

// NewGame create
func NewGame() *Game {
	g := new(Game)
	g.ID = 0
	return g
}

// GobEncode encode game
//func GobEncode(g *GamePacked) ([]byte, error) {
//	var network bytes.Buffer
// enc := gob.NewEncoder(&network)
// err := enc.Encode(g)
//	if err != nil {
//		return nil, err
//	}
//	return network.Bytes(), nil
//}

// LoadMoves from board into ActiveGame
func (g *Game) LoadMoves(cb *BoardType) {
	node := PositionNode{}
	node.nodetype = MoveNode
	node.potential = cb.GenerateLegalMoves()
	g.PositionTree = append(g.PositionTree, node)
	g.CurrentPositionNode = len(g.PositionTree) - 1
	/*
		g.CurrentFen = cb.ToFen()
		g.SideToMove = "w"
		if cb.Turn == Black {
			g.SideToMove = "b"
		}
		g.Ischeck = cb.IsCheck(cb.Turn)
		g.Ischeckmate = cb.Checkmate // Needs to be recalced
		g.CurrentPgn = ""
		g.PotentialMoves = cb.GenerateLegalMoves()
		countpMoves := len(g.PotentialMoves)
		g.FromToMoves = make([]string, countpMoves)
		g.SanMoves = make([]string, countpMoves)
		for pindex, pMove := range g.PotentialMoves {
			g.FromToMoves[pindex] = pMove.from().ToRune() + pMove.to().ToRune()
			g.SanMoves[pindex] = pMove.ToSAN()
		}
		mn := 0
		g.CurrentPgn = ""
		linelength := 0
		for _, Move := range g.Moves {
			if Move.color() == White {
				mn++
				movestr := strconv.Itoa(mn) + ". "
				if (linelength + len(movestr)) > 80 {
					g.CurrentPgn += "\n"
					linelength = 0
				}
				g.CurrentPgn += movestr
				linelength += len(movestr)
			}
			movestr := Move.ToSAN() + " "
			if (linelength + len(movestr)) > 80 {
				g.CurrentPgn += "\n"
				linelength = 0
			}
			g.CurrentPgn += movestr
			linelength += len(movestr)
		}
	*/
}

func (g *Game) settag(tagname string, tagvalue string) {
	switch tagname {
	case "Event":
		g.Event = tagvalue
	case "Site":
		g.Site = tagvalue
	case "Date":
		g.Date = tagvalue
	case "Round":
		g.Round = tagvalue
	case "White":
		g.White = tagvalue
	case "Black":
		g.Black = tagvalue
	case "Result":
		g.Result = tagvalue
	case "WhiteTitle":
		g.WhiteTitle = tagvalue
	case "BlackTitle":
		g.BlackTitle = tagvalue
	case "WhiteElo":
		g.WhiteElo = tagvalue
	case "BlackElo":
		g.BlackElo = tagvalue
	case "WhiteUSCF":
		g.WhiteUSCF = tagvalue
	case "BlackUSCF":
		g.BlackUSCF = tagvalue
	case "WhiteNA":
		g.WhiteNA = tagvalue
	case "BlackNA":
		g.BlackNA = tagvalue
	case "WhiteType":
		g.WhiteType = tagvalue
	case "BlackType":
		g.BlackType = tagvalue
	case "WhiteFideId":
		g.WhiteFideID = tagvalue
	case "BlackFideId":
		g.BlackFideID = tagvalue
	case "EventType":
		g.EventType = tagvalue
	case "EventDate":
		g.EventDate = tagvalue
	case "EventSponsor":
		g.EventSponsor = tagvalue
	case "Section":
		g.Section = tagvalue
	case "Stage":
		g.Stage = tagvalue
	case "Board":
		g.Board = tagvalue
	case "Opening":
		g.Opening = tagvalue
	case "Variation":
		g.Variation = tagvalue
	case "SubVariation":
		g.SubVariation = tagvalue
	case "ECO":
		g.ECO = tagvalue
	case "NIC":
		g.NIC = tagvalue
	case "Time":
		g.Time = tagvalue
	case "UTCTime":
		g.UTCTime = tagvalue
	case "TimeControl":
		g.TimeControl = tagvalue
	case "SetUp":
		g.SetUp = tagvalue
	case "FEN":
		g.FEN = tagvalue
	case "Termination":
		g.Termination = tagvalue
	case "Annotator":
		g.Annotator = tagvalue
	case "Mode":
		g.Mode = tagvalue
	case "PlyCount":
		g.PlyCount = tagvalue
	}
}

func (g *Game) gettag(tagname string) string {
	switch tagname {
	case "Event":
		return g.Event
	case "Site":
		return g.Site
	case "Date":
		return g.Date
	case "Round":
		return g.Round
	case "White":
		return g.White
	case "Black":
		return g.Black
	case "Result":
		return g.Result
	case "WhiteTitle":
		return g.WhiteTitle
	case "BlackTitle":
		return g.BlackTitle
	case "WhiteElo":
		return g.WhiteElo
	case "BlackElo":
		return g.BlackElo
	case "WhiteUSCF":
		return g.WhiteUSCF
	case "BlackUSCF":
		return g.BlackUSCF
	case "WhiteNA":
		return g.WhiteNA
	case "BlackNA":
		return g.BlackNA
	case "WhiteType":
		return g.WhiteType
	case "BlackType":
		return g.BlackType
	case "WhiteFideId":
		return g.WhiteFideID
	case "BlackFideId":
		return g.BlackFideID
	case "EventType":
		return g.EventType
	case "EventDate":
		return g.EventDate
	case "EventSponsor":
		return g.EventSponsor
	case "Section":
		return g.Section
	case "Stage":
		return g.Stage
	case "Board":
		return g.Board
	case "Opening":
		return g.Opening
	case "Variation":
		return g.Variation
	case "SubVariation":
		return g.SubVariation
	case "ECO":
		return g.ECO
	case "NIC":
		return g.NIC
	case "Time":
		return g.Time
	case "UTCTime":
		return g.UTCTime
	case "TimeControl":
		return g.TimeControl
	case "SetUp":
		return g.SetUp
	case "FEN":
		return g.FEN
	case "Termination":
		return g.Termination
	case "Annotator":
		return g.Annotator
	case "Mode":
		return g.Mode
	case "PlyCount":
		return g.PlyCount
	default:
		return ""
	}
}

func (g *Game) MoveChoice(index int) MoveType {
	node := g.PositionTree[g.CurrentPositionNode]
	return node.potential[index]
}

func (g *Game) IsMoveInFromMoves(mvstr string) (bool, int) {
	node := g.PositionTree[g.CurrentPositionNode]
	for index, b := range node.potential {
		if b.ToRune() == mvstr {
			return true, index
		}
	}
	return false, -1
}

func (g *Game) GetTargetSquares(from SquareType) []SquareType {
	targetSquares := []SquareType{}
	node := g.PositionTree[g.CurrentPositionNode]
	for _, pMove := range node.potential {
		if pMove.from() == from {
			targetSquares = append(targetSquares, pMove.to())
		}
	}
	return targetSquares
}

func (g *Game) IsTarget(from SquareType, to SquareType) bool {
	istarget := false
	for _, sq := range g.GetTargetSquares(from) {
		if sq == to {
			istarget = true
			break
		}
	}
	return istarget
}
