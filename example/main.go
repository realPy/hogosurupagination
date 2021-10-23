package main

import (
	"fmt"

	"github.com/realPy/hogosurupagination"

	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/document"
	"github.com/realPy/hogosuru/documentfragment"
	"github.com/realPy/hogosuru/event"
	"github.com/realPy/hogosuru/hogosurudebug"
	"github.com/realPy/hogosuru/htmlanchorelement"
	"github.com/realPy/hogosuru/htmlelement"
	"github.com/realPy/hogosuru/htmltemplateelement"
	"github.com/realPy/hogosuru/node"
	"github.com/realPy/hogosuru/promise"
)

type GlobalContainer struct {
	parentNode node.Node
	pagination hogosurupagination.Pagination
	page       int
}

var template htmltemplateelement.HtmlTemplateElement

func (w *GlobalContainer) OnLoad(d document.Document, n node.Node, route string) (*promise.Promise, []hogosuru.Rendering) {

	w.parentNode = n
	htmltemplateelement.GetInterface()
	documentfragment.GetInterface()
	htmlanchorelement.GetInterface()

	w.page = 6
	w.pagination.IDPatternElem = "item-pattern"
	w.pagination.IDTemplate = "pagination-tpl"

	w.pagination.OnConfigureItem = func(elem htmlelement.HtmlElement, page int) {

		if link, err := elem.QuerySelector("#link-pattern"); hogosuru.AssertErr(err) {

			if aobj, err := link.Discover(); hogosuru.AssertErr(err) {

				if a, ok := aobj.(htmlanchorelement.HtmlAnchorElement); ok {

					if page >= 0 {

						a.SetTextContent(fmt.Sprintf("%d", page+1))

						a.OnClick(func(e event.Event) {
							w.pagination.Select(elem, page)

							e.PreventDefault()

						})
					} else {
						a.SetTextContent("...")
					}

				}
			}

		}

		w.pagination.OnSelectItem = func(elem htmlelement.HtmlElement) {
			class, _ := elem.ClassName()
			elem.SetClassName(class + " selected")
		}

	}

	return nil, []hogosuru.Rendering{&w.pagination}
}

func (w *GlobalContainer) Node(r hogosuru.Rendering) node.Node {

	if r == &w.pagination {
		if d, err := document.New(); hogosuru.AssertErr(err) {

			if elem, err := d.GetElementById("pagination"); hogosuru.AssertErr(err) {
				return elem.Node
			}

		}
	}

	return w.parentNode
}

func (w *GlobalContainer) OnEndChildRendering(r hogosuru.Rendering) {

}

func (w *GlobalContainer) OnEndChildsRendering() {
	w.pagination.SetMax(w.page)
}

func (w *GlobalContainer) OnUnload() {

}

func main() {

	hogosuru.Init()
	hogosurudebug.EnableDebug()
	hogosuru.Router().DefaultRendering(&GlobalContainer{})
	hogosuru.Router().Start(hogosuru.HASHROUTE)
	ch := make(chan struct{})
	<-ch

}
