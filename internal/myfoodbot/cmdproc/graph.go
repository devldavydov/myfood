package cmdproc

import (
	"bytes"
	"text/template"
)

const (
	ChartColorRed    = "rgb(255, 99, 132)"
	ChartColorOrange = "rgb(255, 159, 64)"
	ChartColorYellow = "rgb(255, 205, 86)"
	ChartColorGreen  = "rgb(75, 192, 192)"
	ChartColorBlue   = "rgb(54, 162, 235)"
	ChartColorPurple = "rgb(153, 102, 255)"
	ChartColorGrey   = "rgb(201, 203, 207)"
)

type ChartData struct {
	ElemID   string
	XLabels  []string
	Type     string
	Datasets []ChartDataset
}

type ChartDataset struct {
	Data  []float64
	Label string
	Color string
}

func GetChartSnippet(data *ChartData) (string, error) {
	tmpl := template.Must(template.
		New("").
		Parse(`
<script>
	function plot() {
		const ctx = document.getElementById('{{.ElemID}}');

		new Chart(ctx, {
			type: '{{.Type}}',
			data: {
				labels: [
				{{- range .XLabels }}
					'{{- . }}',
				{{- end }}
				],
				datasets: [
				{{- range .Datasets }}
					{
						label: '{{.Label}}',
						data: [
						{{- range .Data }}
							{{- . }},
						{{- end }}
						],
						borderWidth: 2,
						borderColor: '{{.Color}}',
						backgroundColor: '{{.Color}}'
					},
				{{- end}}					
				]
			}
		});		
	}
	window.onload = plot;
</script>
	`))

	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
