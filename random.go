package main

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
)

const longScale = float64(0xFFFFFFFFFFFFFFF)

func getHash(localSalt string, units []string) (int64, error) {
	key := fmt.Sprintf("%s.%s%s%s", localSalt, config.globalSalt, config.saltSeperator, strings.Join(units, "."))
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	return strconv.ParseInt(hash[:15], 16, 64)
}

func getUniform(localSalt string, units []string, min, max float64) (float64, error) {
	hash, err := getHash(localSalt, units)
	if err != nil {
		return 0, err
	}
	return min + (max-min)*(float64(hash)/longScale), nil
}

func randomFloat(localSalt string, units []string, min, max float64) (float64, error) {
	return getUniform(localSalt, units, min, max)
}

func randomInt(localSalt string, units []string, min, max int64) (int64, error) {
	hash, err := getHash(localSalt, units)
	if err != nil {
		return 0, err
	}
	return min + hash%(max-min+1), nil
}

func bernoulliTrail(p float64) (int64, error)
