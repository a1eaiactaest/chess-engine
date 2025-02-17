package main

import (
	"chess-engine/controller"
	"chess-engine/engine"
	"chess-engine/test"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		controller.Main()
	}()

	test.TestGochess()
	engine.FeedbackEngine()

	wg.Wait()
}
