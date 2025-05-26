package gui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/twitch"
)

const updateVariablesTimeout = 3 * time.Second

var twitchVariableUpdaterTicker *time.Ticker
var twitchVariableUpdaterTickerDone chan struct{}
var updateTwitchVariablesSectionSignal = make(chan struct{}, 1)

func startUpdatingTwitchVariables(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive integer")
	}
	if twitchVariableUpdaterTicker != nil {
		return errors.New("ticker already running - stop it first")
	}

	twitchVariableUpdaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
	twitchVariableUpdaterTickerDone = make(chan struct{})

	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		ctx := context.Background()

		// initial update, before the ticker
		if err := callUpdateTwitchVariablesWithTimeout(ctx); err != nil {
			config.Logger.LogInfof("failed - err: %v", err)
			if errors.Is(err, twitch.Err401Unauthorised) {
				errChan <- err
				return
			}
		}

		for {
			select {
			case <-twitchVariableUpdaterTickerDone:
				config.Logger.LogInfo("updaterTickerDone closed")
				doneChan <- struct{}{}
				return
			case <-twitchVariableUpdaterTicker.C:
				if err := callUpdateTwitchVariablesWithTimeout(ctx); err != nil {
					config.Logger.LogInfof("failed - err: %v", err)
					if errors.Is(err, twitch.Err401Unauthorised) {
						errChan <- err
						return
					}
					continue
				}
				select {
				case updateTwitchVariablesSectionSignal <- struct{}{}:
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
	if twitchVariableUpdaterTicker != nil {
		config.Logger.LogInfo("ticker 'updaterTicker' stopped")
		twitchVariableUpdaterTicker.Stop()
		twitchVariableUpdaterTicker = nil
	}
	if twitchVariableUpdaterTickerDone != nil {
		config.Logger.LogInfo("chan 'updaterTickerDone' closed")
		close(twitchVariableUpdaterTickerDone)
		twitchVariableUpdaterTickerDone = nil
	}
}

// Attempts to update the twitch variables, but cancels if the timeout limit is exceeded
func callUpdateTwitchVariablesWithTimeout(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, updateVariablesTimeout)
	defer cancel()
	return twitch.UpdateTwitchVariables(ctx)
}
