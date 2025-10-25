package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-gota/gota/dataframe"
)

var (
    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Padding(0, 2)
    evenRowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#8ecae6")).Padding(0, 2)
    oddRowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#A0A0A0")).Padding(0, 2)
)

func RenderDf(df *dataframe.DataFrame, width, height, scrollRow, scrollCol int) string {
    if df == nil || df.Nrow() == 0 {
        return lipgloss.NewStyle().Foreground(lipgloss.Color("#de083a")).Render("No data found")
    }
    
    headers := df.Names()
    /**
    colsWidth := calculateColumns(df, width, len(headers))
    **/
    numCols := len(headers)
    visibleCols, colsWidth := calculateColumns(df, width, scrollCol)

    var headerCells []string

    for i, h := range headers {
        cell := padRight(truncate(h, colsWidth[i]), colsWidth[i])
        headerCells = append(headerCells, headerStyle.Render(cell))
    }

    header := lipgloss.JoinHorizontal(lipgloss.Left, headerCells...)
    
    var rows []string
    maxRow := height - 4

    for i := 0; i < df.Nrow() &&  i < maxRow; i++ {
        var cells []string
        style := evenRowStyle

        if i%2 == 1 {
            style = oddRowStyle
        }

        for j := 0; j < df.Ncol(); j++ {
            val := fmt.Sprintf("%v", df.Elem(i, j))
            cell := padRight(truncate(val, colsWidth[j]), colsWidth[j])
            cells = append(cells, style.Render(cell))
        }

        row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)
        rows = append(rows, row)
    }

    content := lipgloss.JoinVertical(lipgloss.Left, header, strings.Join(rows, "\n"))

    return content
}

func calculateColumns(df *dataframe.DataFrame, totalWidth, numCols int) ([]int, []int) {
    width := make([]int, numCols)

    for i, h := range df.Names() {
        width[i] = len(h) + 4
    }

    for i := 0; i < df.Nrow(); i++ {
        for j := 0; j < df.Ncol(); j++ {
            s := fmt.Sprintf("%v", df.Elem(i, j))

            if len(s) + 2 > width[j] {
                width[j] = len(s) + 2
            }
        }
    }

    totalUsed := 0
    for _, w := range width {
        totalUsed += w
    }

    if totalUsed > totalWidth-2 {
        scale := float64(totalWidth-2) / float64(totalUsed)
        for i := range width {
            width[i] = max(3, int(float64(width[i]) * scale))
            headersLen := len(df.Names()[i]) + 4

            if width[i] < headersLen {
                width[i] = headersLen
            }
        }
    }

    remaining := totalWidth - 2 - totalUsed

    if remaining > 0 && numCols > 0 {
        extraPerCol := remaining / numCols

        for i := range width {
            width[i] += extraPerCol
        }
    }

    return width
}

func padRight(s string, width int) string {
    if len(s) >= width {
        return s[:width]
    }

    return s + strings.Repeat(" ", width - len(s))
}

func truncate(s string, maxWidth int) string {
    if len(s) <= maxWidth {
        return s
    }
    if maxWidth < 3 {
        return s[:maxWidth]
    }

    return s[:maxWidth-1] + "..."
}

