package section

import (
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/cursor"
	"github.com/vektorlab/toplib/toggle"
)

type Samples struct {
	Title    string
	Fields   []string
	Cursor   *cursor.Cursor
	SortMenu *toplib.Menu
	Toggles  toggle.Toggles
}

func NewSamples(title string, fields ...string) *Samples {
	return &Samples{
		Title:   title,
		Fields:  fields,
		Cursor:  cursor.NewCursor(),
		Toggles: toggle.NewToggles(&toggle.Toggle{Name: "sort"}),
	}
}

func (d Samples) Name() string { return d.Title }

func (d Samples) Handlers(opts toplib.Options) map[string]func(ui.Event) {
	return map[string]func(ui.Event){
		"/sys/kbd/<up>": func(ui.Event) {
			if d.Cursor.Up(opts.Recorder.Items()) {
				opts.Render()
			}
		},
		"/sys/kbd/<down>": func(ui.Event) {
			if d.Cursor.Down(opts.Recorder.Items()) {
				opts.Render()
			}
		},
		"/sys/kbd/s": func(ui.Event) {
			d.Toggles.Toggle("sort", true)
			opts.Render()
		},
	}
}

func (s Samples) Grid(opts toplib.Options) *ui.Grid {
	samples := opts.Recorder.Samples()
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
	table.BgColors[s.Cursor.IDX(opts.Recorder.Items())] = ui.ColorRed
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
