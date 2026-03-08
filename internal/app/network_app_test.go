package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmine-1252/pingtop-go/internal/checks"
	"github.com/landmine-1252/pingtop-go/internal/pingtop"
)

func TestRunHeadlessOncePrintsSummary(t *testing.T) {
	tempDir := t.TempDir()
	configManager := pingtop.NewConfigManager(filepath.Join(tempDir, "pingtop.json"))
	config := configManager.Update(func(config *AppConfig) {
		config.Targets = nil
	})
	stateStore := pingtop.NewStateStore(config)
	logger := pingtop.NewCSVLogger(filepath.Join(tempDir, "pingtop_log.csv"))
	coordinator := checks.NewCheckCoordinator(checks.NewPingRunner(), nil)
	defer coordinator.Close()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer
	rc := runHeadless(configManager, stateStore, logger, coordinator, true)
	_ = writer.Close()
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read captured stdout: %v", err)
	}
	if rc != 0 {
		t.Fatalf("expected rc 0, got %d", rc)
	}
	if !strings.Contains(string(output), "Session summary:") {
		t.Fatalf("expected summary output, got %q", string(output))
	}
}

func TestRunLongHelpPrintsUsage(t *testing.T) {
	output := captureStdout(t, func() int {
		return Run([]string{"--help"})
	})
	if !strings.Contains(output, "Usage: pingtop [flags]") {
		t.Fatalf("expected usage output, got %q", output)
	}
}

func TestRunShortHelpPrintsUsage(t *testing.T) {
	output := captureStdout(t, func() int {
		return Run([]string{"-h"})
	})
	if !strings.Contains(output, "Usage: pingtop [flags]") {
		t.Fatalf("expected usage output, got %q", output)
	}
}

func TestRunShortVersionPrintsVersion(t *testing.T) {
	output := captureStdout(t, func() int {
		return Run([]string{"-v"})
	})
	if !strings.Contains(output, pingtop.Version) {
		t.Fatalf("expected version output, got %q", output)
	}
}

func TestParseArgsSupportsShortAliases(t *testing.T) {
	args, err := parseArgs([]string{"-n", "-o", "-v"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !args.noUI {
		t.Fatal("expected -n to enable no-ui mode")
	}
	if !args.once {
		t.Fatal("expected -o to enable once mode")
	}
	if !args.showVersion {
		t.Fatal("expected -v to enable version output")
	}
}

func TestParseArgsCapturesPositionalTargets(t *testing.T) {
	args, err := parseArgs([]string{"-n", "example.com", "1.1.1.1"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !args.noUI {
		t.Fatal("expected -n to enable no-ui mode")
	}
	if len(args.targets) != 2 || args.targets[0] != "example.com" || args.targets[1] != "1.1.1.1" {
		t.Fatalf("unexpected targets: %#v", args.targets)
	}
}

func TestBuildServicesWithPositionalTargetsDisablesLogging(t *testing.T) {
	tempDir := t.TempDir()
	runtimePaths := pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	}
	services, err := buildServices(runtimePaths, cliArgs{targets: []string{"example.com", "1.1.1.1"}})
	if err != nil {
		t.Fatalf("unexpected buildServices error: %v", err)
	}
	defer services.coordinator.Close()

	config := services.configManager.Snapshot()
	if len(config.Targets) != 2 || config.Targets[0].Value != "example.com" || config.Targets[1].Value != "1.1.1.1" {
		t.Fatalf("unexpected configured targets: %#v", config.Targets)
	}

	success := true
	latency := 12.0
	services.logger.LogResults([]CheckResult{{
		CycleID:       1,
		Timestamp:     pingtop.NewSessionTotals().StartedAt,
		Target:        "example.com",
		TargetType:    "hostname",
		ResolvedIP:    "93.184.216.34",
		DNSSuccess:    &success,
		PingSuccess:   true,
		LatencyMS:     &latency,
		ErrorCategory: "ok",
	}}, config)

	if _, err := os.Stat(runtimePaths.LogPath); !os.IsNotExist(err) {
		t.Fatalf("expected positional-target run to avoid creating log file, got err=%v", err)
	}
}

func TestHandleKeyEscQuitsWhenNoPromptIsActive(t *testing.T) {
	tempDir := t.TempDir()
	services, err := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	}, cliArgs{})
	if err != nil {
		t.Fatalf("unexpected buildServices error: %v", err)
	}
	defer services.coordinator.Close()

	ui := NewPingTopUI(
		services.runtimePaths,
		services.configManager,
		services.stateStore,
		services.logger,
		services.coordinator,
		services.updateManager,
	)
	ui.handleKey("\x1b")

	if ui.running {
		t.Fatal("expected Esc to stop the UI when no prompt is active")
	}
}

func TestHandleKeyEscCancelsPromptWithoutQuitting(t *testing.T) {
	tempDir := t.TempDir()
	services, err := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	}, cliArgs{})
	if err != nil {
		t.Fatalf("unexpected buildServices error: %v", err)
	}
	defer services.coordinator.Close()

	ui := NewPingTopUI(
		services.runtimePaths,
		services.configManager,
		services.stateStore,
		services.logger,
		services.coordinator,
		services.updateManager,
	)
	ui.prompt = &PromptState{Kind: "add", Message: "test"}
	ui.handleKey("\x1b")

	if !ui.running {
		t.Fatal("expected Esc to keep the UI running while canceling a prompt")
	}
	if ui.prompt != nil {
		t.Fatal("expected Esc to clear the active prompt")
	}
}

func TestNewPingTopUIDefaultsHelpVisibleFromConfig(t *testing.T) {
	tempDir := t.TempDir()
	services, err := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	}, cliArgs{})
	if err != nil {
		t.Fatalf("unexpected buildServices error: %v", err)
	}
	defer services.coordinator.Close()

	ui := NewPingTopUI(
		services.runtimePaths,
		services.configManager,
		services.stateStore,
		services.logger,
		services.coordinator,
		services.updateManager,
	)

	if !ui.helpVisible {
		t.Fatal("expected help to be visible by default")
	}
}

func TestHandleKeyHTogglesAndPersistsHelpVisibility(t *testing.T) {
	tempDir := t.TempDir()
	services, err := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	}, cliArgs{})
	if err != nil {
		t.Fatalf("unexpected buildServices error: %v", err)
	}
	defer services.coordinator.Close()

	ui := NewPingTopUI(
		services.runtimePaths,
		services.configManager,
		services.stateStore,
		services.logger,
		services.coordinator,
		services.updateManager,
	)

	ui.handleKey("h")
	if ui.helpVisible {
		t.Fatal("expected h to hide help")
	}
	if services.configManager.Snapshot().HelpVisible {
		t.Fatal("expected help visibility change to persist to config")
	}

	reopened := NewPingTopUI(
		services.runtimePaths,
		services.configManager,
		services.stateStore,
		services.logger,
		services.coordinator,
		services.updateManager,
	)
	if reopened.helpVisible {
		t.Fatal("expected reopened UI to use saved hidden help state")
	}
}

func captureStdout(t *testing.T, run func() int) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	defer reader.Close()

	originalStdout := os.Stdout
	os.Stdout = writer
	rc := run()
	_ = writer.Close()
	os.Stdout = originalStdout

	if rc != 0 {
		t.Fatalf("expected rc 0, got %d", rc)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read captured stdout: %v", err)
	}
	return string(output)
}
