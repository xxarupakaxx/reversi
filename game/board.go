package game

type Board struct {
	Cells [][]Color
}

func NewBoard() *Board {
	b := &Board{Cells: make([][]Color, 10)}

	for i := 0; i < 10; i++ {
		b.Cells[i] = make([]Color, 10)
	}

	for i := 0; i < 10; i++ {
		b.Cells[0][i] = Wall
	}
	for i := 1; i < 9; i++ {
		b.Cells[i][0] = Wall
		b.Cells[i][9] = Wall
	}

	for i := 0; i < 9; i++ {
		b.Cells[9][i] = Wall
	}

	b.Cells[4][4] = White
	b.Cells[5][5] = White
	b.Cells[5][4] = Blank
	b.Cells[4][5] = Blank

	return b
}
