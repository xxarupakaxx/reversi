package game

type Color int

const (
	Empty Color = iota
	Blank
	White
	Wall
	None
)

func ColorToStr(c Color) string {
	switch c {
	case Blank:
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
	case Blank:
		return White
	case White:
		return Blank
	}

	panic("私の色は何？？？？？？？？？")
}
