package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func renderBlock(title string, lines []string, c *color.Color) string {
	var b strings.Builder

	// 各ブロックの幅を動的に計算（色付けを無視して計算）
	maxWidth := len(title)
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 2 // 両端のスペース分を追加

	// ブロックの上部（色なし）
	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))

	// タイトル行（色付き）
	b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, title))

	// 内容行（色付き）
	for _, line := range lines {
		b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, line))
	}

	// ブロックの下部（色なし）
	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))

	return b.String()
}

func printHorizontalBlocks(blocks []string) {
	if len(blocks) == 0 {
		return
	}

	// 各ブロックの行と幅を計算
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
			// ANSIエスケープシーケンスを無視して幅を計算
			visibleLength := len(stripANSI(line))
			if visibleLength > maxWidths[i] {
				maxWidths[i] = visibleLength
			}
		}
	}

	// 各行を描画
	for l := 0; l < maxLines; l++ {
		for i := 0; i < len(lines); i++ {
			if l < len(lines[i]) {
				// 行を描画（色付きのまま）
				fmt.Print(lines[i][l])
				// 空白を追加して幅を揃える
				fmt.Print(strings.Repeat(" ", maxWidths[i]-len(stripANSI(lines[i][l]))))
			} else {
				// 空行の場合は空白を出力
				fmt.Print(strings.Repeat(" ", maxWidths[i]))
			}
			fmt.Print("  ") // ブロック間のスペース
		}
		fmt.Println()
	}
	fmt.Println()
}

// ANSIエスケープシーケンスを取り除く正規表現
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}
