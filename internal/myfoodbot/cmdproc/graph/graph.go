package graph

import (
	"bytes"
	"io"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type DataPoint struct {
	Value float64
	Title string
}

func NewBarChart(title, ytitle string, points []DataPoint) (io.Reader, error) {
	p := plot.New()
	p.Title.Text = title
	p.Y.Label.Text = ytitle

	vals := make(plotter.Values, len(points))
	nominals := make([]string, len(points))
	for i := range points {
		vals[i] = points[i].Value
		if i == 0 || i == len(points)-1 {
			nominals[i] = points[i].Title
			continue
		}
		nominals[i] = ""
	}
	p.NominalX(nominals...)

	bars, err := plotter.NewBarChart(vals, vg.Points(20))
	if err != nil {
		return nil, err
	}

	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	p.Add(bars)

	buf := bytes.NewBuffer([]byte{})
	wr, err := p.WriterTo(vg.Points(640), vg.Points(480), "png")
	if err != nil {
		return nil, err
	}

	_, err = wr.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
