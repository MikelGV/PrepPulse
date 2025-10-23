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
        Background(lipgloss.Color("b7a67b")).
        Padding(0, 1)
    evenRowStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#8ecae6"))
    oddRowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#A0A0A0"))
)

func RenderDf(df *dataframe.DataFrame, width, height int) string {
    if df == nil || df.Nrow() == 0 {
        return lipgloss.NewStyle().Foreground(lipgloss.Color("#de083a")).Render("No data found")
    }
    
    headers := df.Names()

    colsWidth := calculateColumns(df, width, len(headers))

    var headerCells []string

    for i, h := range headers {
        cell := truncate(h, colsWidth[i])
        headerCells = append(headerCells, cell)
    }

    header := lipgloss.JoinHorizontal(lipgloss.Center, headerCells...)
    
    var rows []string
    maxRow := height - 4

    for i := 0; i < df.Nrow() &&  i < maxRow; i++ {
        var cells []string
        style := evenRowStyle

        if i%2 == 1 {
            style = oddRowStyle
        }

        for j := 0; j < df.Ncol(); j++ {
            val := fmt.Sprintf("%v ", df.Elem(i, j))
            cell := truncate(val, colsWidth[j])
            cells = append(cells, style.Render(cell))
        }

        row := lipgloss.JoinHorizontal(lipgloss.Center, cells...)
        rows = append(rows, row)
    }

    content := lipgloss.JoinVertical(lipgloss.Center, header, strings.Join(rows, "\n"))

    return content
}

func calculateColumns(df *dataframe.DataFrame, totalWidth, numCols int) []int {
    width := make([]int, numCols)

    for i, h := range df.Names() {
        width[i] = len(h)
    }

    for i := 0; i < df.Nrow(); i++ {
        for j := 0; j < df.Ncol(); j++ {
            s := fmt.Sprintf("%v ", df.Elem(i, j))

            if len(s) > width[j] {
                width[j] = len(s)
            }
        }
    }

    available := (totalWidth-4) / numCols

    for i := range width {
        if width[i] > available {
            width[i] = available
        }
    }

    return width
}

func truncate(s string, maxWidth int) string {
    if len(s) <= maxWidth {
        return s
    }
    if maxWidth < 3 {
        return s[:maxWidth]
    }

    return s[:maxWidth] + "..."
}
