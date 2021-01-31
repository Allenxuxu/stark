package limit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Test_updateEstimatedLimit(t *testing.T) {
	minlimit := int64(100)
	maxlimit := int64(10000)

	exceptCPU := float64(30)
	cpus := []float64{0, 10, 10, 20, 80, 60, 24, 50, 10, 10, 20, 80, 10, 5, 1, 2, 60, 24, 50, 60, 24, 50, 30, 30}
	limits := make([]float64, len(cpus))
	limits[0] = float64(updateEstimatedLimit(exceptCPU, cpus[0], int64(limits[0]), minlimit, maxlimit))

	for i := 1; i < len(cpus); i++ {
		limits[i] = float64(updateEstimatedLimit(exceptCPU, cpus[i], int64(limits[i-1]), minlimit, maxlimit))
	}

	var cpu, limit plotter.XYs
	for i := 0; i < len(cpus); i++ {
		cpu = append(cpu, plotter.XY{
			X: float64(i),
			Y: cpus[i],
		})

		limit = append(limit, plotter.XY{
			X: float64(i),
			Y: limits[i],
		})
	}

	assert.Nil(t, draw("cpu", cpu))
	assert.Nil(t, draw("limit", limit))
}

func draw(name string, vs ...interface{}) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = name
	p.X.Label.Text = "Time"
	p.Y.Label.Text = name

	err = plotutil.AddLinePoints(p, vs...)
	if err != nil {
		return err
	}

	return p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("/Users/xuxu/Downloads/%s.png", name))
}
