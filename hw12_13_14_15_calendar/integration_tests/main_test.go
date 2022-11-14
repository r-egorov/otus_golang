package main

import (
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const wait = 13 * time.Second

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                "integration-tests",
		ScenarioInitializer: InitializeScenario,
		Options:             nil,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
