package toplib

import (
	ui "github.com/gizak/termui"
)

type View struct {
	rowFn    func() []*ui.Row
	Handlers map[string]func(ui.Event)
}

func NewView(fn func() []*ui.Row) *View {
	return &View{
		rowFn:    fn,
		Handlers: map[string]func(ui.Event){},
	}
}

func (v *View) Rows() []*ui.Row { return v.rowFn() }

type ViewMap struct {
	view  string
	views map[string]*View
}

func NewViewMap() *ViewMap {
	return &ViewMap{
		views: map[string]*View{},
	}
}

func (v *ViewMap) Add(name string, view *View) {
	v.views[name] = view
}

func (v *ViewMap) Set(name string) {
	v.view = name
}

func (v ViewMap) Get() *View {
	if view, ok := v.views[v.view]; ok {
		return view
	}
	return nil
}
