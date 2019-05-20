package models

import (
	"errors"

	"github.com/billyboar/battleships/helpers"
)

// BoardRow - Board has predifined 10 rows
const BoardRow = 10

// Board will contain battleships
type Board struct {
	IsComputer  bool          `json:"is_computer"` // True: if board is for computer, False if not
	Battleships []*BattleShip `json:"battleships"`
	MissedShots []Cell        `json:"missed_shots"`
}

// NewBoard creates simple board
func NewBoard(isComputer bool) *Board {
	board := Board{
		IsComputer: isComputer,
	}

	return &board
}

// GenerateBoard creates new board with 2 destroyer ships and single battleship
func GenerateBoard(isComputer bool) (*Board, error) {
	board := NewBoard(isComputer)
	destroyerOne, err := NewBattleShip(DestroyerLength)
	if err != nil {
		return nil, err
	}
	destroyerTwo, err := NewBattleShip(DestroyerLength)
	if err != nil {
		return nil, err
	}
	battleShip, err := NewBattleShip(BattleShipLength)
	if err != nil {
		return nil, err
	}

	if err := board.AddNewShip(destroyerOne); err != nil {
		return nil, err
	}
	if err := board.AddNewShip(destroyerTwo); err != nil {
		return nil, err
	}
	if err := board.AddNewShip(battleShip); err != nil {
		return nil, err
	}

	return board, nil
}

// AddNewShip adds new ships to board
func (b *Board) AddNewShip(ship *BattleShip) error {
	if len(b.Battleships) == BattleShipNumber {
		return errors.New("board has full ships")
	}

	headCell := b.FindRandomHeadCell(ship)
	ship.BuildBody(headCell)

	b.Battleships = append(b.Battleships, ship)

	return nil
}

// FindShip returns battleship with given ID
func (b *Board) FindShip(shipID string) *BattleShip {
	for _, battleship := range b.Battleships {
		if battleship.ID == shipID {
			return battleship
		}
	}
	return nil
}

// FindRandomHeadCell finds head cell for a ship randomly satisfying
// condition that ships don't overlap
func (b *Board) FindRandomHeadCell(ship *BattleShip) Cell {
	// 70 total variants for vertical/horizontal battleships
	// 80 total variants for vertical/horizontal destroyers

	var xLimit, yLimit int

	if ship.IsVertical {
		xLimit = BoardRow - 1
		yLimit = BoardRow - ship.Length
	} else {
		xLimit = BoardRow - ship.Length
		yLimit = BoardRow - 1
	}

	var possibleHeadCells []Cell
	for x := 0; x <= xLimit; x++ {
		for y := 0; y <= yLimit; y++ {
			isValidCell := true
			for _, shipOnBoard := range b.Battleships {
				headCell := shipOnBoard.Cells[0]

				if ship.IsVertical && shipOnBoard.IsVertical {
					if x == headCell.X &&
						((headCell.Y > y && headCell.Y-y < ship.Length) ||
							(headCell.Y < y && y-headCell.Y < shipOnBoard.Length) ||
							headCell.Y == y) {
						isValidCell = false
						break
					}
				}

				if ship.IsVertical && !shipOnBoard.IsVertical {
					if (x >= headCell.X && x <= headCell.X+shipOnBoard.Length-1) &&
						((headCell.Y > y && headCell.Y-y < ship.Length) ||
							(headCell.Y < y && y-headCell.Y < shipOnBoard.Length) ||
							headCell.Y == y) {
						isValidCell = false
						break
					}
				}

				if !ship.IsVertical && !shipOnBoard.IsVertical {
					if y == headCell.Y &&
						((headCell.X > x && headCell.X-x < ship.Length) ||
							(headCell.X < x && x-headCell.X < shipOnBoard.Length) ||
							headCell.X == x) {
						isValidCell = false
						break
					}
				}

				if !ship.IsVertical && shipOnBoard.IsVertical {
					if (y >= headCell.Y && y <= headCell.Y+shipOnBoard.Length-1) &&
						((headCell.X > x && headCell.X-x < ship.Length) ||
							(headCell.X < x && x-headCell.X < shipOnBoard.Length) ||
							headCell.X == x) {
						isValidCell = false
						break
					}
				}
			}
			if isValidCell {
				possibleHeadCells = append(possibleHeadCells, Cell{
					X: x,
					Y: y,
				})
			}
		}
	}

	randomCellNumber := helpers.GenerateRandomInt(len(possibleHeadCells))
	return possibleHeadCells[randomCellNumber]
}

