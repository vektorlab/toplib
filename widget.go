package toplib

import (
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib/toggle"
	"strings"
)

type ToggleMenu struct {
	ui.Block
	ToggleFgColor ui.Attribute
	ToggleBgColor ui.Attribute
	Toggles       toggle.Toggles
}

func NewToggleMenu(toggles toggle.Toggles) *ToggleMenu {
	return &ToggleMenu{
		Toggles:       toggles,
		ToggleFgColor: ui.ThemeAttr("list.item.fg"),
		ToggleBgColor: ui.ThemeAttr("list.item.bg"),
	}
}

func (tm *ToggleMenu) names() []string {
	names := []string{}
	for _, toggle := range tm.Toggles {
		names = append(names, toggle.Name)
	}
	return names
}

func (tm *ToggleMenu) Buffer() ui.Buffer {
	buf := tm.Block.Buffer()
	cs := ui.DefaultTxBuilder.Build(strings.Join(tm.names(), "\n"), tm.ToggleFgColor, tm.ToggleBgColor)
	i, j := 0, 0
	for i < len(cs) {
		if cs[j].Ch == '\n' {

		}
	}

	return buf
}

/*

type Gauges struct {
	Fields []string
}

func NewGauges(fields ...string) *Gauges {
	return &Gauges{Fields: fields}
}

func (g *Gauges) Buffers(r *Recorder, c *Cursor) []ui.GridBufferer {
	samples := fromSamples(r.Samples())
	sample := samples[c.IDX(samples)]
	gauges := make([]ui.GridBufferer, len(g.Fields))
	for i, field := range g.Fields {
		gauge := ui.NewGauge()
		gauge.Height = 3
		gauge.Percent = int(sample.GetFloat64(field))
		gauge.BorderLabel = field
		gauge.PercentColor = ui.ColorMagenta
		gauge.BarColor = ui.ColorRed
		gauge.BorderFg = ui.ColorWhite
		gauges[i] = gauge
	}
	return gauges
}

type Table struct {
	Fields []string
}

func NewTable(fields ...string) *Table {
	return &Table{Fields: fields}
}

func (t *Table) Buffers(r *Recorder, c *Cursor) []ui.GridBufferer {
	samples := r.Samples()
	rows := [][]string{t.Fields}
	for _, s := range samples {
		rows = append(rows, s.Strings(t.Fields))
	}
	table := ui.NewTable()
	table.Rows = rows
	table.Separator = false
	table.Border = false
	table.SetSize()
	table.Analysis()
	table.BgColors[c.IDX(samples)] = ui.ColorRed
	return []ui.GridBufferer{table}
}

type Summary struct {
	ui.Par
	count int
}

func NewSummary() *Summary {
	s := &Summary{}
	s.Height = 3
	s.Width = 50
	s.BorderLabel = "CTop Summary"
	s.TextFgColor = ui.ColorWhite
	s.BorderFg = ui.ColorCyan
	return s
}

func (s *Summary) Update(t *Top) {
	s.Text = fmt.Sprintf("Samples collected: %d", s.count)
}

func (s *Summary) Buffers(r *Recorder, _ *Cursor) []ui.GridBufferer {
	p := ui.NewPar(fmt.Sprintf("Samples collected: %d", r.Counter))
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "CTop Summary"
	p.BorderFg = ui.ColorCyan
	return []ui.GridBufferer{p}
}

type Chart struct {
	Field string
}

func NewChart(field string) *Chart {
	return &Chart{Field: field}
}

func (ch *Chart) Buffers(r *Recorder, c *Cursor) []ui.GridBufferer {
	samples := r.Samples()
	sample := samples[c.IDX(samples)]
	values := r.HistFloat64(sample.ID(), ch.Field)
	chart := ui.NewLineChart()
	chart.BorderLabel = ch.Field
	chart.Mode = "dot"
	chart.DotStyle = '+'
	chart.Height = 35
	chart.X = 0
	chart.Y = 100
	chart.Data = values
	chart.AxesColor = ui.ColorWhite
	chart.LineColor = ui.ColorYellow | ui.AttrBold
	return []ui.GridBufferer{chart}
}

type ContainerWidgets interface {
	Row() *ui.Row
	Render()
	Highlight()
	UnHighlight()
	SetCPU(int)
	SetNet(int64, int64)
	SetMem(int64, int64, int)
}

type Compact struct {
	Cid    *ui.Par
	Net    *ui.Par
	Name   *ui.Par
	Cpu    *ui.Gauge
	Memory *ui.Gauge
}

func NewCompact(id string, name string) *Compact {
	return &Compact{
		Cid:    compactPar(id),
		Net:    compactPar("-"),
		Name:   compactPar(name),
		Cpu:    slimGauge(),
		Memory: slimGauge(),
	}
}

func (w *Compact) Render() {
}

func (w *Compact) Row() *ui.Row {
	return ui.NewRow(
		ui.NewCol(2, 0, w.Name),
		ui.NewCol(2, 0, w.Cid),
		ui.NewCol(2, 0, w.Cpu),
		ui.NewCol(2, 0, w.Memory),
		ui.NewCol(2, 0, w.Net),
	)
}

func (w *Compact) Highlight() {
	w.Name.TextFgColor = ui.ColorDefault
	w.Name.TextBgColor = ui.ColorWhite
}

func (w *Compact) UnHighlight() {
	w.Name.TextFgColor = ui.ColorWhite
	w.Name.TextBgColor = ui.ColorDefault
}

func (w *Compact) SetCPU(val int) {
	w.Cpu.BarColor = colorScale(val)
	w.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		w.Cpu.BarColor = ui.ColorBlack
	}
	w.Cpu.Percent = val
}

func (w *Compact) SetNet(rx int64, tx int64) {
	w.Net.Text = fmt.Sprintf("%s / %s", byteFormat(rx), byteFormat(tx))
}

func (w *Compact) SetMem(val int64, limit int64, percent int) {
	w.Memory.Label = fmt.Sprintf("%s / %s", byteFormat(val), byteFormat(limit))
	if percent < 5 {
		percent = 5
		w.Memory.BarColor = ui.ColorBlack
	} else {
		w.Memory.BarColor = ui.ColorGreen
	}
	w.Memory.Percent = percent
}

func slimGauge() *ui.Gauge {
	g := ui.NewGauge()
	g.Height = 1
	g.Border = false
	g.Percent = 0
	g.PaddingBottom = 0
	g.BarColor = ui.ColorGreen
	g.Label = "-"
	return g
}

*/
