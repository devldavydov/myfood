package html

import (
	"fmt"
	"strings"
)

type Accordion struct {
	id    string
	items []*AccordionItem
}

var _ IELement = (*Accordion)(nil)

func NewAccordion(id string) *Accordion {
	return &Accordion{id: id}
}

func (r *Accordion) AddItem(item *AccordionItem) *Accordion {
	item.setAccordionID(r.id)
	r.items = append(r.items, item)
	return r
}

func (r *Accordion) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
	<div class="accordion" id="%s">
	`, r.id))

	for _, item := range r.items {
		sb.WriteString(item.Build())
	}

	sb.WriteString(`
	</div>
	`)

	return sb.String()
}

//
//
//

type AccordionItem struct {
	accordionID string
	id          string
	header      string
	body        IELement
}

var _ IELement = (*AccordionItem)(nil)

func HewAccordionItem(id, header string, body IELement) *AccordionItem {
	return &AccordionItem{id: id, header: header, body: body}
}

func (r *AccordionItem) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
	<div class="accordion-item">
		<h2 class="accordion-header">
			<button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#%s"
					aria-expanded="false" aria-controls="%s">
				<b>%s</b>
			</button>
		</h2>
		<div id="%s" class="accordion-collapse collapse" data-bs-parent="#%s">
			<div class="accordion-body">	
	`,
		r.id,
		r.id,
		r.header,
		r.id,
		r.accordionID))

	sb.WriteString(r.body.Build())

	sb.WriteString(`
			</div>
		</div>
	</div>
	`)

	return sb.String()
}

func (r *AccordionItem) setAccordionID(id string) {
	r.accordionID = id
}
