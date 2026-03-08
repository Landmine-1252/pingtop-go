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

func TestHandleKeyEscQuitsWhenNoPromptIsActive(t *testing.T) {
	tempDir := t.TempDir()
	services := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	})
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
	services := buildServices(pingtop.RuntimePaths{
		ConfigPath: filepath.Join(tempDir, "pingtop.json"),
		LogPath:    filepath.Join(tempDir, "pingtop_log.csv"),
	})
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
