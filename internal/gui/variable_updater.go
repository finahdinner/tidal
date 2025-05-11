package gui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/finahdinner/tidal/internal/twitch"
)

const updateVariablesTimeout = 5 * time.Second

var updaterTicker *time.Ticker
var updaterTickerDone chan struct{}
var updateVariablesSectionSignal = make(chan struct{}, 1)

func startUpdatingVariables(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive integer")
	}
	if updaterTicker != nil {
		return errors.New("ticker already running - stop it first")
	}

	updaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
	updaterTickerDone = make(chan struct{})

	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		ctx := context.Background()

		// initial update, before the ticker
		if err := callUpdateStreamVariablesWithTimeout(ctx); err != nil {
			log.Printf("failed - err: %v", err)
			if errors.Is(err, twitch.Err401Unauthorised) {
				errChan <- err
				return
			}
		}

		for {
			select {
			case <-updaterTickerDone:
				log.Println("updaterTickerDone closed")
				doneChan <- struct{}{}
				return
			case <-updaterTicker.C:
				if err := callUpdateStreamVariablesWithTimeout(ctx); err != nil {
					log.Printf("failed - err: %v", err)
					if errors.Is(err, twitch.Err401Unauthorised) {
						errChan <- err
						return
					}
					continue
				}
				select {
				case updateVariablesSectionSignal <- struct{}{}:
					// signal to update widgets in variables section
				default:
					// reached if updateSignalChan is full
				}
			}
		}
	}()

	select {
	case err := <-errChan:
		stopUpdaterTicker()
		return fmt.Errorf("unauthorised 401 - cancelling - err: %w", err)
	case <-doneChan:
		return nil
	}

}

func stopUpdaterTicker() {
	if updaterTicker != nil {
		log.Println("ticker 'updaterTicker' stopped")
		updaterTicker.Stop()
		updaterTicker = nil
	}
	if updaterTickerDone != nil {
		log.Println("chan 'updaterTickerDone' closed")
		close(updaterTickerDone)
		updaterTickerDone = nil
	}
}

// attempts to update the stream variables, but cancels if the timeout limit is exceeded
func callUpdateStreamVariablesWithTimeout(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, updateVariablesTimeout)
	defer cancel()
	return twitch.UpdateStreamVariables(ctx)
}
