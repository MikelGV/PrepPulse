package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
    dataframe "github.com/rocketlaunchr/dataframe-go"
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
    selectedCellRow = lipgloss.NewStyle().
        Background(lipgloss.Color("#126782")).
        Foreground(lipgloss.Color("#FFFFFF")).Padding(0, 2)
    editInputCellStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("#ffb703")).
        Foreground(lipgloss.Color("#6969696")).Padding(0, 2)
 
)

func RenderDf(df *dataframe.DataFrame, width, height, cursorRow, cursorCol, scrollRow, scrollCol int, filteredRows []int, editMode bool, editBuffer string) string {
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

    var displayRows []int
    if len(filteredRows) > 0 {
        displayRows = filteredRows
    } else {
        displayRows = make([]int, df.Nrow())
        for i := 0; i < df.Nrow(); i++ {
            displayRows[i] = i
        }
    }


    start := clamp(scrollRow, 0, len(displayRows)-1)
    end := clamp(start + maxRow, 0, len(displayRows))

    for viewIdx := start; viewIdx < end; viewIdx++ {
        realRow := displayRows[viewIdx]


        var cells []string
        currentWidth := 0

        style := evenRowStyle

        if viewIdx%2 == 1 {
            style = oddRowStyle
        }

        for j := scrollCol; j < numCols; j ++ {
            w := colsWidth[j]

            isCursorCell := (viewIdx == cursorRow) && (j == cursorCol)

            if currentWidth+w  > width - 2 {
                break
            }

            val := fmt.Sprintf("%v", df.Elem(realRow, j))

            if editMode && isCursorCell {
                input := editBuffer + "█"
                cell := padRight(input, w)
                cells = append(cells, editInputCellStyle.Render(cell))
            } else if isCursorCell {
                // TODO: i have to fix this becasue for some reason the selectedStyle is not rendering
                cell := padRight(truncate(val, w ), w)
                cells = append(cells, selectedStyle.Render(cell))
            } else {
                cell := padRight(truncate(val, w ), w)
                cells = append(cells, style.Render(cell))
            }
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
