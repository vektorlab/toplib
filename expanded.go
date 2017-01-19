package toplib

/*
TODO: Need to hook these into existing "sections"
import (
	ui "github.com/gizak/termui"
)

type Expanded struct {
	Info *ui.Table
	Net  *ExpandedNet
	Cpu  *ExpandedCpu
	Mem  *ExpandedMem
}

func NewExpanded(id, name string) *Expanded {
	return &Expanded{
		Info: NewInfo(id, name),
		Net:  NewExpandedNet(),
		Cpu:  NewExpandedCpu(),
		Mem:  NewExpandedMem(),
	}
}

func NewInfo(id, name string) *ui.Table {
	p := ui.NewTable()
	p.Rows = [][]string{
		[]string{"name", name},
		[]string{"id", id},
	}
	p.Height = 4
	p.Width = 50
	p.FgColor = ui.ColorWhite
	p.Seperator = false
	return p
}

func (w *Expanded) Render() {
	ui.Render(w.Info, w.Cpu, w.Mem, w.Net)
	ui.Handle("/timer/1s", func(ui.Event) {
		ui.Render(w.Info, w.Cpu, w.Mem, w.Net)
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func (w *Expanded) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Mem),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Expanded) Highlight() {
}

func (w *Expanded) UnHighlight() {
}

func (w *Expanded) SetCPU(val int) {
	w.Cpu.Update(val)
}

func (w *Expanded) SetNet(rx int64, tx int64) {
	w.Net.Update(rx, tx)
}

func (w *Expanded) SetMem(val int64, limit int64, percent int) {
	w.Mem.Update(int(val), int(limit))
}

type ExpandedCpu struct {
	*ui.LineChart
	hist FloatHistData
}

func NewExpandedCpu() *ExpandedCpu {
	cpu := &ExpandedCpu{ui.NewLineChart(), NewFloatHistData(60)}
	cpu.BorderLabel = "CPU"
	cpu.Height = 10
	cpu.Width = 50
	cpu.X = 0
	cpu.Y = 4
	cpu.Data = cpu.hist.data
	cpu.DataLabels = cpu.hist.labels
	cpu.AxesColor = ui.ColorDefault
	cpu.LineColor = ui.ColorGreen
	return cpu
}

func (w *ExpandedCpu) Update(val int) {
	w.hist.Append(float64(val))
}

type ExpandedMem struct {
	*ui.BarChart
	hist IntHistData
}

func NewExpandedMem() *ExpandedMem {
	mem := &ExpandedMem{
		ui.NewBarChart(),
		NewIntHistData(8),
	}
	mem.BorderLabel = "MEM"
	mem.Height = 10
	mem.Width = 50
	mem.BarWidth = 5
	mem.BarGap = 1
	mem.X = 0
	mem.Y = 14
	mem.TextColor = ui.ColorDefault
	mem.Data = mem.hist.data
	mem.BarColor = ui.ColorGreen
	mem.DataLabels = mem.hist.labels
	mem.NumFmt = byteFormatInt
	return mem
}

func (w *ExpandedMem) Update(val int, limit int) {
	// implement our own scaling for mem graph
	if val*4 < limit {
		w.SetMax(val * 4)
	} else {
		w.SetMax(limit)
	}
	w.hist.Append(val)
}

type ExpandedNet struct {
	*ui.Sparklines
	rxHist DiffHistData
	txHist DiffHistData
}

func NewExpandedNet() *ExpandedNet {
	net := &ExpandedNet{ui.NewSparklines(), NewDiffHistData(50), NewDiffHistData(50)}
	net.BorderLabel = "NET"
	net.Height = 6
	net.Width = 50
	net.X = 0
	net.Y = 24

	rx := ui.NewSparkline()
	rx.Title = "RX"
	rx.Height = 1
	rx.Data = net.rxHist.data
	rx.TitleColor = ui.ColorDefault
	rx.LineColor = ui.ColorGreen

	tx := ui.NewSparkline()
	tx.Title = "TX"
	tx.Height = 1
	tx.Data = net.txHist.data
	tx.TitleColor = ui.ColorDefault
	tx.LineColor = ui.ColorYellow

	net.Lines = []ui.Sparkline{rx, tx}
	return net
}

func (w *ExpandedNet) Update(rx int64, tx int64) {
	var rate string

	w.rxHist.Append(int(rx))
	rate = strings.ToLower(byteFormatInt(w.rxHist.Last()))
	w.Lines[0].Title = fmt.Sprintf("RX [%s/s]", rate)

	w.txHist.Append(int(tx))
	rate = strings.ToLower(byteFormatInt(w.txHist.Last()))
	w.Lines[1].Title = fmt.Sprintf("TX [%s/s]", rate)
}
*/
