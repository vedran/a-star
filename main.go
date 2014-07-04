package main

import (
	"container/list"
	"fmt"
	"math"
)

// Grid - 2D Array of cells
type Grid [][]*Cell

type CellState int

const (
	UNSEEN   = 0
	OPEN     = 1
	CLOSED   = 2
	DISABLED = 3
	PATH     = 4
)

// Cell - X, Y, H, G, state, parent
type Cell struct {
	X      int
	Y      int
	H      int
	G      int
	State  CellState
	Parent *Cell
}

func (cell *Cell) F() int {
	return cell.G + cell.H
}

func (cell *Cell) Walkable() bool {
	return cell.State == UNSEEN || cell.State == OPEN
}

func calcHeuristic(curX int, curY int, targetX int, targetY int) int {
	// Manhattan
	return int(10*math.Abs(float64(curX-targetX)) + 10*math.Abs(float64(curY-targetY)))
}

func getLowestFScoreElement(openList *list.List) *list.Element {
	if openList.Len() == 0 {
		return nil
	}

	lowestElem := openList.Front()

	for e := lowestElem.Next(); e != nil; e = e.Next() {
		cell := e.Value.(*Cell)

		if cell.F() < lowestElem.Value.(*Cell).F() {
			lowestElem = e
		}
	}

	return lowestElem
}

func GetNeighbourCells(grid Grid, cell *Cell) ([]*Cell, []int) {
	var neighbours [8]*Cell
	var costs [8]int
	neighbourCount := 0

	// left
	if cell.X > 0 && grid[cell.Y][cell.X-1].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y][cell.X-1]
		costs[neighbourCount] = 10

		neighbourCount++
	}

	// upper left
	if cell.X > 0 && cell.Y+1 < len(grid) && grid[cell.Y+1][cell.X-1].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y+1][cell.X-1]
		costs[neighbourCount] = 14

		neighbourCount++
	}

	// top
	if cell.Y+1 < len(grid) && grid[cell.Y+1][cell.X].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y+1][cell.X]
		costs[neighbourCount] = 10

		neighbourCount++
	}

	// top right
	if cell.Y+1 < len(grid) && cell.X+1 < len(grid[cell.Y+1]) && grid[cell.Y+1][cell.X].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y+1][cell.X+1]
		costs[neighbourCount] = 14

		neighbourCount++
	}

	// right
	if cell.X+1 < len(grid[cell.Y]) && grid[cell.Y][cell.X+1].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y][cell.X+1]
		costs[neighbourCount] = 10

		neighbourCount++
	}

	// bottom right
	if cell.X+1 < len(grid[cell.Y]) && cell.Y > 0 && grid[cell.Y][cell.X+1].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y-1][cell.X+1]
		costs[neighbourCount] = 14

		neighbourCount++
	}

	// bottom
	if cell.Y > 0 && grid[cell.Y-1][cell.X].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y-1][cell.X]
		costs[neighbourCount] = 10

		neighbourCount++
	}

	// bottom left
	if cell.X > 0 && cell.Y > 0 && grid[cell.Y-1][cell.X-1].Walkable() {
		neighbours[neighbourCount] = grid[cell.Y-1][cell.X-1]
		costs[neighbourCount] = 14

		neighbourCount++
	}

	return neighbours[:neighbourCount], costs[:neighbourCount]
}

func ProcessNeighbours(curCell *Cell, targetX int, targetY int, grid Grid, openCells *list.List) {
	neighbours, costs := GetNeighbourCells(grid, curCell)

	for n := range neighbours {
		newG := curCell.G + costs[n]

		if neighbours[n].State == OPEN && newG < neighbours[n].G {
			// If neighbour is already in the open list
			// then check if my G + cost to that node < its existing G,
			// and if so, update that neighbour and set parent to me
			neighbours[n].G = newG
			neighbours[n].Parent = curCell
		} else if neighbours[n].State == UNSEEN {
			// If my neighbour is not already on the open list, calculate G and H and add it to the open list
			neighbours[n].G = newG
			neighbours[n].H = calcHeuristic(neighbours[n].X, neighbours[n].Y, targetX, targetY)
			neighbours[n].State = OPEN
			neighbours[n].Parent = curCell

			openCells.PushBack(neighbours[n])
		}
	}
}

func PrintGrid(startX int, startY int, targetX int, targetY int, grid Grid) {
	for y := range grid {
		for x := range grid[y] {
			if x == startX && y == startY {
				fmt.Printf("[O] ")
			} else if x == targetX && y == targetY {
				fmt.Printf("[X] ")
			} else if grid[y][x].State == PATH {
				fmt.Printf("[*] ")
			} else if grid[y][x].State == DISABLED {
				fmt.Printf("[|] ")
			} else {
				fmt.Printf("[ ] ")
			}
		}

		fmt.Printf("\n")
	}
}

func main() {
	// Build the grid
	gridWidth := 7
	gridHeight := 5

	startX := 1
	startY := 2

	targetX := 5
	targetY := 2

	grid := make([][]*Cell, gridHeight)

	for y := range grid {
		grid[y] = make([]*Cell, gridWidth)

		for x := range grid[y] {
			h := calcHeuristic(x, y, targetX, targetY)
			grid[y][x] = &Cell{x, y, h, 0, UNSEEN, nil}
		}
	}

	// Make a wall
	grid[1][3].State = DISABLED
	grid[2][3].State = DISABLED
	grid[3][3].State = DISABLED

	// Init the starting cell
	startCell := grid[startY][startX]
	startCell.H = 0

	// Add the start cell to the list of open cells
	openCells := list.New()
	openCells.PushBack(startCell)

	for openCells.Len() > 0 {
		lowestElem := getLowestFScoreElement(openCells)
		if lowestElem == nil {
			panic(fmt.Sprintf("No lowest elem found"))
		}

		// Remove the lowest cost element of the open list
		openCells.Remove(lowestElem)

		curCell := lowestElem.Value.(*Cell)
		curCell.State = CLOSED

		if curCell.X == targetX && curCell.Y == targetY {
			for curCell != nil {
				curCell.State = PATH
				curCell = curCell.Parent
			}

			PrintGrid(startX, startY, targetX, targetY, grid)
			return
		}

		ProcessNeighbours(curCell, targetX, targetY, grid, openCells)
	}
}
