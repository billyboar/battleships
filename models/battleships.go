package models

import (
	"github.com/billyboar/battleships/helpers"
	"github.com/gofrs/uuid"
)

// BattleShipNumber is the total number
// of battleships in a single board
const BattleShipNumber = 3

// BattleShip lengths types
const (
	DestroyerLength  = 4
	BattleShipLength = 5
)

// BattleShip health status
const (
	Dead  = true
	Alive = false
)

// BattleShip has 2 types of battleships
// with lengths of 4 and 5
type BattleShip struct {
	ID         string `json:"id"`          // UUID
	Length     int    `json:"length"`      // they must be either 4 or 5
	IsVertical bool   `json:"is_vertical"` // indicates ship orientation
	Cells      []Cell `json:"cells"`
	IsDead     bool   `json:"is_dead"` // indicates if ships is dead or alive (true = Dead, false = Alive)
}

// NewBattleShip creates new battleship struct
func NewBattleShip(shipLength int) (*BattleShip, error) {
	ship := BattleShip{
		Length: shipLength,
	}

	shipID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	ship.ID = shipID.String()

	// generating random ship positions
	randomPosition := helpers.GenerateRandomInt(100)
	if randomPosition%2 == 0 {
		ship.IsVertical = true
	}

	return &ship, nil
}

// BuildBody adds head cell and rest of the cells depending on length
func (b *BattleShip) BuildBody(headCell Cell) {
	b.Cells = append(b.Cells, headCell)
	for i := 1; i < b.Length; i++ {
		if b.IsVertical {
			b.Cells = append(b.Cells, Cell{
				X: headCell.X,
				Y: headCell.Y + i,
			})
		} else {
			b.Cells = append(b.Cells, Cell{
				X: headCell.X + i,
				Y: headCell.Y,
			})
		}
	}
}

func (b *BattleShip) GetDamageCount() int {
	return len(b.GetDamagedCells())
}

func (b *BattleShip) GetDamagedCells() (damagedCells []Cell) {
	for _, cell := range b.Cells {
		if cell.IsDead {
			damagedCells = append(damagedCells, cell)
		}
	}
	return
}
