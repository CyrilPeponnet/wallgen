package main

import (
	"flag"
	"image/png"
	"math/rand"
	"os"
	"sort"

	"github.com/CyrilPeponnet/wallgen/internal/grid"
	"github.com/CyrilPeponnet/wallgen/internal/utils"
)

func main() {

	outPtr := flag.String("output", "output.png", "The output file")
	sizeX := flag.Int("x", 800, "The the x size of the image")
	sizeY := flag.Int("y", 800, "The the y size of the image")
	step := flag.Int("step", 60, "The step between each cells")
	cellDst := flag.Int("cells", 40, "The cell distribution on the grid (0-100)")
	innerDst := flag.Int("inner", 80, "The inner circle distribution on a cell (0-100)")
	outerDst := flag.Int("outer", 60, "The outer circle distribution on a cell (0-100)")
	palPtr := flag.String("palette", "#fa32f3:100,#00a1cd:100,#4b1ff7:30", "The palette passed as sorted weighed colors")
	bgptr := flag.String("backround", "#000e12", "The background color")
	flag.Parse()

	pal := utils.ParsePalette(*palPtr)

	parsedPal := [][3]int{}
	for p, w := range pal {
		col, err := utils.ParseHexColor(p)
		if err != nil {
			panic(err)
		}
		for i := 0; i < w; i++ {
			parsedPal = append(parsedPal, col)
		}
	}

	bg, err := utils.ParseHexColor(*bgptr)
	if err != nil {
		panic(err)
	}

	c := grid.GridConfiguration{
		X:                 *sizeX,
		Y:                 *sizeY,
		Palete:            parsedPal,
		InnerDistribution: *innerDst,
		OuterDistribution: *outerDst,
		CellDistribution:  *cellDst,
		BackGroundColor:   bg,
		Step:              *step,
	}

	g := grid.NewGrid(c)

	// randomly active cells
	for cidx := range g.Cells {
		if rand.Intn(100) < c.CellDistribution {
			g.Cells[cidx].Alive = true
		}
	}

	// active cell in cluster
	for cidx := range g.Cells {
		if len(g.CellNeighBoors(cidx)) > 4 {
			g.Cells[cidx].Alive = true
		}
	}

	// active cell in cluster
	for cidx := range g.Cells {
		if len(g.CellNeighBoors(cidx)) > 7 {
			g.Cells[cidx].Alive = true
		}
	}

	// Randomly assign colors and set alpha to neighbors weigh
	for cidx := range g.Cells {
		// pick a random color for inner circle
		w := len(g.CellNeighBoors(cidx))

		color := c.Palete[rand.Intn(len(c.Palete))]

		if rand.Intn(100) < c.InnerDistribution {
			g.Cells[cidx].InnerColor[0] = color[0]
			g.Cells[cidx].InnerColor[1] = color[1]
			g.Cells[cidx].InnerColor[2] = color[2]
			g.Cells[cidx].InnerColor[3] = 28*w - rand.Intn(c.InnerDistribution)
		}

		// Change color if we dont have enouygh neighbors
		if w < 2 {
			color = c.Palete[rand.Intn(len(c.Palete))]
		}

		if w > 1 && rand.Intn(100) < c.OuterDistribution {
			g.Cells[cidx].OuterColor[0] = color[0]
			g.Cells[cidx].OuterColor[1] = color[1]
			g.Cells[cidx].OuterColor[2] = color[2]
			g.Cells[cidx].OuterColor[3] = 28*w - rand.Intn(c.OuterDistribution)
		} else {
			// if we have no outer circle use the less used color instead for inner
			lv := 0
			lc := ""

			for mc, mv := range pal {
				if lv == 0 {
					lv = mv
					lc = mc
					continue
				}
				if mv < lv {
					lv = mv
					lc = mc
				}
			}

			color, err := utils.ParseHexColor(lc)
			if err != nil {
				panic(err)
			}

			g.Cells[cidx].InnerColor[0] = color[0]
			g.Cells[cidx].InnerColor[1] = color[1]
			g.Cells[cidx].InnerColor[2] = color[2]

		}
	}

	// if we have a neighbors that is w higher than us we can randomly take its colors
	for cidx := range g.Cells {

		if !g.Cells[cidx].Alive {
			continue
		}

		nb := g.CellNeighBoors(cidx)
		sort.Ints(nb)
		sort.Sort(sort.Reverse(sort.IntSlice(nb)))

		for _, n := range nb {
			if len(g.CellNeighBoors(n)) > len(nb) {

				// only if they are set
				if g.Cells[n].InnerColor[0] != 0 && g.Cells[n].InnerColor[1] != 0 && g.Cells[2].InnerColor[2] != 0 {

					g.Cells[cidx].InnerColor[0] = g.Cells[n].InnerColor[0]
					g.Cells[cidx].InnerColor[1] = g.Cells[n].InnerColor[1]
					g.Cells[cidx].InnerColor[2] = g.Cells[n].InnerColor[2]
				}

				if g.Cells[n].OuterColor[0] != 0 && g.Cells[n].OuterColor[1] != 0 && g.Cells[2].OuterColor[2] != 0 {

					g.Cells[cidx].OuterColor[0] = g.Cells[n].OuterColor[0]
					g.Cells[cidx].OuterColor[1] = g.Cells[n].OuterColor[1]
					g.Cells[cidx].OuterColor[2] = g.Cells[n].OuterColor[2]

					g.Cells[cidx].OuterColor[3] = g.Cells[n].InnerColor[3]
				}
				continue
			}
		}

	}

	// Write image
	im := g.Draw()
	file, err := os.Create(*outPtr)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, im); err != nil {
		panic(err)
	}
}
