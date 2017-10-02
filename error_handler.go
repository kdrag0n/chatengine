package main

func handlePanic() {
	if err := recover(); err != nil {
		logger.Error(err)
	}
}
