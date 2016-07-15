package main

import (
	"github.com/rubyist/circuitbreaker"
	"time"
	log "github.com/Sirupsen/logrus"
)

type ExtendedCircuitBreakerMeta struct {
	CircuitBreakerMeta
	CB *circuit.Breaker
}

func NewCircuitBreaker(apiSpec *APISpec) *ExtendedCircuitBreakerMeta {
	breakerMeta := &ExtendedCircuitBreakerMeta{CircuitBreakerMeta: apiSpec.CircuitBreaker}
	breakerMeta.CB = circuit.NewRateBreaker(apiSpec.CircuitBreaker.ThresholdPercent, apiSpec.CircuitBreaker.Samples)

	events := breakerMeta.CB.Subscribe()

	go func() {
		path := apiSpec.Proxy.ListenPath
		timerActive := false
		for {
			e := <-events
			switch e {
			case circuit.BreakerTripped:
				log.Warning("[PROXY] [CIRCUIT BREKER] Breaker tripped for path: ", path)
				log.Debug("Breaker tripped: ", e)

				// Start a timer function
				if !timerActive {
					go func(timeout int, breaker *circuit.Breaker) {
						log.Debug("-- Sleeping for (s): ", timeout)
						time.Sleep(time.Duration(timeout) * time.Second)
						log.Debug("-- Resetting breaker")
						breaker.Reset()
						timerActive = false
					}(apiSpec.CircuitBreaker.ReturnToServiceAfter, breakerMeta.CB)
					timerActive = true
				}
			}
		}
	}()

	return breakerMeta
}
