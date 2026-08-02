package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

// Test output() with simpleCounter
func TestOptOutput(t *testing.T) {
	opt := &Opt{
		KeyPrefix: "test",
	}
	logCounter := &simpleCounter{}
	baseLogCounter := &simpleCounter{}

	logCounter.total = 10
	logCounter.duration = 5

	baseLogCounter.total = 20
	baseLogCounter.duration = 10

	now := time.Now()
	output := opt.output(logCounter, baseLogCounter, now)
	expectedOutput := "log-incr-rate.test_count.log\t2.000000\t" + fmt.Sprint(uint64(now.Unix())) + "\n" +
		"log-incr-rate.test_count.base\t2.000000\t" + fmt.Sprint(uint64(now.Unix())) + "\n" +
		"log-incr-rate.test_rate.log\t1.000000\t" + fmt.Sprint(uint64(now.Unix())) + "\n"

	if output != expectedOutput {
		t.Errorf("output() = %v, want %v", output, expectedOutput)
	}
}

func TestOptOutputZeroDuration(t *testing.T) {
	opt := &Opt{
		KeyPrefix: "test",
	}
	logCounter := &simpleCounter{}
	baseLogCounter := &simpleCounter{}

	logCounter.total = 10
	logCounter.duration = 0

	baseLogCounter.total = 20
	baseLogCounter.duration = 0

	now := time.Now()
	output := opt.output(logCounter, baseLogCounter, now)
	expectedOutput := ""

	if output != expectedOutput {
		t.Errorf("output() = %v, want %v", output, expectedOutput)
	}
}

func appendToFile(t *testing.T, filePath string, content string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("Failed to write to file %s: %v", filePath, err)
	}
}

func TestOptRun(t *testing.T) {
	workdir := t.TempDir()
	// create test log files
	logFile := path.Join(workdir, "test.log")
	baseLogFile := path.Join(workdir, "base.log")

	if err := os.WriteFile(logFile, []byte("test log content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}
	if err := os.WriteFile(baseLogFile, []byte("test base log content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test base log file: %v", err)
	}

	opt := &Opt{
		KeyPrefix:   "test",
		LogFile:     logFile,
		BaseLogFile: baseLogFile,
	}

	output, err := opt.run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if output != "" {
		t.Errorf("Run() output = %v, want empty string", output)
	}

	// add more logs to the log file and run again
	for i := range 5 {
		appendToFile(t, logFile, fmt.Sprintf("new log line %d\n", i))
	}
	for i := range 15 {
		appendToFile(t, baseLogFile, fmt.Sprintf("new base log line %d\n", i))
	}
	time.Sleep(1 * time.Second) // wait for a second to ensure the duration is non-zero

	// run again
	output, err = opt.run()
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if output == "" {
		t.Errorf("Run() output = %v, want non-empty", output)
	}
	if !strings.Contains(output, "log-incr-rate.test_count.log\t5.000000") || !strings.Contains(output, "log-incr-rate.test_count.base\t15.000000") {
		t.Errorf("Run() output = %v, want output to contain metric keys", output)
	}
	if !strings.Contains(output, "log-incr-rate.test_rate.log\t0.3333") {
		t.Errorf("Run() output = %v, want output to contain log-incr-rate.test_rate.log", output)
	}
}
