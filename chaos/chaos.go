package chaos

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	maxLatencyMilli = float64(getEnvIntOrDefault("CHAOS_MAX_LATENCY_MILLI", 2000))
	salt            = getEnvIntOrDefault("CHAOS_SALT", 10)
)

func getEnvIntOrDefault(envName string, defaultValue int) int {
	valueString := os.Getenv(envName)

	if valueString == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueString)

	if err != nil {
		log.Printf("couldn't parse env %s: %s", envName, err)
		return defaultValue
	}

	return value
}

func Do(fn func(), dur time.Duration) func() {
	var stopCh chan struct{}
	go func(stopCh chan struct{}) {
		for {
			select {
			case <-stopCh:
				return
			case <-time.Tick(dur):
				fn()
			}
		}
	}(stopCh)

	return func() {
		stopCh <- struct{}{}
	}
}

func Latency() {
	factor := -float64(rand.Intn(salt))

	latency := time.Duration(maxLatencyMilli*(1-math.Exp2(factor))) * time.Millisecond

	time.Sleep(latency)
}

func Error() error {
	var err error
	errFactor := rand.Intn(salt)

	if errFactor == 0 {
		err = fmt.Errorf("some ramdom error")
	}

	return err
}
