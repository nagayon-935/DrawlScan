package utils

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/fatih/color"
)

func TestRenderBlock(t *testing.T) {
	tests := []struct {
		name  string
		title string
		lines []string
		color *color.Color
		want  string
	}{
		{
			name:  "basic",
			title: "Title",
			lines: []string{"line1", "line2"},
			color: color.New(color.FgWhite),
			want: "+---------+\n" +
				"| Title   |\n" +
				"| line1   |\n" +
				"| line2   |\n" +
				"+---------+\n",
		},
		{
			name:  "longest line",
			title: "T",
			lines: []string{"short", "this is longest"},
			color: color.New(color.FgWhite),
			want: "+-------------------+\n" +
				"| T                 |\n" +
				"| short             |\n" +
				"| this is longest   |\n" +
				"+-------------------+\n",
		},
		{
			name:  "no lines",
			title: "Empty",
			lines: []string{},
			color: color.New(color.FgWhite),
			want: "+---------+\n" +
				"| Empty   |\n" +
				"+---------+\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderBlock(tt.title, tt.lines, tt.color)
			// 色なしで比較
			got = stripANSI(got)
			if got != tt.want {
				t.Errorf("RenderBlock() = \n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}

func TestRenderBlock_MultibyteWidth(t *testing.T) {
	// "日本語" is 3 runes but 9 bytes in UTF-8; width must be based on
	// visible rune count, not byte length, or the box is padded far wider
	// than the widest visible line requires.
	title := "日本語"
	lines := []string{"ab"}

	contentWidth := utf8.RuneCountInString(title) // widest line: 3 runes
	maxWidth := contentWidth + 2

	border := "+-" + strings.Repeat("-", maxWidth) + "-+\n"
	want := border +
		fmt.Sprintf("| %-*s |\n", maxWidth, title) +
		fmt.Sprintf("| %-*s |\n", maxWidth, lines[0]) +
		border

	got := stripANSI(RenderBlock(title, lines, color.New(color.FgWhite)))
	if got != want {
		t.Errorf("RenderBlock() = \n%q\nwant:\n%q", got, want)
	}
}

func TestPrintHorizontalBlocks(t *testing.T) {
	// この関数は出力を目視確認する用途が主なので、テストではpanicしないことのみ確認
	blocks := []string{
		RenderBlock("A", []string{"a1", "a2"}, color.New(color.FgWhite)),
		RenderBlock("B", []string{"b1", "b2", "b3"}, color.New(color.FgWhite)),
	}
	PrintHorizontalBlocks(blocks)
	PrintHorizontalBlocks([]string{}) // 空でもpanicしない
}

func Test_stripANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no ansi",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:  "with ansi",
			input: "\x1b[31mred\x1b[0m text",
			want:  "red text",
		},
		{
			name:  "multiple ansi",
			input: "\x1b[32mgreen\x1b[0m and \x1b[34mblue\x1b[0m",
			want:  "green and blue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripANSI(tt.input); got != tt.want {
				t.Errorf("stripANSI() = %v, want %v", got, tt.want)
			}
		})
	}
}
