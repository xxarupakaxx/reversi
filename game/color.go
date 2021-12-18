package game

type Color int

const (
	Empty Color = iota
	Black
	White
	Wall
	None
)

func ColorToStr(c Color) string {
	switch c {
	case Black:
		return "o"
	case White:
		return "x"
	case Empty:
		return " "
	}

	return ""
}

func OpponentColor(me Color) Color {
	switch me {
	case Black:
		return White
	case White:
		return Black
	}

	panic("私の色は何？？？？？？？？？")
}
