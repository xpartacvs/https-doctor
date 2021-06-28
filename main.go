package main

import (
	"https-doctor/packages/logger"
	"https-doctor/packages/worker"
)

func main() {
	logger.Log().Info().Msg("Application stared.")

	if err := worker.Start(); err != nil {
		logger.Log().Err(err).Msg("Application has been terminated due to unexpected error.")
	}
}
