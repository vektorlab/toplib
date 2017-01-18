// +build ignore

package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
	"github.com/kevinschoon/ctop"
)

type Container struct {
	id     string
	name   string
	done   chan bool
	stats  chan *docker.Stats
	reader *StatReader
}

func NewContainer(c docker.APIContainers) *Container {
	id := c.ID[:12]
	name := strings.Replace(c.Names[0], "/", "", 1) // use primary container name
	return &Container{
		id:     id,
		name:   name,
		done:   make(chan bool),
		stats:  make(chan *docker.Stats),
		reader: &StatReader{},
	}
}

func (c *Container) Collect(client *docker.Client) {
	go func() {
		opts := docker.StatsOptions{
			ID:     c.id,
			Stats:  c.stats,
			Stream: true,
			Done:   c.done,
		}
		client.Stats(opts)
	}()
}

var filters = map[string][]string{
	"status": []string{"running"},
}

func NewContainerMap(host string) *ContainerMap {
	// init docker client
	client, err := docker.NewClient(host)
	if err != nil {
		panic(err)
	}
	cm := &ContainerMap{
		client:     client,
		containers: make(map[string]*Container),
	}
	return cm
}

type ContainerMap struct {
	client     *docker.Client
	containers map[string]*Container
}

func (cm *ContainerMap) Samples() []*ctop.Sample {
	samples := []*ctop.Sample{}
	opts := docker.ListContainersOptions{
		Filters: filters,
	}
	containers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		id := c.ID[:12]
		if _, ok := cm.containers[id]; !ok {
			cm.containers[id] = NewContainer(c)
			cm.containers[id].Collect(cm.client)
		}
	}
loop:
	for id, container := range cm.containers {
		select {
		case stats := <-container.stats:
			if stats == nil {
				delete(cm.containers, id[:12])
				continue loop
			}
			container.reader.Read(stats)
			cm.containers[id].reader.Read(stats)
			samples = append(samples, container.reader.Sample(id[:12], container.name))
		case <-container.done:
			delete(cm.containers, id[:12])
		}
	}
	return samples
}

type StatReader struct {
	CPUUtil    int
	NetTx      int64
	NetRx      int64
	MemLimit   int64
	MemPercent int
	MemUsage   int64
	lastCpu    float64
	lastSysCpu float64
}

// TODO: Remove redundant type casting
func (s StatReader) Sample(id, name string) *ctop.Sample {
	sample := ctop.NewSample(id)
	sample.SetString("NAME", name)
	sample.SetFloat64("CPUUtil", float64(s.CPUUtil))
	sample.SetFloat64("NetTx", float64(s.NetTx))
	sample.SetFloat64("NetRx", float64(s.NetRx))
	sample.SetFloat64("MemLimit", float64(s.MemLimit))
	sample.SetFloat64("MemPercent", float64(s.MemPercent))
	sample.SetFloat64("MemUsage", float64(s.MemUsage))
	return sample
}

func (s *StatReader) Read(stats *docker.Stats) {
	s.ReadCPU(stats)
	s.ReadMem(stats)
	s.ReadNet(stats)
}

func (s *StatReader) ReadCPU(stats *docker.Stats) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - s.lastCpu
	syscpudiff := system - s.lastSysCpu
	s.CPUUtil = round((cpudiff / syscpudiff * 100) * ncpus)
	s.lastCpu = total
	s.lastSysCpu = system
}

func (s *StatReader) ReadMem(stats *docker.Stats) {
	s.MemUsage = int64(stats.MemoryStats.Usage)
	s.MemLimit = int64(stats.MemoryStats.Limit)
	s.MemPercent = round((float64(s.MemUsage) / float64(s.MemLimit)) * 100)
}

func (s *StatReader) ReadNet(stats *docker.Stats) {
	s.NetTx, s.NetRx = 0, 0
	for _, network := range stats.Networks {
		s.NetTx += int64(network.TxBytes)
		s.NetRx += int64(network.RxBytes)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func main() {
	top := ctop.NewTop()
	var (
		cursor = ctop.NewCursor()
		table  = ctop.NewTable(
			"ID", "NAME", "CPUUtil", "MemLimit",
			"MemPercent", "MemUsage", "NetTx", "NetRx",
		)
		sortMenu = ctop.NewMenu(
			"ID", "NAME", "CPUUtil", "MemLimit",
			"MemPercent", "MemUsage", "NetTx", "NetRx",
		)
		toggles = ctop.NewToggles(&ctop.Toggle{Name: "sort"})
	)
	defaultView := ctop.NewView(func() []*ui.Row {
		return []*ui.Row{
			ctop.NewHeader().Row(),
			ui.NewRow(
				ui.NewCol(12, 0, table.Buffers(top.Recorder, cursor)...),
			),
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

	sortView := ctop.NewView(func() []*ui.Row {
		return []*ui.Row{
			ctop.NewHeader().Row(),
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
	top.Views.Set("default")

	cm := NewContainerMap("unix:///var/run/docker.sock")
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
	loop:
		for {
			select {
			case <-top.Exit:
				for _, container := range cm.containers {
					container.done <- true
				}
				ui.StopLoop()
				break loop
			case <-ticker.C:
				top.Samples <- cm.Samples()
			}
		}
		ticker.Stop()
	}()
	if err := ctop.Run(top); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
