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
	samples map[sample.Namespace][]*sample.Sample
}

func (m *MockSource) Collect() ([]*sample.Sample, error) {
	results := []*sample.Sample{}
	for _, samples := range m.samples {
		for _, smpl := range samples {
			// Uncomment to change iteration size
			if rand.Intn(10) > 5 {
				continue
			}

			s := sample.NewSample(smpl.ID(), smpl.Namespace())
			s.SetFloat64("CPU", float64(rand.Intn(5)))
			s.SetFloat64("MEM", float64(rand.Intn(40)))
			s.SetFloat64("DISK", float64(rand.Intn(60)))
			s.SetFloat64("GPU", float64(rand.Intn(90)))
			s.SetString("THING", RandString(50))
			s.SetString("OTHER THING", RandString(20))
			results = append(results, s)
		}
	}
	return results, nil
}

func main() {
	namespaces := []sample.Namespace{
		sample.Namespace("ns-0"),
		sample.Namespace("ns-1"),
		sample.Namespace("ns-2"),
	}
	source := &MockSource{
		samples: map[sample.Namespace][]*sample.Sample{
			namespaces[0]: []*sample.Sample{},
			namespaces[1]: []*sample.Sample{},
			namespaces[2]: []*sample.Sample{},
		},
	}
	for i := 0; i < 35; i++ {
		source.samples[namespaces[0]] = append(source.samples[namespaces[0]], sample.NewSample(RandString(10), namespaces[0]))
		source.samples[namespaces[1]] = append(source.samples[namespaces[1]], sample.NewSample(RandString(10), namespaces[1]))
		source.samples[namespaces[2]] = append(source.samples[namespaces[2]], sample.NewSample(RandString(10), namespaces[2]))
	}
	sections := []toplib.Section{
		section.NewSamples(namespaces[0], "ID", "CPU", "MEM", "DISK", "GPU"),
		section.NewSamples(namespaces[1], "ID", "CPU", "MEM", "DISK", "GPU", "THING"),
		section.NewSamples(namespaces[2], "ID", "CPU", "MEM", "DISK", "GPU", "THING", "OTHER THING"),
		&section.Debug{Namespaces: namespaces},
	}
	if err := toplib.Run(toplib.NewTop(sections), source.Collect); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
