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
	maxLatencyMilli = float64(getEnvIntOrDefault("CHAOS_MAX_LATENCY_MS", 2000))
	latencySalt     = getEnvIntOrDefault("CHAOS_LATENCY_SALT", 10)
	errorSalt       = getEnvIntOrDefault("CHAOS_ERROR_SALT", 10)
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
		ticker := time.NewTicker(dur)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				stopCh <- struct{}{}
				return
			case <-ticker.C:
				fn()
			}
		}
	}(stopCh)

	return func() {
		stopCh <- struct{}{}
		<-stopCh
	}
}

func Latency() {
	factor := -float64(rand.Intn(latencySalt))

	latency := time.Duration(maxLatencyMilli*(1-math.Exp2(factor))) * time.Millisecond

	time.Sleep(latency)
}

func Error() error {
	var err error
	errFactor := rand.Intn(errorSalt)

	if errFactor == 0 {
		err = fmt.Errorf("some ramdom error")
	}

	return err
}
