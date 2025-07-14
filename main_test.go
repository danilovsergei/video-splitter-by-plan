package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParsePlan(t *testing.T) {
	content := `exercise|timestamp
forehand|1:20-5:50
backhand|7:00-11:15
serve|1:05:10-1:07:55
`
	tmpfile := filepath.Join(t.TempDir(), "plan.csv")
	if err := os.WriteFile(tmpfile, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}

	expected := []PlanEntry{
		{Exercise: "forehand", Timestamp: "1:20-5:50", StartTime: "1:20", EndTime: "5:50"},
		{Exercise: "backhand", Timestamp: "7:00-11:15", StartTime: "7:00", EndTime: "11:15"},
		{Exercise: "serve", Timestamp: "1:05:10-1:07:55", StartTime: "1:05:10", EndTime: "1:07:55"},
	}

	entries, err := parsePlan(tmpfile)
	if err != nil {
		t.Fatalf("parsePlan() error = %v", err)
	}

	if !reflect.DeepEqual(entries, expected) {
		t.Errorf("parsePlan() = %v, want %v", entries, expected)
	}
}

func TestGenerateFFmpegCommand(t *testing.T) {
	videoPath := "/tmp/video.mp4"
	videoDir := "/tmp"
	entry := PlanEntry{
		Exercise:  "forehand",
		StartTime: "1:20",
		EndTime:   "5:50",
	}
	index := 1
	outputFileName := fmt.Sprintf("%02d_%s.mp4", index, entry.Exercise)
	outputPath := filepath.Join(videoDir, outputFileName)

	expectedCmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", entry.StartTime,
		"-to", entry.EndTime,
		"-c", "copy",
		outputPath,
	)
	actualCmd := generateFFmpegCmd(videoPath, videoDir, entry, index)

	fmt.Printf("Expected command: %s\n", strings.Join(expectedCmd.Args, " "))
	fmt.Printf("Actual command: %s\n", strings.Join(actualCmd.Args, " "))

	if !reflect.DeepEqual(actualCmd.Args, expectedCmd.Args) {
		t.Errorf("generateFFmpegCmd() args = %v, want %v", actualCmd.Args, expectedCmd.Args)
	}
}
