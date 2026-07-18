package main

type simpleCounter struct {
	total    float64
	duration float64
}

func (lc *simpleCounter) Parse(b []byte) error {
	lc.total = lc.total + 1
	return nil
}

func (lc *simpleCounter) Finish(duration float64) {
	lc.duration = duration
}

func (lc *simpleCounter) GetTotal() float64 {
	return lc.total
}

func (lc *simpleCounter) GetDuration() float64 {
	return lc.duration
}
