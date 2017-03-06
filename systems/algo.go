//PathFinding Algorithm

package main

import "container/heap"
import "fmt"
import "math"
import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"fmt"
	"image/color"
	"sync"
)

var list [][]bool

type grid struct {
	x   int
	y   int
	f   int
	g   int
	par *grid
}

var start grid
var end grid

type gridHeap []grid

func (h gridHeap) Len() int          { return len(h) }
func (h gridHeap) Les(i, j int) bool { return h[i].f < h[j].f }
func (h gridHeap) Swap(i, j int)     { h[i], h[j] = h[j], h[i] }

func (h *gridHeap) Push(x interface{}) {
	*h = append(*h, x.(grid))
}

func (h *gridHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var h *gridHeap

func hvalue(x, y int) int {
	var a, b, h int
	X := endgrid.x
	Y := endgrid.y
	a = math.Abs(X - x)
	b = math.Abs(Y - y)
	if a > b {
		h = a
	} else {
		h = b
	}
	return h
}

func eval(b1, block *grid) bool {
	if b1.x == engrid.x && b1.y == endgrid.y {
		endgrid.par = b1
		return true
	}
	b1.g = block.g + 1
	b1.f = hvalue(b1.x, b1.y) + b1.g
	b1.par = block
	list[b1.x][b1.y] = true
	heap.Push(h, *b1)
	return false
}

func open(block grid) {
	var b1 grid
	if block.x-1 >= 0 && block.y-1 >= 0 {
		if !Grid[block.x-1][block.y-1] {
			if !list[block.x-1][block.y-1] {
				b1.x = block.x - 1
				b1.y = block.y - 1

				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.y-1 >= 0 {
		if !Grid[block.x][block.y-1] {
			if !list[block.x][block.y-1] {
				b1.x = block.x
				b1.y = block.y - 1
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize && block.y-1 >= 0 {
		if !Grid[block.x+1][block.y-1] {
			if !list[block.x+1][block.y-1] {
				b1.x = block.x + 1
				b1.y = block.y - 1
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.x-1 >= 0 {
		if !Grid[block.x-1][block.y] {
			if !list[block.x-1][block.y] {
				b1.x = block.x - 1
				b1.y = block.y
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize {
		if !Grid[block.x+1][block.y] {
			if !list[block.x+1][block.y] {
				b1.x = block.x + 1
				b1.y = block.y
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.x-1 >= 0 && block.y+1 < int(engo.GameHeight()*ScaleFactor)/GridSize {
		if !Grid[block.x-1][block.y+1] {
			if !list[block.x-1][block.y+1] {
				b1.x = block.x - 1
				b1.y = block.y + 1
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.y+1 < (engo.GameHeight()*ScaleFactor)/GridSize {
		if !Grid[block.x][block.y+1] {
			if !list[block.x][block.y+1] {
				b1.x = block.x
				b1.y = block.y + 1
				if eval(&b1, &block) {
					return
				}
			}
		}
	}
	if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize && block.y+1 < int(engo.GameHeight()*ScaleFactor)/GridSize {
		if !grid[block.x+1][block.y+1] {
			if !list[block.x+1][block.y+1] {
				b1.x = block.x + 1
				b1.y = block.y + 1
				if eval(&b1, &block) {
					return
				}
			}
		}
	}

	p := (*h)[0].x
	q := (*h)[0].y
	list[p][q] = true
	block = (*h)[0]
	heap.Pop(h)
	if len(h) <= 0 {
		return
	}
	open(block) //function calling recursively
}

func main() {

	var x1, y1 int
	var startgrid, temp grid

	x1 = startgrid.x
	y1 = startgrid.y
	h = &gridHeap{}
	heap.Init(h)
	startgrid.par = startgrid
	startgrid.g = 0
	list = make([][]bool, int(engo.WindowWidth()*ScaleFactor)/GridSize)
	for i, _ := range list {
		list[i] = make([]bool, int(engo.WindowHeight()*ScaleFactor)/GridSize)
	}
	list[x1][y1] = 1

	open(startgrid)

	var path []grid

	temp = endgrid
	path = append(path, temp)
	for temp.par != &startgrid {
		temp = *temp.par
		path = append(path, temp)
	}
	append(path, temp)
	path = Reverse(path)
}
