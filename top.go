package toplib

import (
	ui "github.com/gizak/termui"
)

// Top displays periodically updated
// samples of data in a top-like fashion
type Top struct {
	Samples  chan []*Sample // Incoming Samples
	Exit     chan bool
	Recorder *Recorder // Holds samples
	Views    *ViewMap  // Holds views
}

func NewTop() *Top {
	return &Top{
		Samples:  make(chan []*Sample),
		Exit:     make(chan bool),
		Recorder: NewRecorder(),
		Views:    NewViewMap(),
	}
}

func (t *Top) Render() {
	view := t.Views.Get()
	ui.Body.Rows = view.Rows()
	for path, fn := range view.Handlers {
		ui.Handle(path, fn)
	}
	Reset()
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

	ui.DefaultEvtStream.ResetHandlers()
	//ui.Handle("/timer/1s", func(ui.Event) {
	//	top.Render()
	//})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		top.Exit <- true
	})

	ui.Loop()

	return err
}

func Reset() {
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)
}
