package toplib

import (
	"github.com/vektorlab/toplib/sample"
	"sync"
)

const MaxSamples = 350

// Recorder saves recorded samples
type Recorder struct {
	SortField string
	samples   map[sample.Namespace]map[string][]*sample.Sample
	latest    map[sample.Namespace][]*sample.Sample
	Counter   int
	mu        sync.RWMutex
}

func NewRecorder() *Recorder {
	return &Recorder{
		SortField: "ID",
		samples:   map[sample.Namespace]map[string][]*sample.Sample{},
	}
}

func (r *Recorder) load(s *sample.Sample) {
	id := s.ID()
	namespace := s.Namespace()
	if _, ok := r.samples[namespace]; !ok {
		r.samples[namespace] = map[string][]*sample.Sample{}
	}
	if samples, ok := r.samples[namespace][id]; !ok {
		r.samples[namespace][id] = []*sample.Sample{}
	} else {
		if len(samples) >= MaxSamples {
			//Pop
			r.samples[namespace][id] = samples[1:]
		}
	}
	r.samples[namespace][id] = append(r.samples[namespace][id], s)
}

func (r *Recorder) Latest(namespace sample.Namespace) []*sample.Sample {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.latest[namespace]; !ok {
		return nil
	}
	return r.latest[namespace]
}

func (r *Recorder) Load(samples []*sample.Sample) {
	r.mu.Lock()
	defer r.mu.Unlock()
	latest := map[sample.Namespace][]*sample.Sample{}
	for _, smpl := range samples {
		namespace := smpl.Namespace()
		if _, ok := latest[namespace]; !ok {
			latest[namespace] = []*sample.Sample{}
		}
		latest[namespace] = append(latest[namespace], smpl)
		r.load(smpl)
		r.Counter++
	}
	r.latest = latest
}
