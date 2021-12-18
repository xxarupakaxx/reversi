package game

import "fmt"

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
	b.Cells[5][4] = Black
	b.Cells[4][5] = Black

	return b
}

func (b *Board) PutStone(x, y int, c Color) error {
	if !b.CanPutStone(x, y, c) {
		return fmt.Errorf("failed to put stone x=%v, y=%v color=%v", x, y, ColorToStr(c))
	}

	b.Cells[x][y] = c

	for i := -1; i < 1; i++ {
		for j := -1; j < 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if b.CountTurnableStonesByDirection(x, y, i, j, c) > 0 {
				b.TurnStonesByDirection(x, y, i, j, c)
			}
		}
	}

	return nil
}

// CanPutStone 石を置けるかどうかの判定
func (b *Board) CanPutStone(x, y int, c Color) bool {
	if b.Cells[x][y] != Empty {
		return false
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if b.CountTurnableStonesByDirection(x, y, i, j, c) > 0 {
				return true
			}

		}
	}
	return false
}

// CountTurnableStonesByDirection ある方向にひっくり返すことのできる石がいくつあるか
func (b *Board) CountTurnableStonesByDirection(x, y, dx, dy int, c Color) int {
	cnt := 0

	nx := x + dx
	ny := y + dy
	for {
		nc := b.Cells[nx][ny]

		if nc != OpponentColor(c) {
			break
		}

		cnt++

		nx += dx
		ny += dy
	}

	if cnt > 0 && b.Cells[nx][ny] == c {
		return cnt
	}

	return 0
}

// TurnStonesByDirection  石をひっくり返す
func (b *Board) TurnStonesByDirection(x, y, dx, dy int, c Color) {
	nx := x + dx
	ny := y + dy

	for true {
		nc := b.Cells[nx][ny]

		if nc != OpponentColor(c) {
			break
		}

		b.Cells[nx][ny] = c

		nx += dx
		ny += dy
	}
}

// AvailableCellCount 盤面の置けるマスの数
func (b *Board) AvailableCellCount(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.CanPutStone(i, j, c) {
				cnt++
			}
		}
	}

	return cnt
}

// Score 盤面の石の数を数える
func (b *Board) Score(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] != c {
				continue
			}
			cnt++
		}
	}

	return cnt
}

// Rest 盤面の置けるマスの数を数える
func (b *Board) Rest() int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] == Empty {
				cnt++
			}
		}
	}

	return cnt
}