// RegisterShot marks cell as dead if shot cell equals
// to any cell of any battleships
func (b *Board) RegisterShot(shot Cell) (shotStatus bool, shipID string) {
	for _, battleship := range b.Battleships {
		for cellNum, cell := range battleship.Cells {
			if cell.Compare(&shot) {
				// when shot is correct mark single cell as dead
				battleship.Cells[cellNum].IsDead = true
				shotStatus = true
				shipID = battleship.ID
				return
			}
		}
	}

	if !shotStatus {
		b.MissedShots = append(b.MissedShots, shot)
	}

	return
}

func (b *Board) MarkShipIfDead(shipID string) *BattleShip {
	battleShip := b.FindShip(shipID)
	if battleShip == nil {
		return nil
	}

	damagedCells := battleShip.GetDamagedCells()

	if len(damagedCells) == battleShip.Length {
		battleShip.IsDead = true
		return battleShip
	}

	return nil
}

func (b *Board) GetAllShipWounds() []Cell {
	woundedCells := []Cell{}
	for _, battleship := range b.Battleships {
		woundedCells = append(woundedCells, battleship.GetDamagedCells()...)
	}

	return woundedCells
}

func (b *Board) GetDeadShips() []BattleShip {
	deadShips := []BattleShip{}
	for _, battleship := range b.Battleships {
		if battleship.IsDead {
			deadShips = append(deadShips, *battleship)
		}
	}

	return deadShips
}

type CellMap map[int]map[int]bool

func (b *Board) CalculateShot() *Cell {
	possibleCells := []Cell{}
	missedShotsMap := CellMap{}

	// populate missedShotsMap
	for _, missedShot := range b.MissedShots {
		if missedShotsMap[missedShot.X] == nil {
			missedShotsMap[missedShot.X] = make(map[int]bool)
		}
		missedShotsMap[missedShot.X][missedShot.Y] = true
	}

	woundedCellsMap := CellMap{}
	woundedShips := []*BattleShip{}
	for _, battleShip := range b.Battleships {
		if !battleShip.IsDead {
			woundedCells := battleShip.GetDamagedCells()
			for _, woundedCell := range woundedCells {
				if woundedCellsMap[woundedCell.X] == nil {
					woundedCellsMap[woundedCell.X] = make(map[int]bool)
				}
				woundedCellsMap[woundedCell.X][woundedCell.Y] = true
			}

			// wounded yet not dead ships
			if len(woundedCells) > 0 {
				woundedShips = append(woundedShips, battleShip)
			}

		}
	}

	if len(woundedShips) == 0 {
		//when computer hasn't wounded any ships
		for x := 0; x < BoardRow; x++ {
			for y := 0; y < BoardRow; y++ {
				if !missedShotsMap[x][y] && !woundedCellsMap[x][y] {
					possibleCells = append(possibleCells, Cell{
						X: x,
						Y: y,
					})
				}
			}
		}

		return &possibleCells[helpers.GenerateRandomInt(len(possibleCells))]
	}

	// calculate shots on wounded yet not dead ships
	for _, woundedShip := range woundedShips {
		woundedCells := woundedShip.GetDamagedCells()
		if woundedShip.GetDamageCount() == 1 {
			// when ship is hit only once
			possibleCells := checkAllSides(missedShotsMap, woundedCellsMap, woundedCells[0])
			return &possibleCells[helpers.GenerateRandomInt(len(possibleCells))]
		}

		// when ship is hit multiple times
		// damaged cells will adjacent all the time
		// since we are only shooting adjacent cells when ship is
		// hit only once.
		lastWoundedCell := woundedCells[len(woundedCells)-1]
		if woundedShip.IsVertical {
			possibleCells = checkVerticalCells(missedShotsMap, woundedCellsMap, lastWoundedCell)
			if len(possibleCells) == 0 {
				possibleCells = checkVerticalCells(missedShotsMap, woundedCellsMap, woundedCells[0])
			}
		} else {
			possibleCells = checkHorizontalCells(missedShotsMap, woundedCellsMap, lastWoundedCell)
			if len(possibleCells) == 0 {
				possibleCells = checkHorizontalCells(missedShotsMap, woundedCellsMap, woundedCells[0])
			}
		}

		return &possibleCells[helpers.GenerateRandomInt(len(possibleCells))]
	}

	return nil
}

