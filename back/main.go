package main

import "fmt"

type Field struct {
	cells []*Cell
	turn  int
	size  int
}

func (f *Field) createCells() {
	f.cells = make([]*Cell, f.size)
	for i := 0; i < f.size; i++ {
		f.cells[i].i = i
	}
}

func NewField(size int) *Field {
	field := Field{size: size}
	field.createCells()
	for i := 0; i < field.size; i++ {
		field.cells[i].selectNeighbors(&field)
	}

	return &field
}

func (f *Field) run() {
	for i := 0; i < f.size; i++ {
		f.cells[i].setNewCondition()
	}
	for i := 0; i < f.size; i++ {
		f.cells[i].state = f.cells[i].tempState
	}
}

type Cell struct {
	state     bool
	tempState bool
	neighbors []*Cell
	i         int
}

func (c *Cell) setNewCondition() {
	aliveCount := 0

	for i := 0; i < len(c.neighbors); i++ {
		if c.neighbors[i].state {
			aliveCount++
		}
	}

	if aliveCount < 2 || aliveCount > 3 {
		c.tempState = false
		return
	}

	if aliveCount == 3 && !c.state {
		c.tempState = true
		return
	}
}

func (c *Cell) selectNeighbors(f *Field) {
	// top left
	if c.i == 0 {
		c.neighbors = []*Cell{f.cells[1], f.cells[f.size], f.cells[f.size+1]}
		return
	}
	// top right
	if c.i == f.size-1 {
		c.neighbors = []*Cell{f.cells[f.size-2], f.cells[f.size*2-1], f.cells[f.size*2-2]}
		return
	}
	// bottom left
	if c.i == f.size*(f.size-1) {
		c.neighbors = []*Cell{f.cells[f.size*(f.size-2)], f.cells[f.size*(f.size-2)+1], f.cells[f.size*(f.size-1)+1]}
		return
	}
	// bottom right
	if c.i == f.size*f.size-1 {
		c.neighbors = []*Cell{f.cells[f.size*f.size-2], f.cells[f.size*(f.size-1)-1], f.cells[f.size*(f.size-1)-2]}
		return
	}
	// top row no corners
	if c.i > 0 && c.i < f.size-1 {
		c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i+f.size-1], f.cells[c.i+f.size], f.cells[c.i+f.size+1]}
		return
	}
	// bottom row no corners
	if c.i > f.size*(f.size-1) && c.i < f.size*f.size-1 {
		c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i-f.size-1], f.cells[c.i-f.size], f.cells[c.i-f.size+1]}
		return
	}
	// left column no corners
	if c.i%f.size == 0 {
		c.neighbors = []*Cell{f.cells[c.i-f.size], f.cells[c.i+f.size], f.cells[c.i-f.size+1], f.cells[c.i+1], f.cells[c.i+f.size+1]}
		return
	}
	// right column no corners
	if c.i%f.size == 0 {
		c.neighbors = []*Cell{f.cells[c.i-f.size], f.cells[c.i+f.size], f.cells[c.i-f.size-1], f.cells[c.i-1], f.cells[c.i+f.size-1]}
		return
	}
	// all the rest with 8 normal neighbors
	c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i-f.size-1], f.cells[c.i-f.size], f.cells[c.i-f.size+1], f.cells[c.i+f.size-1], f.cells[c.i+f.size], f.cells[c.i+f.size+1]}
}

func main() {
	f := NewField(3)

	fmt.Println(f)
}
