package cmdproc

import (
	"bytes"
	"text/template"
)

type ChardData struct {
	ElemID  string
	XLabels []string
	Data    []float64
	Label   string
	Type    string
}

func GetChartSnippet(data *ChardData) (string, error) {
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
	`))

	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