func checkAllSides(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) []Cell {
	possibleVerticalCells := checkVerticalCells(missedShotsMap, woundedCellsMap, currentCell)
	possibleHorizontalCells := checkHorizontalCells(missedShotsMap, woundedCellsMap, currentCell)

	return append(possibleVerticalCells, possibleHorizontalCells...)
}

func checkVerticalCells(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) []Cell {
	possibleCells := []Cell{}
	// checking above cell
	if cell := checkAboveCell(missedShotsMap, woundedCellsMap, currentCell); cell != nil {
		possibleCells = append(possibleCells, *cell)
	}

	// checking below cell
	if cell := checkBelowCell(missedShotsMap, woundedCellsMap, currentCell); cell != nil {
		possibleCells = append(possibleCells, *cell)
	}
	return possibleCells
}

func checkHorizontalCells(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) []Cell {
	possibleCells := []Cell{}
	if cell := checkLeftCell(missedShotsMap, woundedCellsMap, currentCell); cell != nil {
		possibleCells = append(possibleCells, *cell)
	}

	// checking right cell
	if cell := checkRightCell(missedShotsMap, woundedCellsMap, currentCell); cell != nil {
		possibleCells = append(possibleCells, *cell)
	}
	return possibleCells
}

func checkAboveCell(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) *Cell {
	tempCoordinate := currentCell.Y - 1
	if tempCoordinate >= 0 && !woundedCellsMap[currentCell.X][tempCoordinate] && !missedShotsMap[currentCell.X][tempCoordinate] {
		return &Cell{
			X: currentCell.X,
			Y: tempCoordinate,
		}
	}
	return nil
}

func checkBelowCell(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) *Cell {
	tempCoordinate := currentCell.Y + 1
	if tempCoordinate < BoardRow && !woundedCellsMap[currentCell.X][tempCoordinate] && !missedShotsMap[currentCell.X][tempCoordinate] {
		return &Cell{
			X: currentCell.X,
			Y: tempCoordinate,
		}
	}
	return nil
}

func checkLeftCell(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) *Cell {
	tempCoordinate := currentCell.X - 1
	if tempCoordinate >= 0 && !woundedCellsMap[tempCoordinate][currentCell.Y] && !missedShotsMap[tempCoordinate][currentCell.Y] {
		return &Cell{
			X: tempCoordinate,
			Y: currentCell.Y,
		}
	}
	return nil
}

func checkRightCell(missedShotsMap, woundedCellsMap map[int]map[int]bool, currentCell Cell) *Cell {
	tempCoordinate := currentCell.X + 1
	if tempCoordinate < BoardRow && !woundedCellsMap[tempCoordinate][currentCell.Y] && !missedShotsMap[tempCoordinate][currentCell.Y] {
		return &Cell{
			X: tempCoordinate,
			Y: currentCell.Y,
		}
	}
	return nil
}
