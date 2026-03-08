package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pingtop/internal/checks"
	"pingtop/internal/pingtop"
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
