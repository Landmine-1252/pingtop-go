package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/landmine-1252/pingtop-go/internal/pingtop"
)

func TestBuildExitSummaryIncludesTitleAndRuntime(t *testing.T) {
	renderer := &Renderer{}
	startedAt := time.Unix(100, 0)
	completedAt := time.Unix(225, 0)
	output := renderer.BuildExitSummary(pingtop.StateSnapshot{
		Diagnosis:            "All monitored targets are reachable",
		LastCycleCompletedAt: completedAt,
		Session: pingtop.SessionTotals{
			StartedAt:       startedAt,
			CyclesCompleted: 2,
			TotalChecks:     10,
			Successes:       10,
		},
	})

	if !strings.Contains(output, "pingtop  Session summary") {
		t.Fatalf("expected title line, got %q", output)
	}
	if !strings.Contains(output, "ran 2m5s") {
		t.Fatalf("expected runtime in output, got %q", output)
	}
}

func TestBuildScreenIncludesVersionInHeader(t *testing.T) {
	renderer := &Renderer{}
	output := renderer.BuildScreen(
		pingtop.StateSnapshot{
			Diagnosis:          "All monitored targets are reachable",
			StatsWindowSeconds: 3600,
		},
		pingtop.AppConfig{
			CheckIntervalSeconds:     1,
			PingTimeoutMS:            1200,
			UIRefreshIntervalSeconds: 0.5,
			LatencyWarningMS:         100,
			LatencyCriticalMS:        250,
			DiagnosisConfirmCycles:   2,
			RecoveryConfirmCycles:    2,
			LoggingMode:              "around_failure",
			VisibleEventLines:        8,
		},
		false,
		false,
		nil,
		UpdateStatus{},
	)

	if !strings.Contains(output, "pingtop "+pingtop.Version) {
		t.Fatalf("expected versioned header, got %q", output)
	}
}

func TestBuildScreenOmitsPromptHelpWhenNoPromptIsActive(t *testing.T) {
	renderer := &Renderer{}
	output := renderer.BuildScreen(
		pingtop.StateSnapshot{
			Diagnosis:          "All monitored targets are reachable",
			StatsWindowSeconds: 3600,
		},
		pingtop.AppConfig{
			CheckIntervalSeconds:     1,
			PingTimeoutMS:            1200,
			UIRefreshIntervalSeconds: 0.5,
			HelpVisible:              true,
			LatencyWarningMS:         100,
			LatencyCriticalMS:        250,
			DiagnosisConfirmCycles:   2,
			RecoveryConfirmCycles:    2,
			LoggingMode:              "around_failure",
			VisibleEventLines:        8,
		},
		false,
		true,
		nil,
		UpdateStatus{},
	)

	if strings.Contains(output, "Enter submit") {
		t.Fatalf("expected prompt help to stay hidden without an active prompt, got %q", output)
	}
}

func TestBuildScreenShowsPromptHelpInsideActivePrompt(t *testing.T) {
	renderer := &Renderer{}
	output := renderer.BuildScreen(
		pingtop.StateSnapshot{
			Diagnosis:          "All monitored targets are reachable",
			StatsWindowSeconds: 3600,
		},
		pingtop.AppConfig{
			CheckIntervalSeconds:     1,
			PingTimeoutMS:            1200,
			UIRefreshIntervalSeconds: 0.5,
			HelpVisible:              true,
			LatencyWarningMS:         100,
			LatencyCriticalMS:        250,
			DiagnosisConfirmCycles:   2,
			RecoveryConfirmCycles:    2,
			LoggingMode:              "around_failure",
			VisibleEventLines:        8,
		},
		false,
		true,
		&PromptState{Kind: "add", Message: "enter a hostname or IP to add"},
		UpdateStatus{},
	)

	if !strings.Contains(output, "Enter submit | Esc cancel | Backspace edit") {
		t.Fatalf("expected active prompt to include prompt help, got %q", output)
	}
}
