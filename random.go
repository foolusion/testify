package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const longScale = float64(0xFFFFFFFFFFFFFFF)

type hashedUnit int64

func mapJoin(vals url.Values, sep, eq string) string {
	list := make([]string, 0, len(vals))
	for k := range vals {
		list = append(list, k)
	}
	sort.Strings(list)

	var buf bytes.Buffer
	for i, k := range list {
		buf.WriteString(k)
		buf.WriteString(eq)
		buf.WriteString(strings.Join(vals[k], ","))
		if i < len(list)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

func newHashedUnit(localSalt string, units url.Values) (hashedUnit, error) {
	key := fmt.Sprintf("%s.%s%s%s", localSalt, config.globalSalt, config.saltSeperator, mapJoin(units, ":", "="))
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	i, err := strconv.ParseInt(hash[:15], 16, 64)
	return hashedUnit(i), err
}

func (h hashedUnit) getUniform(min, max float64) float64 {
	return min + (max-min)*(float64(h)/longScale)
}

func (h hashedUnit) randomFloat(min, max float64) float64 {
	return h.getUniform(min, max)
}

func (h hashedUnit) randomInt(min, max int64) int64 {
	return min + int64(h)%(max-min+1)
}

func (h hashedUnit) bernoulliTrial(p float64) (bool, error) {
	if p < 0 || p > 1 {
		return false, fmt.Errorf("p must be between 0 and 1: %v", p)
	}
	return h.getUniform(0, 1) > p, nil
}

func (h hashedUnit) uniformChoice(choices []string) string {
	if len(choices) == 0 {
		return ""
	}
	return choices[int(h)%len(choices)]
}

func (h hashedUnit) weightedChoice(choices []string, weights []float64) (string, error) {
	if len(choices) != len(weights) {
		return "", fmt.Errorf("len(choices) != len(weights): %v != %v", len(choices), len(weights))
	}
	selection := make([]float64, len(weights))
	cumSum := 0.0
	for i, v := range weights {
		cumSum += v
		selection[i] = cumSum
	}
	choice := h.getUniform(0, cumSum)
	for i, v := range selection {
		if v > choice {
			return choices[i], nil
		}
	}
	return "", fmt.Errorf("weightedChoice: wtf: shouldn't be here")
}
