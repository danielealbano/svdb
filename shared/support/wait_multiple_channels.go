package shared_support

import "time"

func WaitMultipleChannels(waitInterval time.Duration, channels ...<-chan struct{}) {
	exit := false

	for !exit {
		for _, ch := range channels {
			select {
			case <-ch:
				exit = true
				break
			default:
				// do nothing
			}

			// sleep to avoid busy waiting
			time.Sleep(waitInterval)
		}
	}
}
