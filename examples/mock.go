// +build ignore

package main

import (
	"fmt"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/sample"
	"github.com/vektorlab/toplib/section"
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
	samples []*sample.Sample
}

func (m *MockSource) Collect() []*sample.Sample {
	if m.samples == nil {
		m.samples = []*sample.Sample{}
		for i := 0; i < 35; i++ {
			m.samples = append(m.samples, sample.NewSample(RandString(10)))
		}
	}
	samples := []*sample.Sample{}
	for _, smpl := range m.samples {
		/*
			// Uncomment to change iteration size
			if rand.Intn(10) > 5 {
				continue
			}
		*/
		s := sample.NewSample(smpl.ID())
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
	sections := []toplib.Section{
		section.NewSamples("ID", "CPU", "MEM", "DISK", "GPU", "THING", "OTHER THING"),
		section.NewDebug(),
	}

	top := toplib.NewTop(sections)
	tick := time.NewTicker(500 * time.Millisecond)

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

	if err := toplib.Run(top); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
