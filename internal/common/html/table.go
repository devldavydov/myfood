package html

import (
	"fmt"
	"strings"
)

type Table struct {
	header []string
	rows   []*Tr
	footer []IELement
}

var _ IELement = (*Table)(nil)

func NewTable(header []string) *Table {
	return &Table{header: header}
}

func (r *Table) AddRow(row *Tr) *Table {
	r.rows = append(r.rows, row)
	return r
}

func (r *Table) AddFooterElement(elem IELement) *Table {
	r.footer = append(r.footer, elem)
	return r
}

func (r *Table) Build() string {
	var sb strings.Builder

	// Header
	sb.WriteString(`
	<table class="table table-bordered table-hover">
		<thead class="table-light">
			<tr>
	`)
	for _, h := range r.header {
		sb.WriteString(fmt.Sprintf("<th>%s</th>", h))
	}
	sb.WriteString(`
			</tr>
		</thead>
		<tbody>
	`)

	// Rows
	for _, row := range r.rows {
		sb.WriteString(row.Build())
	}
	sb.WriteString(`
		</tbody>
	`)

	// Footer
	sb.WriteString(`
		<tfoot>
	`)
	for _, item := range r.footer {
		sb.WriteString(item.Build())
	}
	sb.WriteString(`
		</tfoot>
	`)

	// End
	sb.WriteString(`
	</table>
	`)

	return sb.String()
}

//
//
//

type Td struct {
	val   IELement
	attrs Attrs
}

var _ IELement = (*Td)(nil)

func NewTd(val IELement, attrs Attrs) *Td {
	return &Td{val: val, attrs: attrs}
}

func (r *Td) Build() string {
	return fmt.Sprintf("<td %s>%s</td>", r.attrs.String(), r.val.Build())
}

//
//
//

type Tr struct {
	items []*Td
	attrs Attrs
}

var _ IELement = (*Tr)(nil)

func NewTr(attr Attrs) *Tr {
	return &Tr{attrs: attr}
}

func (r *Tr) AddTd(td *Td) *Tr {
	r.items = append(r.items, td)
	return r
}

func (r *Tr) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
	<tr %s>
	`, r.attrs.String()))

	for _, item := range r.items {
		sb.WriteString(item.Build())
	}

	sb.WriteString(`
	</tr>
	`)

	return sb.String()
}
