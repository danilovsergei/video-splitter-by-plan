package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParsePlan(t *testing.T) {
	expected := []PlanEntry{
		{Exercise: "forehand", Timestamp: "1:20-5:50", StartTime: "1:20", EndTime: "5:50"},
		{Exercise: "backhand", Timestamp: "7:00-11:15", StartTime: "7:00", EndTime: "11:15"},
		{Exercise: "serve", Timestamp: "1:05:10-1:07:55", StartTime: "1:05:10", EndTime: "1:07:55"},
	}

	entries, err := parsePlan("testdata/plan.csv")
	if err != nil {
		t.Fatalf("parsePlan() error = %v", err)
	}

	if !reflect.DeepEqual(entries, expected) {
		t.Errorf("parsePlan() = %v, want %v", entries, expected)
	}
}

func TestGenerateFFmpegCommandWithTitle(t *testing.T) {
	videoPath := "/tmp/video.mp4"
	videoDir := "/tmp"
	entry := PlanEntry{
		Exercise:  "forehand drill",
		StartTime: "1:20",
		EndTime:   "5:50",
	}
	index := 1
	outputFileName := fmt.Sprintf("%02d - %s.mp4", index, entry.Exercise)
	outputPath := filepath.Join(videoDir, outputFileName)

	expectedCmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", entry.StartTime,
		"-to", entry.EndTime,
		"-c", "copy",
		"-metadata", fmt.Sprintf("title=%s", entry.Exercise),
		outputPath,
	)
	actualCmd := generateFFmpegCmd(videoPath, videoDir, entry, index)

	fmt.Printf("Expected command: %s\n", strings.Join(expectedCmd.Args, " "))
	fmt.Printf("Actual command: %s\n", strings.Join(actualCmd.Args, " "))

	if !reflect.DeepEqual(actualCmd.Args, expectedCmd.Args) {
		t.Errorf("generateFFmpegCmd() args = %v, want %v", actualCmd.Args, expectedCmd.Args)
	}
}
