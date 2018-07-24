package math

import (
	"math/rand"
	"time"
)

func Rand(min, max int) int {
	if min > max {
		return 0
	}

	tsp := time.Now().UnixNano()
	source := rand.NewSource(tsp)
	r := rand.New(source)

	randNum := r.Intn(max-min) + min
	return randNum
}
