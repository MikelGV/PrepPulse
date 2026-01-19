package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ColumnsStat struct {
	Name         string
	Type         string
	Min          float64
	Max          float64
	Mean         float64
	Median       float64
	Unique       int
	Missing      int
	Buckets      []int
	BucketLabels []string
}

var (
	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8ecae6")).
			Padding(1, 2)

	summaryHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8ecae6")).
			Width(12)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	barStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffb703"))

	barColors = []lipgloss.Color{
		lipgloss.Color("#FFE8DC"),
		lipgloss.Color("#FFD2B5"),
		lipgloss.Color("#FFB89E"),
		lipgloss.Color("#FF9E87"),
		lipgloss.Color("#FF8470"),
		lipgloss.Color("#FF6A5A"),
		lipgloss.Color("#FF5044"),
		lipgloss.Color("#FF362D"),
		lipgloss.Color("#FF1C16"),
		lipgloss.Color("#FF0200"),
	}

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	axisStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4a90a4"))

	bucketLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8ecae6"))

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8ecae6")).
				Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	suggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#126782"))

	suggestionSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#8ecae6")).
				Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2e7d32"))

	gapStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#126782")).
			Foreground(lipgloss.Color("#126782"))

	tickMarkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8ecae6"))
)

func RenderHistogramPopup(stat ColumnsStat, width, height int) string {
	var sections []string

	title := summaryHeaderStyle.Render(
		fmt.Sprintf("Histogram: %s (%d buckets)", stat.Name, len(stat.Buckets)),
	)
	sections = append(sections, title)
	sections = append(sections, "")

	summary := renderSummary(stat)
	sections = append(sections, summary)

	if stat.Type != "string" {
		histogramSection := renderHistogram(stat, width, height)
		sections = append(sections, histogramSection)
	}

	sections = append(sections, "")
	hint := hintStyle.Render("[ ] = Adjust buckets  |  Esc = Close")
	sections = append(sections, hint)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	popup := popupStyle.Render(content)

	return popup
}

func RenderColumnInputPopup(inputBuffer string, suggestions []string, selectedIndex int, width, height int) string {
	var sections []string

	prompt := inputPromptStyle.Render("Enter column name:")
	inputDisplay := inputBuffer + "_"
	inputRow := inputStyle.Render(inputDisplay)

	sections = append(sections, prompt)
	sections = append(sections, inputRow)
	sections = append(sections, "")

	if len(suggestions) > 0 {
		suggestionHeader := hintStyle.Render("Suggestions (use ↑↓ to navigate, Enter to select):")
		sections = append(sections, suggestionHeader)
		sections = append(sections, "")

		maxSuggestions := 10
		if len(suggestions) < maxSuggestions {
			maxSuggestions = len(suggestions)
		}

		for i := 0; i < maxSuggestions; i++ {
			sugg := suggestions[i]
			if i == selectedIndex {
				sections = append(sections, suggestionSelectedStyle.Render("  > "+sugg))
			} else {
				sections = append(sections, suggestionStyle.Render("    "+sugg))
			}
		}
	} else {
		noMatch := hintStyle.Render("No matching columns found")
		sections = append(sections, noMatch)
	}

	sections = append(sections, "")
	hint := hintStyle.Render("Esc=Cancel")
	sections = append(sections, hint)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	popup := popupStyle.Render(content)

	return popup
}

func renderSummary(stat ColumnsStat) string {
	rows := [][]string{}

	rows = append(rows,
		[]string{"Min", fmt.Sprintf("%.2f", stat.Min)},
		[]string{"Max", fmt.Sprintf("%.2f", stat.Max)},
		[]string{"Mean", fmt.Sprintf("%.2f", stat.Mean)},
		[]string{"Median", fmt.Sprintf("%.2f", stat.Median)},
	)

	var renderedRows []string
	for _, row := range rows {
		label := labelStyle.Render(row[0])
		value := valueStyle.Render(row[1])
		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Left, label, ": ", value))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}

func getGradientColor(count, maxCount int) lipgloss.Color {
	if maxCount == 0 {
		return barColors[0]
	}

	ratio := float64(count) / float64(maxCount)
	index := int(ratio * float64(len(barColors)-1))

	if index < 0 {
		index = 0
	}
	if index >= len(barColors) {
		index = len(barColors) - 1
	}

	return barColors[index]
}

func renderHistogram(stat ColumnsStat, popupWidth, popupHeight int) string {
	if len(stat.Buckets) == 0 || len(stat.BucketLabels) == 0 {
		return valueStyle.Render("No data to display")
	}

	maxCount := 0
	for _, count := range stat.Buckets {
		if count > maxCount {
			maxCount = count
		}
	}

	if maxCount == 0 {
		return valueStyle.Render("No data to display")
	}

	maxBarHeight := popupHeight - 8
	if maxBarHeight < 10 {
		maxBarHeight = 10
	}

	var histogramRows []string

	for row := maxBarHeight; row >= 0; row-- {
		countValue := row * (maxCount / maxBarHeight)
		if row == 0 {
			countValue = 0
		}

		showLabel := false
		for _, count := range stat.Buckets {
			if count == countValue {
				showLabel = true
				break
			}
		}

		var tickLabel string
		if showLabel {
			tickLabel = tickMarkStyle.Render(leftPad(strconv.Itoa(countValue), 5) + " │ ")
		} else {
			tickLabel = leftPad("", 8)
		}

		var bars string
		for bucketIdx, count := range stat.Buckets {
			barHeight := int(float64(count) / float64(maxCount) * float64(maxBarHeight))

			shouldShowBar := row <= barHeight && barHeight > 0

			if bucketIdx > 0 {
				bars += gapStyle.Render("  ")
			}

			if shouldShowBar {
				barColor := getGradientColor(count, maxCount)
				barStyle := lipgloss.NewStyle().Foreground(barColor)
				bars += barStyle.Render("█")
			} else {
				bars += " "
			}
		}

		histogramRows = append(histogramRows, tickLabel+bars)
	}

	var axisLine string
	axisLine = leftPad("", 8)
	for i := range stat.Buckets {
		if i > 0 {
			axisLine += gapStyle.Render("  ")
		}
		axisLine += "─"
	}
	axisLine = tickMarkStyle.Render(axisLine)
	histogramRows = append(histogramRows, axisLine)

	var labelRow string
	labelRow = leftPad("", 8)

	totalGaps := (len(stat.Buckets) - 1) * 2
	availableWidth := popupWidth - 8 - totalGaps
	labelWidth := availableWidth / len(stat.Buckets)
	if labelWidth < 4 {
		labelWidth = 4
	}

	for _, label := range stat.BucketLabels {
		labelRow += bucketLabelStyle.Render(centerInWidth(truncateLabel(label, labelWidth), labelWidth))
		labelRow += gapStyle.Render("  ")
	}
	histogramRows = append(histogramRows, labelRow)

	return lipgloss.JoinVertical(lipgloss.Left, histogramRows...)
}

func centerInWidth(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	padding := width - len(s)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func leftPad(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func truncateLabel(label string, maxLen int) string {
	if len(label) <= maxLen {
		return label
	}
	if maxLen < 4 {
		return label[:maxLen]
	}
	return label[:maxLen-1] + "."
}
