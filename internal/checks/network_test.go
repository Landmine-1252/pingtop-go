package checks

import (
	"strings"
	"testing"
	"time"
)

func TestPingRunnerBuildsWindowsLinuxAndDarwinCommands(t *testing.T) {
	runner := NewPingRunner()
	runner.goos = "windows"
	windows := runner.buildCommand("1.1.1.1", 1200)
	expectedWindows := []string{"ping", "-n", "1", "-w", "1200", "1.1.1.1"}
	if !equalStrings(windows, expectedWindows) {
		t.Fatalf("unexpected windows command: %#v", windows)
	}

	runner.goos = "linux"
	linux := runner.buildCommand("1.1.1.1", 1200)
	expectedLinux := []string{"ping", "-n", "-c", "1", "-W", "2", "1.1.1.1"}
	if !equalStrings(linux, expectedLinux) {
		t.Fatalf("unexpected linux command: %#v", linux)
	}

	runner.goos = "darwin"
	darwin := runner.buildCommand("1.1.1.1", 1200)
	expectedDarwin := []string{"ping", "-n", "-c", "1", "-W", "1200", "1.1.1.1"}
	if !equalStrings(darwin, expectedDarwin) {
		t.Fatalf("unexpected darwin command: %#v", darwin)
	}
}

func TestDNSResolverTimesOutWithoutSpawningDuplicateLookups(t *testing.T) {
	calls := make([]string, 0, 1)
	resolver := NewDNSResolver(func(hostname string) (bool, string, string) {
		calls = append(calls, hostname)
		time.Sleep(200 * time.Millisecond)
		return false, "", "slow failure"
	})

	ok, _, errMessage := resolver.Resolve("hhh", 50)
	if ok || !strings.Contains(errMessage, "exceeded 50 ms timeout") {
		t.Fatalf("unexpected first resolve result: ok=%v err=%q", ok, errMessage)
	}
	ok, _, errMessage = resolver.Resolve("hhh", 50)
	if ok || !strings.Contains(errMessage, "still pending") {
		t.Fatalf("unexpected second resolve result: ok=%v err=%q", ok, errMessage)
	}
	if len(calls) != 1 || calls[0] != "hhh" {
		t.Fatalf("expected one lookup call, got %#v", calls)
	}
	time.Sleep(250 * time.Millisecond)
	ok, _, errMessage = resolver.Resolve("hhh", 50)
	if ok || errMessage != "slow failure" {
		t.Fatalf("unexpected final resolve result: ok=%v err=%q", ok, errMessage)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
