package models

import (
	"fmt"
	"testing"
)

func printBoard(b *Board) {
	for _, ship := range b.Battleships {
		fmt.Println(ship.Cells, ship.IsVertical)
	}
	fmt.Println("--------------------------")
}

func TestFindRandomSpace(t *testing.T) {
	for i := 0; i < 1000; i++ {
		board, err := GenerateBoard(false)
		if err != nil {
			t.Fatal("failed to create a board:", err)
		}

		cellMap := map[int][]int{} //map[X][]Y

		for _, ship := range board.Battleships {
			for _, cell := range ship.Cells {
				for _, y := range cellMap[cell.X] {
					if y == cell.Y {
						t.Errorf(fmt.Sprintf("duplicate cells at (%d, %d)", cell.X, cell.Y))
						printBoard(board)
					}
				}
				cellMap[cell.X] = append(cellMap[cell.X], cell.Y)
			}
		}
	}
}
