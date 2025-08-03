package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	tagDir string
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Add titles to video files based on their filenames",
	Run: func(cmd *cobra.Command, args []string) {
		if tagDir == "" {
			log.Fatal("--dir flag is required")
		}

		err := filepath.Walk(tagDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp4") {
				title := generateTitleFromFilename(info.Name())
				tempOutputPath := path + ".tmp.mp4"

				ffmpegCmd := exec.Command("ffmpeg",
					"-i", path,
					"-c", "copy",
					"-metadata", fmt.Sprintf("title=%s", title),
					tempOutputPath,
				)

				fmt.Printf("Executing: %s\n", strings.Join(ffmpegCmd.Args, " "))
				output, err := ffmpegCmd.CombinedOutput()
				if err != nil {
					fmt.Printf("Failed to execute ffmpeg for '%s'. Error: %s\n", path, err)
					fmt.Printf("ffmpeg output: %s\n", string(output))
					// Clean up the temporary file if ffmpeg fails
					os.Remove(tempOutputPath)
				} else {
					// On success, replace the original file with the temporary one
					err := os.Rename(tempOutputPath, path)
					if err != nil {
						fmt.Printf("Failed to replace original file '%s': %v\n", path, err)
					} else {
						fmt.Printf("Successfully tagged '%s' with title '%s'\n", path, title)
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Fatalf("Error walking the path %q: %v\n", tagDir, err)
		}
	},
}

func generateTitleFromFilename(filename string) string {
	// Remove .mp4 extension
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Regex to find prefixes like "01_", "02 - ", etc.
	re := regexp.MustCompile(`^\d+\s*[-_]\s*`)
	title = re.ReplaceAllString(title, "")

	return title
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().StringVar(&tagDir, "dir", "", "directory with mp4 files to tag")
}
