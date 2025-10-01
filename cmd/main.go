package main

import (
	"github.com/MaximGubanov/googger"
)

func main() {
	log, err := googger.NewLogger("test.log", "test", "D")
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Debug("Тест DEBUG")
	log.Info("Тест INFO")
	log.Warn("Тест WARN")
	log.Error("Тест ERROR")
}
