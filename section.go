package toplib

import (
	"fmt"
	ui "github.com/gizak/termui"
)

/*
Section returns a renderable ui.Grid. A Top
may contain an arbitrary number of Sections within.
*/
type Section interface {
	Name() string
	Grid(*Top) *ui.Grid
	Handlers(*Top) map[string]func(ui.Event)
}

type section struct {
	name string
}

func (s section) Name() string { return s.name }

type SamplesSection struct {
	section
	Fields   []string
	Cursor   *Cursor
	SortMenu *Menu
	Toggles  Toggles
}

func NewSamplesSection(fields ...string) *SamplesSection {
	return &SamplesSection{
		section: section{name: "samples"},
		Fields:  fields,
		Cursor:  NewCursor(),
		Toggles: NewToggles(&Toggle{Name: "sort"}),
	}
}

func (d SamplesSection) Handlers(t *Top) map[string]func(ui.Event) {
	return map[string]func(ui.Event){
		"/sys/kbd/<up>": func(ui.Event) {
			if d.Cursor.Up(t.Recorder.Items()) {
				t.Render()
			}
		},
		"/sys/kbd/<down>": func(ui.Event) {
			if d.Cursor.Down(t.Recorder.Items()) {
				t.Render()
			}
		},
		"/sys/kbd/s": func(ui.Event) {
			d.Toggles.Toggle("sort", true)
			t.Render()
		},
	}
}

func (s SamplesSection) Grid(t *Top) *ui.Grid {
	samples := t.Recorder.Samples()
	rows := [][]string{s.Fields}
	for _, sample := range samples {
		rows = append(rows, sample.Strings(s.Fields))
	}
	table := ui.NewTable()
	table.Rows = rows
	table.Separator = false
	table.Border = false
	table.SetSize()
	table.Analysis()
	table.BgColors[s.Cursor.IDX(t.Recorder.Items())] = ui.ColorRed
	l := ui.NewList()
	l.Items = s.Fields
	l.Height = 30
	l.Width = 25
	if s.Toggles.State("sort") {
		return ui.NewGrid(
			ui.NewRow(
				ui.NewCol(3, 0, l),
				ui.NewCol(9, 0, table),
			),
		)
	}
	return ui.NewGrid(
		ui.NewRow(
			ui.NewCol(12, 0, table),
		),
	)
}

type DebugSection struct {
	Cursor *Cursor
}

func NewDebugSection() *DebugSection {
	return &DebugSection{
		Cursor: NewCursor(),
	}
}

func (d DebugSection) Name() string { return "debug" }

func (d DebugSection) Handlers(t *Top) map[string]func(ui.Event) {
	return map[string]func(ui.Event){
		"/sys/kbd/<up>": func(ui.Event) {
			if d.Cursor.Up(t.Recorder.Items()) {
				t.Render()
			}
		},
		"/sys/kbd/<down>": func(ui.Event) {
			if d.Cursor.Down(t.Recorder.Items()) {
				t.Render()
			}
		},
	}
}

func (d *DebugSection) Grid(t *Top) *ui.Grid {
	p := ui.NewPar(fmt.Sprintf("Samples Loaded: %d", t.Recorder.Counter))
	p.Height = 3
	p.Width = 10
	c := ui.NewPar(fmt.Sprintf("Cursor Position: %s", d.Cursor.ID))
	c.Height = 3
	c.Width = 10
	l := ui.NewList()
	l.BorderLabel = "Handlers"
	l.Items = listHandlers()
	l.Width = 25
	l.Height = 20
	return ui.NewGrid(
		ui.NewRow(
			ui.NewCol(6, 0, p),
			ui.NewCol(6, 0, c),
		),
		ui.NewRow(
			ui.NewCol(12, 0, l),
		),
	)
}
