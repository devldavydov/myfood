package cmdproc

import (
	"bytes"
	"fmt"
	"text/template"
)

type ChardData struct {
	XLabels []string
	Data    []float64
	Label   string
	Type    string
}

func GetChartSnippet(data *ChardData) (string, error) {
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
						{{- range .Data }}
							{{- . }},
						{{- end }}
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
