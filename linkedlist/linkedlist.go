package linkedlist

import "github.com/gookit/color"

type Cell struct {
	Row  int
	Col  int
	Vals []int
	Prev *Cell
	Next *Cell
}

type LinkedList struct {
	Head        *Cell
	Last        *Cell
	CurrentCell *Cell
}

func createlinkedList() *LinkedList {
	return &LinkedList{}
}

func (p *LinkedList) AddCell(row, col int, arr []int) error {
	c := &Cell{
		Row:  row,
		Col:  col,
		Vals: arr,
	}

	if p.Head == nil {
		p.Head = c
		c.Prev = nil
	} else {
		currentCell := p.Head
		for currentCell.Next != nil {
			currentCell = currentCell.Next
		}
		currentCell.Next = c
		c.Prev = currentCell
		p.Last = c
	}
	return nil
}

func (p *LinkedList) ShowAllEmptyCells() error {
	currentCell := p.Head
	if currentCell == nil {
		color.Red.Println("EmptyCell list is empty.")
		return nil
	}
	color.Green.Printf("%+v\n", *currentCell)
	for currentCell.Next != nil {
		currentCell = currentCell.Next
		color.Green.Printf("%+v\n", *currentCell)
	}
	return nil
}

func (p *LinkedList) CountNodes() int {
	count := 0
	currN := p.Head
	for currN != nil {
		currN = currN.Next
		count += 1
	}
	return count
}
