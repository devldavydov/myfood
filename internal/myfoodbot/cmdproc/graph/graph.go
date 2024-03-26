package graph

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type DataPoint struct {
	Value float64
	Title string
}

func NewLine(title, xtitle, ytitle string, points []DataPoint) (io.Reader, error) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xtitle
	p.Y.Label.Text = ytitle

	vals := make(plotter.XYs, len(points))
	nominals := make([]string, len(points))
	min := math.MaxFloat64
	max := float64(0)

	for i := range points {
		vals[i].Y = points[i].Value
		vals[i].X = float64(i)

		if i == 0 || i == len(points)-1 {
			nominals[i] = points[i].Title
		} else {
			nominals[i] = ""
		}

		if points[i].Value > max {
			max = points[i].Value
		}
		if points[i].Value < min {
			min = points[i].Value
		}
	}
	p.NominalX(nominals...)

	ticks := []plot.Tick{
		{Value: min, Label: fmt.Sprintf("%.1f", min)},
		{Value: max, Label: fmt.Sprintf("%.1f", max)},
	}
	p.Y.Tick.Marker = plot.ConstantTicks(ticks)

	// points
	ln, pts, err := plotter.NewLinePoints(vals)
	if err != nil {
		return nil, err
	}
	ln.LineStyle.Width = vg.Length(1)
	ln.Color = plotutil.Color(0)

	p.Add(ln, pts)

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

type ChardData struct {
	XLabels []string
	Data    []float64
	Label   string
	Type    string
}

const _jsChartURL = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"

func GeChartSnippet(data *ChardData) (string, error) {
	tmpl := template.Must(template.
		New("").
		Parse(fmt.Sprintf(`
<script src="%s"></script>
<script>
	function plot() {
		const ctx = document.getElementById('chart');

		new Chart(ctx, {
			type: '{{.Type}}',
			data: {
				labels: [
				{{- range .XLabels }}
					'{{- . }}',
				{{- end }}
				],
				datasets: [
					{
						label: '{{.Label}}',
						data: [
						{{ range .Data }}
							{{.}},
						{{ end }}
						],
						borderWidth: 2,
						borderColor: 'rgb(255, 99, 132)',
						backgroundColor: 'rgb(255, 99, 132)'
					}
				]
			}
		});		
	}
	window.onload = plot;
</script>
	`, _jsChartURL)))

	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
