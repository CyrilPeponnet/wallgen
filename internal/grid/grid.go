package grid

import (
	"image"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/fogleman/gg"
)

// gridConfiguration is the configuration set by flags
type GridConfiguration struct {
	Palete            [][3]int
	BackGroundColor   [3]int
	X                 int
	Y                 int
	InnerDistribution int
	OuterDistribution int
	CellDistribution  int
	Step              int
}

// Grid represent a grid of cells
type Grid struct {
	canevas        *gg.Context
	final          *gg.Context
	Cells          []Cell
	X              int
	Y              int
	Step           int
	xCells, yCells int
}

// A Cell
type Cell struct {
	Alive bool
	// Color in RGBA
	InnerColor, OuterColor [4]int
	// Postion of the cell within the grid
	X, Y int
}

// Return a new grid
func NewGrid(cfg GridConfiguration) *Grid {

	// compute how many full cells we have for a given canvas
	// on x axis
	xCells := int(cfg.X/cfg.Step) - 1
	yCells := int(cfg.Y/cfg.Step) - 1

	g := Grid{
		Cells:   make([]Cell, xCells*yCells),
		X:       cfg.X,
		Y:       cfg.Y,
		Step:    cfg.Step,
		xCells:  xCells,
		yCells:  yCells,
		canevas: gg.NewContext(cfg.X, cfg.Y),
		final:   gg.NewContext(cfg.X, cfg.Y),
	}

	// Create our cells
	row := 0
	for idx := range g.Cells {
		if idx > 0 && idx%xCells == 0 {
			row++
		}
		// The X and Y starts at 1
		g.Cells[idx].X = idx - row*xCells
		g.Cells[idx].Y = row
	}

	// create the final rendering canevas
	g.final.SetRGB255(cfg.BackGroundColor[0], cfg.BackGroundColor[1], cfg.BackGroundColor[2])
	g.final.DrawRectangle(0, 0, float64(g.X), float64(g.Y))
	g.final.Fill()
	return &g

}

// Draw everything
func (g Grid) Draw() image.Image {

	// Draw our cells
	var cell Cell

	xMargin := (g.X - (g.xCells)*(g.Step) + g.Step/2) / 2
	yMargin := (g.Y - (g.yCells)*(g.Step) + g.Step/2) / 2

	g.canevas.SetLineWidth(2)

	// Draw the cells
	for cidx := range g.Cells {
		cell = g.Cells[cidx]
		if !cell.Alive {
			continue
		}

		if cell.OuterColor[3] > 0 {
			g.canevas.SetRGBA255(int(cell.OuterColor[0]), int(cell.OuterColor[1]), int(cell.OuterColor[2]), int(cell.OuterColor[3]))
			g.canevas.DrawCircle(float64(xMargin+cell.X*g.Step), float64(yMargin+cell.Y*g.Step), 25)
			g.canevas.Stroke()
		}

		if cell.InnerColor[3] > 0 {
			g.canevas.SetRGBA255(int(cell.InnerColor[0]), int(cell.InnerColor[1]), int(cell.InnerColor[2]), int(cell.InnerColor[3]))
			g.canevas.DrawCircle(float64(xMargin+cell.X*g.Step), float64(yMargin+cell.Y*g.Step), 10)
			g.canevas.Fill()

		}
	}

	// store the original source of light to draw it back later
	original := g.canevas.Image()

	// bloom this source of light
	// create a larger image
	size := g.canevas.Image().Bounds().Size()
	newSize := image.Rect(0, 0, size.X+10, size.Y+10)

	// copy the original in this larger image, slightly translated to the center
	extended := image.NewRGBA(newSize)
	xOffset := 10
	yOffset := 10
	extendedSize := g.canevas.Image().Bounds().Size()
	for x := 0; x < extendedSize.X; x++ {
		for y := 0; y < extendedSize.Y; y++ {
			extended.Set(xOffset+x, yOffset+y, g.canevas.Image().At(x, y))
		}
	}

	// dilate the image to have a bigger source of light
	dilated := effect.Dilate(extended, 2)

	// blur the image
	bloomed := blur.Gaussian(dilated, 4.0)

	// draw our bloomed light
	g.final.DrawImage(bloomed, 0, 0)

	// re-apply the original source of light
	g.final.DrawImage(original, 10, 10)

	return g.final.Image()
}

// Return a list of index of my neighbors
func (g Grid) CellNeighBoors(idx int) (nidx []int) {

	// Left
	if idx-1 >= 0 && idx%g.xCells != 0 {
		if g.Cells[idx-1].Alive {
			nidx = append(nidx, idx-1)
		}
	}
	// right
	if idx+1 < len(g.Cells) && (idx+1)%g.xCells != 0 {
		if g.Cells[idx+1].Alive {
			nidx = append(nidx, idx+1)
		}
	}
	// up
	if idx-g.xCells >= 0 {
		if g.Cells[idx-g.xCells].Alive {
			nidx = append(nidx, idx-g.xCells)
		}
	}
	// down
	if idx+g.xCells < len(g.Cells) {
		if g.Cells[idx+g.xCells].Alive {
			nidx = append(nidx, idx+g.xCells)
		}
	}
	// up left
	if idx-g.xCells-1 >= 0 && (idx)%g.xCells != 0 {
		if g.Cells[idx-g.xCells-1].Alive {
			nidx = append(nidx, idx-g.xCells-1)
		}
	}
	//down left
	if idx+g.xCells-1 < len(g.Cells) && (idx)%g.xCells != 0 {
		if g.Cells[idx+g.xCells-1].Alive {
			nidx = append(nidx, idx+g.xCells-1)
		}
	}
	// up right
	if idx-g.xCells+1 >= 0 && (idx+1)%g.xCells != 0 {
		if g.Cells[idx-g.xCells+1].Alive {
			nidx = append(nidx, idx-g.xCells+1)
		}
	}
	//down right
	if idx+g.xCells+1 < len(g.Cells) && (idx+1)%g.xCells != 0 {
		if g.Cells[idx+g.xCells+1].Alive {
			nidx = append(nidx, idx+g.xCells+1)
		}
	}

	return nidx
}
