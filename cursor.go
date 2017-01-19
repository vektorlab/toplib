package toplib

// Cursor stores the currently selected Sample
type Cursor struct {
	ID string
}

func NewCursor() *Cursor {
	return &Cursor{}
}

func (c *Cursor) IDX(samples []*Sample) int {
	for n, sample := range samples {
		if sample.ID() == c.ID {
			return n
		}
	}
	return 0
}

func (c *Cursor) Up(samples []*Sample) bool {
	idx := c.IDX(samples)
	if idx > 0 {
		c.ID = samples[idx-1].ID()
		return true
	}
	return false
}

func (c *Cursor) Down(samples []*Sample) bool {
	idx := c.IDX(samples)
	if idx < (len(samples) - 1) {
		c.ID = samples[idx+1].ID()
		return true
	}
	return false
}
