package toplib

import "sync"

const MaxSamples = 350

type Recorder struct {
	SortField string
	samples   map[string][]*Sample
	Counter   int
	latest    []*Sample
	mu        sync.RWMutex
}

func NewRecorder() *Recorder {
	return &Recorder{
		SortField: "ID",
		samples:   map[string][]*Sample{},
		latest:    []*Sample{},
	}
}

func (r *Recorder) store(sample *Sample) {
	if samples, ok := r.samples[sample.ID()]; !ok {
		r.samples[sample.ID()] = []*Sample{}
	} else {
		if len(samples) >= MaxSamples {
			//Pop
			r.samples[sample.ID()] = samples[1:]
		}
	}
	r.samples[sample.ID()] = append(r.samples[sample.ID()], sample)
}

func (r *Recorder) HistFloat64(id, field string) []float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	values := []float64{}
	if samples, ok := r.samples[id]; ok {
		for _, sample := range samples {
			values = append(values, sample.GetFloat64(field))
		}
	}
	return values
}

func (r *Recorder) HistString(id, field string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	values := []string{}
	if samples, ok := r.samples[id]; ok {
		for _, sample := range samples {
			values = append(values, sample.GetString(field))
		}
	}
	return values
}

func (r *Recorder) Samples() []*Sample {
	r.mu.Lock()
	defer r.mu.Unlock()
	Sort(r.SortField, r.latest)
	return r.latest
}

func (r *Recorder) Load(samples []*Sample) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latest = samples
	for _, sample := range r.latest {
		r.store(sample)
		r.Counter++
	}
}
