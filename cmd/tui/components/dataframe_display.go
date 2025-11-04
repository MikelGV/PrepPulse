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

func RenderDf(df *dataframe.DataFrame, width, height, scrollRow, scrollCol int, filteredRows []int) string {
    if df == nil || df.Nrow() == 0 {
        return lipgloss.NewStyle().Foreground(lipgloss.Color("#de083a")).Render("No data found")
    }
    
    headers := df.Names()
    numCols := len(headers)
    colsWidth := calculateColumns(df, width, numCols)

    var headerCells []string
    currentWidth := 0

    for i := scrollCol; i < numCols; i++ {
        w := colsWidth[i]

        if currentWidth+w > width-2 {
            break
        }

        cell := padRight(truncate(headers[i], w), w)
        headerCells = append(headerCells, headerStyle.Render(cell))
        currentWidth += w
    }

    header := lipgloss.JoinHorizontal(lipgloss.Left, headerCells...)
    
    var rows []string
    maxRow := height - 4

    for i := scrollRow; i < df.Nrow() && i < maxRow; i++ {
        var cells []string
        currentWidth := 0

        style := evenRowStyle

        if i%2 == 1 {
            style = oddRowStyle
        }
        for j := scrollCol; j < numCols; j ++ {
            w := colsWidth[j]

            if currentWidth+w  > width - 2 {
                break
            }

            val := fmt.Sprintf("%v", df.Elem(i, j))
            cell := padRight(truncate(val, w ), w)
            cells = append(cells, style.Render(cell))
            currentWidth += w
        }


        row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)
        rows = append(rows, row)
    }

    content := lipgloss.JoinVertical(lipgloss.Left, header, strings.Join(rows, "\n"))

    return content
}

func calculateColumns(df *dataframe.DataFrame, totalWidth, numCols int) []int {
    widths := make([]int, numCols)

    for i, h := range df.Names() {
        widths[i] = len(h) + 4
    }

    for i := 0; i < df.Nrow(); i++ {
        for j := 0; j < numCols; j++ {
            s := fmt.Sprintf("%v", df.Elem(i, j))
            if len(s)+2 > widths[j] {
                widths[j] = len(s) + 2
            }
        }

    }


    total := 0
    for _, w := range widths { total += w }
    if total > totalWidth - 2 {
        scale := float64(totalWidth -2) / float64(total)
        for i := range widths {
            newWidths := int(float64(widths[i]) * scale)
            minWidth := len(df.Names()[i]) + 4
            widths[i] = max(newWidths, minWidth)
        }
    }

    return widths
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


func clamp(v, lo, hi int) int {
    if v < lo {
        return lo
    }

    if v > hi {
        return hi
    }

    return v
}
