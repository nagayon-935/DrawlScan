package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

func RenderBlock(title string, lines []string, c *color.Color) string {
	var b strings.Builder

	maxWidth := utf8.RuneCountInString(title)
	for _, line := range lines {
		if w := utf8.RuneCountInString(line); w > maxWidth {
			maxWidth = w
		}
	}
	maxWidth += 2

	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))
	b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, title))

	for _, line := range lines {
		b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, line))
	}

	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))

	return b.String()
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func PrintHorizontalBlocks(blocks []string) {
	if len(blocks) == 0 {
		return
	}

	lines := make([][]string, len(blocks))
	maxWidths := make([]int, len(blocks))
	maxLines := 0

	for i, b := range blocks {
		blockLines := strings.Split(strings.TrimSuffix(b, "\n"), "\n")
		lines[i] = blockLines
		if len(blockLines) > maxLines {
			maxLines = len(blockLines)
		}
		for _, line := range blockLines {
			visibleLength := utf8.RuneCountInString(stripANSI(line))
			if visibleLength > maxWidths[i] {
				maxWidths[i] = visibleLength
			}
		}
	}

	for l := 0; l < maxLines; l++ {
		for i := 0; i < len(lines); i++ {
			if l < len(lines[i]) {
				fmt.Print(lines[i][l])
				fmt.Print(strings.Repeat(" ", maxWidths[i]-utf8.RuneCountInString(stripANSI(lines[i][l]))))
			} else {
				fmt.Print(strings.Repeat(" ", maxWidths[i]))
			}
			fmt.Print("  ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func stripANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}
