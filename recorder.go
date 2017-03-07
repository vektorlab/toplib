package toplib

import (
	"github.com/vektorlab/toplib/sample"
	"sync"
)

const MaxSamples = 350

// Recorder saves recorded samples
type Recorder struct {
	SortField string
	samples   map[string][]*sample.Sample
	Counter   int
	latest    []*sample.Sample
	mu        sync.RWMutex
}

func NewRecorder() *Recorder {
	return &Recorder{
		SortField: "ID",
		samples:   map[string][]*sample.Sample{},
		latest:    []*sample.Sample{},
	}
}

func (r *Recorder) store(s *sample.Sample) {
	if samples, ok := r.samples[s.ID()]; !ok {
		r.samples[s.ID()] = []*sample.Sample{}
	} else {
		if len(samples) >= MaxSamples {
			//Pop
			r.samples[s.ID()] = samples[1:]
		}
	}
	r.samples[s.ID()] = append(r.samples[s.ID()], s)
}

func (r *Recorder) HistFloat64(id, field string) []float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	values := []float64{}
	if samples, ok := r.samples[id]; ok {
		for _, s := range samples {
			values = append(values, s.GetFloat64(field))
		}
	}
	return values
}

func (r *Recorder) HistString(id, field string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	values := []string{}
	if samples, ok := r.samples[id]; ok {
		for _, s := range samples {
			values = append(values, s.GetString(field))
		}
	}
	return values
}

func (r *Recorder) Items() []Item {
	items := []Item{}
	for _, sample := range r.Samples() {
		items = append(items, sample)
	}
	return items
}

func (r *Recorder) Samples() []*sample.Sample {
	r.mu.Lock()
	defer r.mu.Unlock()
	sample.Sort(r.SortField, r.latest)
	return r.latest
}

func (r *Recorder) Load(samples []*sample.Sample) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latest = samples
	for _, s := range r.latest {
		r.store(s)
		r.Counter++
	}
}
