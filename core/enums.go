package core

import (
	"fmt"
	"time"
)

type Sampling int

const (
	SamplingRaw Sampling = iota
	SamplingHour
	SamplingDay
)

var (
	samplingStrToValue      = map[string]Sampling{}
	samplingValueToStr      = map[Sampling]string{}
	samplingValueToDuration = map[Sampling]time.Duration{}
)

func init() {
	fullSamplingMapping := []struct {
		Value    Sampling
		String   string
		Duration time.Duration
	}{
		{SamplingRaw, "", 0},
		{SamplingHour, "hour", time.Hour},
		{SamplingDay, "day", 24 * time.Hour},
	}

	for _, m := range fullSamplingMapping {
		samplingStrToValue[m.String] = m.Value
		samplingValueToStr[m.Value] = m.String
		samplingValueToDuration[m.Value] = m.Duration
	}
}

func ParseSampling(value string) (Sampling, error) {
	var err error

	sampling, ok := samplingStrToValue[value]
	if !ok {
		err = fmt.Errorf("Could not prase %v as sampling", value)
	}
	return sampling, err
}

func (s Sampling) Duration() time.Duration {
	// Ignore error as it has to exit
	duration, _ := samplingValueToDuration[s]
	return duration
}

func (s Sampling) String() string {
	// Ignore error as it has to exit
	str, _ := samplingValueToStr[s]
	return str
}
