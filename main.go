package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PlanEntry represents a single line item from the plan file.
type PlanEntry struct {
	Exercise  string
	Timestamp string
	StartTime string
	EndTime   string
}

func parsePlan(planPath string) ([]PlanEntry, error) {
	file, err := os.Open(planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plan file: %w", err)
	}
	defer file.Close()

	var entries []PlanEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.TrimSpace(strings.ToLower(line)) == "exercise|timestamp" {
			continue
		}

		record := strings.Split(line, "|")
		if len(record) != 2 {
			continue
		}

		exercise := strings.TrimSpace(record[0])
		timestamp := strings.TrimSpace(record[1])

		parts := strings.Split(timestamp, "-")
		if len(parts) != 2 {
			continue
		}
		startTime := strings.TrimSpace(parts[0])
		endTime := strings.TrimSpace(parts[1])

		entries = append(entries, PlanEntry{
			Exercise:  exercise,
			Timestamp: timestamp,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading plan file: %w", err)
	}

	return entries, nil
}

func generateFFmpegCmd(videoPath, videoDir string, entry PlanEntry, index int) *exec.Cmd {
	outputFileName := fmt.Sprintf("%02d_%s.mp4", index, strings.ReplaceAll(entry.Exercise, " ", "_"))
	outputPath := filepath.Join(videoDir, outputFileName)

	return exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", entry.StartTime,
		"-to", entry.EndTime,
		"-c", "copy",
		outputPath,
	)
}

func main() {
	videoPath := flag.String("video", "", "path to the video file")
	planPath := flag.String("plan", "", "path to the plan CSV file")
	flag.Parse()

	if *videoPath == "" || *planPath == "" {
		log.Fatal("both --video and --plan flags are required")
	}

	entries, err := parsePlan(*planPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		log.Fatalf("No valid entries found in plan file: %s", *planPath)
	}

	videoDir := filepath.Dir(*videoPath)
	for i, entry := range entries {
		cmd := generateFFmpegCmd(*videoPath, videoDir, entry, i+1)

		fmt.Printf("Executing: %s\n", strings.Join(cmd.Args, " "))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("failed to execute ffmpeg for exercise '%s'. Error: %s\n", entry.Exercise, err)
			fmt.Printf("ffmpeg output: %s\n", string(output))
		} else {
			outputFileName := fmt.Sprintf("%02d_%s.mp4", i+1, strings.ReplaceAll(entry.Exercise, " ", "_"))
			outputPath := filepath.Join(videoDir, outputFileName)
			fmt.Printf("Save fragment %s %s into %s\n", entry.Timestamp, outputFileName, outputPath)
		}
	}

	fmt.Println("Video splitting complete.")
}
