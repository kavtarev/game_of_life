package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Field struct {
	cells []*Cell
	turn  int
	size  int
	total int
}

func (f *Field) createCells() {
	f.cells = make([]*Cell, f.total)
	for i := 0; i < f.total; i++ {
		f.cells[i] = &Cell{i: i}
	}
}

func NewField(size int) *Field {
	field := Field{size: size, total: size * size}
	field.createCells()
	for i := 0; i < len(field.cells); i++ {
		field.cells[i].selectNeighbors(&field)
	}

	return &field
}

func (f *Field) update(str string) error {
	if len(str) != f.total {
		return errors.New("invalid string")
	}

	for i := 0; i < len(str); i++ {
		f.cells[i].state = str[i] == '1'
	}

	return nil
}

func (f *Field) getStateString() string {
	b := strings.Builder{}
	b.Grow(f.total)

	for i := 0; i < len(f.cells); i++ {
		if f.cells[i].state {
			b.WriteString("1")
			continue
		}
		b.WriteString("0")
	}

	return b.String()
}

func (f *Field) run() {
	for i := 0; i < f.total; i++ {
		f.cells[i].setNewCondition()
	}
	for i := 0; i < f.total; i++ {
		f.cells[i].state = f.cells[i].tempState
	}
	f.turn++
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

	// fmt.Println(c.i, aliveCount)

	if aliveCount < 2 || aliveCount > 3 {
		c.tempState = false
		return
	}

	if aliveCount == 3 {
		c.tempState = true
		return
	}

	if aliveCount == 2 && c.state {
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
	if c.i == f.total-1 {
		c.neighbors = []*Cell{f.cells[f.total-2], f.cells[f.size*(f.size-1)-1], f.cells[f.size*(f.size-1)-2]}
		return
	}
	// top row no corners
	if c.i > 0 && c.i < f.size-1 {
		c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i+f.size-1], f.cells[c.i+f.size], f.cells[c.i+f.size+1]}
		return
	}
	// bottom row no corners
	if c.i > f.size*(f.size-1) && c.i < f.total-1 {
		c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i-f.size-1], f.cells[c.i-f.size], f.cells[c.i-f.size+1]}
		return
	}
	// left column no corners
	if c.i%f.size == 0 {
		c.neighbors = []*Cell{f.cells[c.i-f.size], f.cells[c.i+f.size], f.cells[c.i-f.size+1], f.cells[c.i+1], f.cells[c.i+f.size+1]}
		return
	}
	// right column no corners
	if (c.i+1)%f.size == 0 {
		c.neighbors = []*Cell{f.cells[c.i-f.size], f.cells[c.i+f.size], f.cells[c.i-f.size-1], f.cells[c.i-1], f.cells[c.i+f.size-1]}
		return
	}
	// all the rest with 8 neighbors
	c.neighbors = []*Cell{f.cells[c.i-1], f.cells[c.i+1], f.cells[c.i-f.size-1], f.cells[c.i-f.size], f.cells[c.i-f.size+1], f.cells[c.i+f.size-1], f.cells[c.i+f.size], f.cells[c.i+f.size+1]}
}

func main() {
	regexp.MustCompile(`\D`)

	f := NewField(3)

	f.update("000111000")
	fmt.Println(f.getStateString())
	f.run()
	fmt.Println(f.getStateString())
	f.run()
	fmt.Println(f.getStateString())

}
