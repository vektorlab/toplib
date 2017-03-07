package section

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/cursor"
	"sort"
)

// Debug is a section that displays debug information
// which may be useful for developing toplib.
type Debug struct {
	Cursor *cursor.Cursor
}

func NewDebug() *Debug {
	return &Debug{
		Cursor: cursor.NewCursor(),
	}
}

func (d Debug) Name() string { return "debug" }

func (d Debug) Handlers(opts toplib.Options) map[string]func(ui.Event) {
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
	}
}

func (d *Debug) Grid(opts toplib.Options) *ui.Grid {
	p := ui.NewPar(fmt.Sprintf("Samples Loaded: %d", opts.Recorder.Counter))
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

func listHandlers() []string {
	strs := []string{}
	for path, _ := range ui.DefaultEvtStream.Handlers {
		strs = append(strs, path)
	}
	sort.Strings(strs)
	return strs
}
