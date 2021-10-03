package hogosurupagination

import (
	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/document"
	"github.com/realPy/hogosuru/documentfragment"
	"github.com/realPy/hogosuru/element"
	"github.com/realPy/hogosuru/htmldivelement"
	"github.com/realPy/hogosuru/htmlelement"
	"github.com/realPy/hogosuru/htmllielement"
	"github.com/realPy/hogosuru/htmltemplateelement"
	"github.com/realPy/hogosuru/node"
	"github.com/realPy/hogosuru/promise"
)

type Pagination struct {
	IDPatternElem   string
	IDTemplate      string
	parentNode      node.Node
	template        htmltemplateelement.HtmlTemplateElement
	page            int
	current         int
	container       htmldivelement.HtmlDivElement
	OnConfigureItem func(elem htmlelement.HtmlElement, page int)
	OnSelectItem    func(elem htmlelement.HtmlElement)
}

func (p *Pagination) SetMax(page int) {
	p.page = page
	p.Refresh()

}
func (p *Pagination) Select(elem htmlelement.HtmlElement, page int) {
	p.current = page

	p.Refresh()

}

func (p *Pagination) additem(d document.Document, pattern element.Element, page int) {

	if clone, err := d.ImportNode(pattern.Node, true); hogosuru.AssertErr(err) {

		if elemfrom, ok := clone.(htmlelement.HtmlElementFrom); ok {

			var elem htmlelement.HtmlElement
			elem = elemfrom.HtmlElement()

			elem.RemoveAttribute("id")
			if p.OnConfigureItem != nil {
				p.OnConfigureItem(elem, page)
			}

			if p.current == page {

				if p.OnSelectItem != nil {
					p.OnSelectItem(elem)
				}
			}

			p.parentNode.AppendChild(elem.Node)
		}

	}
}
func (p *Pagination) Refresh() {

	if d, err := document.New(); hogosuru.AssertErr(err) {

		if fragment, err := p.template.Content(); hogosuru.AssertErr(err) {
			if linkpattern, err := fragment.GetElementById(p.IDPatternElem); hogosuru.AssertErr(err) {

				for r, err := p.parentNode.FirstChild(); err == nil; r, err = p.parentNode.FirstChild() {
					p.parentNode.RemoveChild(r)
				}

				var offset int = 2

				if p.page > 1 {
					p.additem(d, linkpattern, 0)
				}

				if p.current-offset > 1 {
					p.additem(d, linkpattern, -1)
				}

				var a int
				a = p.current - offset
				if a < 1 {
					a = 1
				}

				var max int

				max = p.current

				if max == (p.page - 1) {
					max = p.page - 2
				}

				for i := a; i <= max; i++ {

					p.additem(d, linkpattern, i)
				}

				a = p.current + offset
				if a > p.page-2 {
					a = p.page - 2
				}

				for i := p.current + 1; i <= a; i++ {

					p.additem(d, linkpattern, i)
				}

				if a < (p.page - 2) {

					p.additem(d, linkpattern, -1)
				}

				p.additem(d, linkpattern, (p.page - 1))

			}
		}

	}

}

func (p *Pagination) OnLoad(d document.Document, n node.Node, route string) (*promise.Promise, []hogosuru.Rendering) {

	p.parentNode = n

	htmltemplateelement.GetInterface()
	documentfragment.GetInterface()
	htmllielement.GetInterface()

	if elem, err := d.GetElementById(p.IDTemplate); hogosuru.AssertErr(err) {

		if elem, err := elem.Discover(); hogosuru.AssertErr(err) {

			if t, ok := elem.(htmltemplateelement.HtmlTemplateElement); ok {
				p.template = t
			}
		}

	}

	return nil, nil
}

func (p *Pagination) OnEndChildsRendering() {

}

func (p *Pagination) OnEndChildRendering(r hogosuru.Rendering) {

}

func (p *Pagination) Node(r hogosuru.Rendering) node.Node {

	return p.parentNode
}

func (p *Pagination) OnUnload() {

}
