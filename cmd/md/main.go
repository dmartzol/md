package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if input is being piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "error: no input provided")
		fmt.Fprintln(os.Stderr, "usage: cat input.csv | md-convert")
		os.Exit(1)
	}

	// Read stdin as text using a scanner
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Filter out empty lines
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Create a new CSV reader with the filtered lines
	reader := csv.NewReader(strings.NewReader(strings.Join(filteredLines, "\n")))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading csv: %v\n", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, "error: input csv is empty")
		os.Exit(1)
	}

	out, err := convertCSVToMarkdown(records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to markdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(out)
}

// convertCSVToMarkdown converts a CSV string to a markdown table with aligned columns.
func convertCSVToMarkdown(in [][]string) (string, error) {
	if len(in) < 1 {
		return "", fmt.Errorf("input slice is empty")
	}

	// Calculate the maximum width for each column
	colWidths := make([]int, len(in[0]))
	for _, row := range in {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var result string

	// Add header row
	result += "|"
	for i, header := range in[0] {
		result += " " + header + strings.Repeat(" ", colWidths[i]-len(header)) + " |"
	}
	result += "\n|"

	// Add separator row
	for _, width := range colWidths {
		result += " " + strings.Repeat("-", width) + " |"
	}
	result += "\n"

	// Add data rows
	for _, row := range in[1:] {
		result += "|"
		for i, cell := range row {
			result += " " + cell + strings.Repeat(" ", colWidths[i]-len(cell)) + " |"
		}
		result += "\n"
	}

	return result, nil
}
