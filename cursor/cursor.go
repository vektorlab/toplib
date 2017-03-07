package cursor

import "github.com/vektorlab/toplib"

// Cursor stores the currently selected Item.
type Cursor struct {
	ID string
}

func NewCursor() *Cursor {
	return &Cursor{}
}

func (c *Cursor) IDX(items []toplib.Item) int {
	for n, item := range items {
		if item.ID() == c.ID {
			return n
		}
	}
	return 0
}

func (c *Cursor) Up(items []toplib.Item) bool {
	idx := c.IDX(items)
	if idx > 0 {
		c.ID = items[idx-1].ID()
		return true
	}
	return false
}

func (c *Cursor) Down(items []toplib.Item) bool {
	idx := c.IDX(items)
	if idx < (len(items) - 1) {
		c.ID = items[idx+1].ID()
		return true
	}
	return false
}
