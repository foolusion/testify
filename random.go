package main

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
)

const longScale = float64(0xFFFFFFFFFFFFFFF)

type hashedUnit int64

func newHashedUnit(localSalt string, units []string) (hashedUnit, error) {
	key := fmt.Sprintf("%s.%s%s%s", localSalt, config.globalSalt, config.saltSeperator, strings.Join(units, "."))
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

func (h hashedUnit) bernoulliTrail(p float64) (bool, error) {
	if p < 0 || p > 1 {
		return false, fmt.Errorf("p must be between 0 and 1: %v", p)
	}
	return h.getUniform(0, 1) > p, nil
}
