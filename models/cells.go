package models

// Cell is part of battleships
type Cell struct {
	X      int  `json:"x"`       // X coordinate
	Y      int  `json:"y"`       // Y coordinate
	IsDead bool `json:"is_dead"` // indicates single ship cell status (True = Dead, false = Alive)
}

// Compare compares if two cells has same
// coordinates
func (c Cell) Compare(input *Cell) bool {
	if c.X == input.X && c.Y == input.Y {
		return true
	}
	return false
}

func (c Cell) IsValid() bool {
	if c.X > BoardRow-1 || c.X < 0 || c.Y > BoardRow-1 || c.Y < 0 {
		return false
	}
	return true
}
