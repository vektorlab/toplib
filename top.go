package toplib

import (
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/extra"
)

/*
Top renders periodically updated
samples of data in a top-like fashion.
*/
type Top struct {
	Samples  chan []*Sample // Incoming Samples
	Exit     chan bool
	Recorder *Recorder // Holds samples
	Sections []Section
	Tabpane  *extra.Tabpane
	Grid     *ui.Grid
	Handlers []map[string]func(ui.Event)
	section  int
}

func NewTop(sections []Section) *Top {
	top := &Top{
		Samples:  make(chan []*Sample),
		Exit:     make(chan bool),
		Recorder: NewRecorder(),
		Sections: sections,
		Tabpane:  extra.NewTabpane(),
		Grid:     ui.NewGrid(),
		Handlers: []map[string]func(ui.Event){},
	}
	return top
}

func (t *Top) handlers() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		t.Exit <- true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		t.section = t.Tabpane.SetActiveLeft()
		t.Render()
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		t.section = t.Tabpane.SetActiveRight()
		t.Render()
	})
	for path, fn := range t.Sections[t.section].Handlers(t) {
		ui.Handle(path, fn)
	}
}

func (t *Top) Render() {
	t.handlers()
	tabs := []extra.Tab{}
	for _, section := range t.Sections {
		grid := section.Grid(t)
		grid.Width = ui.TermWidth()
		grid.Align()
		tab := extra.NewTab(section.Name())
		tab.AddBlocks(grid)
		tabs = append(tabs, *tab)
	}
	t.Tabpane.SetTabs(tabs...)
	t.Tabpane.Width = ui.TermWidth()
	t.Tabpane.Align()
	ui.Clear()
	ui.Render(t.Tabpane)
}

func Run(top *Top) (err error) {
	if err = ui.Init(); err != nil {
		return err
	}
	defer ui.Close()
	go func() {
		for samples := range top.Samples {
			top.Recorder.Load(samples)
			top.Render()
		}
		ui.StopLoop()
	}()
	top.Render()
	ui.Loop()
	return err
}
