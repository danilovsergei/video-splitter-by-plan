# Bulk Video Splitter By Plan

A command-line tool to split a single video file into multiple smaller fragments based on a plan defined in a CSV file. This is particularly useful for splitting long training videos into individual exercises or segments.

## How it Works

The tool reads a `plan.csv` file that specifies the name of each video fragment and its start and end timestamps. It then uses `ffmpeg` to extract these fragments from the source video and saves them as separate `.mp4` files in the same directory as the source video.

The output files are named with a zero-padded index and the exercise name from the plan file (e.g., `01_forehand.mp4`).

## Prerequisites
- [ffmpeg](https://ffmpeg.org/download.html) (must be installed and available in your system's PATH)

## Build

To build the application, run the following command in the project directory:

```bash
go build
```

This will create an executable file named `video-splitter` (or `video-splitter.exe` on Windows).

## Usage

Run the application from the command line, providing the paths to your video file and your plan file.

```bash
./video-splitter --video="/path/to/your/video.mp4" --plan="/path/to/your/plan.csv"
```

### Arguments

- `--video`: The full path to the source video file.
- `--plan`: The full path to the plan CSV file.

## `plan.csv` Format

The `plan.csv` file defines how the video should be split. It must contain a header and use the `|` character as a delimiter.

- **Header**: The file must start with the header line: `exercise|timestamp`
- **Columns**:
    1.  `exercise`: The name of the exercise or segment. This will be used to name the output file.
    2.  `timestamp`: The start and end times for the fragment, separated by a hyphen (`-`).
- **Timestamp Format**: Timestamps can be in `minutes:seconds` (`MM:SS`) or `hours:minutes:seconds` (`HH:MM:SS`) format.

### Example `plan.csv`

```csv
exercise|timestamp
forehand|1:20-5:50
backhand|7:00-11:15
serve|1:05:10-1:07:55
```

## Example Output

When the application runs, it will print the `ffmpeg` command being executed for each fragment and a success message upon completion.

```
Executing: ffmpeg -i /path/to/your/video.mp4 -ss 1:20 -to 5:50 -c copy /path/to/your/01_forehand.mp4
Save fragment 1:20-5:50 01_forehand.mp4 into /path/to/your/01_forehand.mp4

Executing: ffmpeg -i /path/to/your/video.mp4 -ss 7:00 -to 11:15 -c copy /path/to/your/02_backhand.mp4
Save fragment 7:00-11:15 02_backhand.mp4 into /path/to/your/02_backhand.mp4

Video splitting complete.
