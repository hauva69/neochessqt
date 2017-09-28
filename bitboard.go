package main

import (
    "strings"  
    "fmt"
)

// BitBoardType enum
type BitBoardType uint64

var (
	// MS1BTABLE Table of Most Sig Bits
	MS1BTABLE = [256]uint{
		0, 0, 1, 2, 2, 2, 2, 2, 3, 3,
		3, 3, 3, 3, 3, 3, 4, 4, 4, 4,
		4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		4, 4, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		6, 6, 6, 6, 6, 6, 6, 6, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7,
	}

	// LS1BTABLE Table of Least Sig Bits
	LS1BTABLE = [64]uint{
		0, 1, 48, 2, 57, 49, 28, 3,
		61, 58, 50, 42, 38, 29, 17, 4,
		62, 55, 59, 36, 53, 51, 43, 22,
		45, 39, 33, 30, 24, 18, 12, 5,
		63, 47, 56, 27, 60, 41, 37, 16,
		54, 35, 52, 21, 44, 32, 23, 11,
		46, 26, 40, 15, 34, 20, 31, 10,
		25, 14, 19, 9, 13, 8, 7, 6,
	}
)


//Get least significant bit
func (bb BitBoardType) lsb() BitBoardType {
	return bb & (-bb)
}

//Get index of least significant(of Martin Läuter)
func (bb BitBoardType) lsbIndex() uint {
	return LS1BTABLE[(bb.lsb()*0x03f79d71b4cb0a89)>>58]
}

//Get index of most significant bit(of Eugene Nalimov)
func (bb BitBoardType) msbIndex() int {
	var msb int
	if bb > 0xFFFFFFFF {
		bb >>= 32
		msb = 32
	}
	if bb > 0xFFFF {
		bb >>= 16
		msb += 16
	}
	if bb > 0xFF {
		bb >>= 8
		msb += 8
	}
	return msb + int(MS1BTABLE[bb])
}

//Pop least significant bit and return it's index
func (bb *BitBoardType) popLSB() uint {
	lsb := (*bb).lsb()
	*bb -= lsb
	return lsb.lsbIndex()
}

// PrintBitboard d
func PrintBitboard(bb BitBoardType, title string) {
    bitboard := sPrintBitboard(bb, title)
    for _, line := range bitboard {
        fmt.Printf("%s\n", line)
    }
}

func sPrintBitboard(bb BitBoardType, title string) []string {
    s := []string(nil)
    title = strings.Trim(title," ")
    if (len(title) > 22) {
        title = title[:22]
    }
    ls:=(24-len(title))/2
    rs:=24-len(title)-ls
    
    title = strings.Repeat("_",ls)+title+strings.Repeat("_",rs)    
	var shiftMask uint64 = 1
	//bb.printB("pp2")
	s = append(s,fmt.Sprintf("%s", title))
	for rank := 7; rank >= 0; rank-- {
        var line string 
		for file := 7; file >= 0; file-- {
			var squareIdx = uint(rank*8 + file)
			if bb&BitBoardType(shiftMask<<squareIdx) > 0 {
				line += fmt.Sprintf(" X ")
			} else {
				line += fmt.Sprintf(" _ ")
			}
		}
        s = append(s,line)		
	}
    return s
}