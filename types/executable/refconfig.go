package executable

func (e *ParallelRefConfig) RecordAttempt() {
	e.retryAttempts++
}

func (e *ParallelRefConfig) AttemptedMaxTimes() bool {
	return e.retryAttempts >= e.Retries+1
}

func (e *SerialRefConfig) RecordAttempt() {
	e.retryAttempts++
}

func (e *SerialRefConfig) AttemptedMaxTimes() bool {
	return e.retryAttempts >= e.Retries+1
}
