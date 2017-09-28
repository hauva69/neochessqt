package main

// ColorType For Move Color
type ColorType uint

// Constants for Color
const (
	Black   ColorType = 16
	White   ColorType = 8
	NoColor ColorType = 0
)

func (ct ColorType) String() string {
	switch {
	case ct == 0:
		return "NoColor"
	case ct == 8:
		return "White"
	case ct == 16:
		return "Black"
	default:
		return "InvalidColor"
	}
}
