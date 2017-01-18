// +build ignore

package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/kevinschoon/ctop"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type MockSource struct {
	samples []*ctop.Sample
}

func (m *MockSource) Collect() []*ctop.Sample {
	if m.samples == nil {
		m.samples = []*ctop.Sample{}
		for i := 0; i < 35; i++ {
			m.samples = append(m.samples, ctop.NewSample(RandString(10)))
		}
	}
	samples := []*ctop.Sample{}
	for _, sample := range m.samples {
		/*
			// Uncomment to change iteration size
			if rand.Intn(10) > 5 {
				continue
			}
		*/
		s := ctop.NewSample(sample.ID())
		s.SetFloat64("CPU", float64(rand.Intn(5)))
		s.SetFloat64("MEM", float64(rand.Intn(40)))
		s.SetFloat64("DISK", float64(rand.Intn(60)))
		s.SetFloat64("GPU", float64(rand.Intn(90)))
		s.SetString("THING", RandString(50))
		s.SetString("OTHER THING", RandString(20))
		samples = append(samples, s)
	}
	return samples
}

func main() {
	source := &MockSource{}
	top := ctop.NewTop()
	var (
		cursor  = ctop.NewCursor()
		toggles = ctop.NewToggles(
			&ctop.Toggle{Name: "sort"},
			&ctop.Toggle{Name: "expanded"},
		)
		sortMenu  = ctop.NewMenu("ID", "CPU", "MEM", "DISK", "GPU")
		gauges    = ctop.NewGauges("CPU", "MEM", "DISK")
		summary   = ctop.NewSummary()
		chartCPU  = ctop.NewChart("CPU")
		chartMem  = ctop.NewChart("MEM")
		chartDisk = ctop.NewChart("DISK")
		chartGPU  = ctop.NewChart("GPU")
		table     = ctop.NewTable("ID", "CPU", "MEM",
			"DISK", "GPU", "THING", "OTHER THING")
	)

	defaultView := ctop.NewView(func() []*ui.Row {
		return []*ui.Row{
			ctop.NewHeader().Row(),
			// Top two sections
			ui.NewRow(
				ui.NewCol(6, 0, gauges.Buffers(top.Recorder, cursor)...),
				ui.NewCol(6, 0, summary.Buffers(top.Recorder, cursor)...),
			),
			// Main section
			ui.NewRow(
				ui.NewCol(12, 0, table.Buffers(top.Recorder, cursor)...),
			),
			// Bottom toggles
			ui.NewRow(
				ui.NewCol(12, 0, toggles.Buffers()...),
			),
		}
	})

	defaultView.Handlers["/sys/kbd/<up>"] = func(ui.Event) {
		if cursor.Up(top.Recorder.Samples()) {
			top.Render()
		}
	}

	defaultView.Handlers["/sys/kbd/<down>"] = func(ui.Event) {
		if cursor.Down(top.Recorder.Samples()) {
			top.Render()
		}
	}

	defaultView.Handlers["/sys/kbd/s"] = func(ui.Event) {
		if toggles.Toggle("sort", true) {
			top.Views.Set("sort")
		} else {
			top.Views.Set("default")
		}
		top.Render()
	}

	defaultView.Handlers["/sys/kbd/x"] = func(ui.Event) {
		if toggles.Toggle("expanded", true) {
			top.Views.Set("expanded")
		} else {
			top.Views.Set("default")
		}
		top.Render()
	}

	expandedView := ctop.NewView(func() []*ui.Row {
		return []*ui.Row{
			ui.NewRow(
				ui.NewCol(3, 0, chartCPU.Buffers(top.Recorder, cursor)...),
				ui.NewCol(3, 0, chartMem.Buffers(top.Recorder, cursor)...),
				ui.NewCol(3, 0, chartDisk.Buffers(top.Recorder, cursor)...),
				ui.NewCol(3, 0, chartGPU.Buffers(top.Recorder, cursor)...),
			),
			ui.NewRow(
				ui.NewCol(12, 0, toggles.Buffers()...),
			),
		}
	})

	sortView := ctop.NewView(func() []*ui.Row {
		return []*ui.Row{
			ui.NewRow(
				ui.NewCol(6, 0, gauges.Buffers(top.Recorder, cursor)...),
				ui.NewCol(6, 0, summary.Buffers(top.Recorder, cursor)...),
			),
			ui.NewRow(
				ui.NewCol(3, 0, sortMenu),
				ui.NewCol(9, 0, table.Buffers(top.Recorder, cursor)...),
			),
			ui.NewRow(
				ui.NewCol(12, 0, toggles.Buffers()...),
			),
		}
	})

	sortView.Handlers["/sys/kbd/<up>"] = func(ui.Event) {
		top.Recorder.SortField = sortMenu.Up()
		top.Render()
	}
	sortView.Handlers["/sys/kbd/<down>"] = func(ui.Event) {
		top.Recorder.SortField = sortMenu.Down()
		top.Render()
	}

	top.Views.Add("default", defaultView)
	top.Views.Add("sort", sortView)
	top.Views.Add("expanded", expandedView)
	top.Views.Set("default")

	tick := time.NewTicker(100 * time.Millisecond)

	go func() {
	loop:
		for {
			select {
			case <-top.Exit:
				close(top.Samples)
				break loop
			case <-tick.C:
				top.Samples <- source.Collect()
			}
		}
		tick.Stop()
	}()

	if err := ctop.Run(top); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
