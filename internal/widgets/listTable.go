package widgets

import (
	"github.com/AnatolyRugalev/kube-commander/internal/theme"
	ui "github.com/gizak/termui/v3"
	"image"
	"unicode/utf8"
)

type ListTable struct {
	*ui.Block
	Rows             [][]string
	ColumnWidths     []int
	RowSeparator     bool
	TextAlignment    ui.Alignment
	RowStyles        map[int]ui.Style
	FillRow          bool
	topRow           int
	SelectedRowStyle ui.Style
	SelectedRow      int

	// ColumnResizer is called on each Draw. Can be used for custom column sizing.
	ColumnResizer func()

	handler     ListTableHandler
	RowStyle    ui.Style
	HeaderStyle ui.Style
}

type ListTableHandler interface {
	GetHeaderRow() []string
	LoadData() ([][]string, error)
}

type ListTableEventable interface {
	ListTableHandler
	OnEvent(event *ui.Event, item []string) bool
}

type ListTableSelectable interface {
	ListTableHandler
	OnSelect(item []string) bool
}

type ListTableDeletable interface {
	ListTableHandler
	OnDelete(item []string) bool
}

func NewListTable(handler ListTableHandler) *ListTable {
	lt := &ListTable{
		Block:            ui.NewBlock(),
		RowSeparator:     false,
		RowStyles:        make(map[int]ui.Style),
		ColumnResizer:    func() {},
		handler:          handler,
		RowStyle:         theme.Theme["listItem"].Inactive,
		HeaderStyle:      theme.Theme["listHeader"].Inactive,
		SelectedRowStyle: theme.Theme["listItemSelected"].Inactive,
		FillRow:          true,
	}
	lt.BorderStyle = theme.Theme["grid"].Inactive
	lt.TitleStyle = theme.Theme["title"].Inactive
	lt.ColumnResizer = func() {
		rows := append([][]string{lt.handler.GetHeaderRow()}, lt.Rows...)
		colCount := len(rows[0])
		var widths []int
		for i := range rows[0] {
			var width = 1
			if i == colCount-1 {
				// Last column
				width = 999
			} else {
				for _, row := range rows {
					if utf8.RuneCountInString(row[i]) > width {
						width = len(row[i])
					}
				}
			}
			widths = append(widths, width+1)
		}
		lt.ColumnWidths = widths
	}
	return lt
}

func (lt *ListTable) Draw(buf *ui.Buffer) {
	for i := range lt.Rows {
		if i == lt.SelectedRow {
			lt.RowStyles[i] = lt.SelectedRowStyle
		} else {
			lt.RowStyles[i] = lt.RowStyle
		}
	}
	lt.Block.Draw(buf)

	lt.ColumnResizer()

	columnWidths := lt.ColumnWidths
	if len(columnWidths) == 0 {
		columnCount := len(lt.handler.GetHeaderRow())
		columnWidth := lt.Inner.Dx() / columnCount
		for i := 0; i < columnCount; i++ {
			columnWidths = append(columnWidths, columnWidth)
		}
	}

	// adjusts view into widget
	if lt.SelectedRow >= lt.Inner.Dy()+lt.topRow-1 {
		viewport := lt.Inner.Dy() - 2
		lt.topRow = lt.SelectedRow - viewport
	} else if lt.SelectedRow < lt.topRow {
		lt.topRow = lt.SelectedRow
	}

	// draw header
	yCoordinate := lt.drawRow(buf, columnWidths, lt.handler.GetHeaderRow(), lt.HeaderStyle, lt.Inner.Min.Y)

	// draw rows
	for i := lt.topRow; i < len(lt.Rows) && yCoordinate < lt.Inner.Max.Y; i++ {
		rowStyle := lt.RowStyle
		if style, ok := lt.RowStyles[i]; ok {
			rowStyle = style
		}
		yCoordinate = lt.drawRow(buf, columnWidths, lt.Rows[i], rowStyle, yCoordinate)
	}

	// draw UP_ARROW if needed
	if lt.topRow > 0 {
		buf.SetCell(
			ui.NewCell(ui.UP_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(lt.Inner.Max.X-1, lt.Inner.Min.Y+1),
		)
	}

	// draw DOWN_ARROW if needed
	if len(lt.Rows) > int(lt.topRow)+lt.Inner.Dy() {
		buf.SetCell(
			ui.NewCell(ui.DOWN_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(lt.Inner.Max.X-1, lt.Inner.Max.Y-1),
		)
	}
}

func (lt *ListTable) drawRow(buf *ui.Buffer, columnWidths []int, row []string, rowStyle ui.Style, yCoordinate int) int {
	if lt.FillRow {
		blankCell := ui.NewCell(' ', rowStyle)
		buf.Fill(blankCell, image.Rect(lt.Inner.Min.X, yCoordinate, lt.Inner.Max.X, yCoordinate+1))
	}

	colXCoordinate := lt.Inner.Min.X
	// draw row cells
	for j := 0; j < len(row); j++ {
		col := ui.ParseStyles(row[j], rowStyle)
		// draw row cell
		if len(col) > columnWidths[j] || lt.TextAlignment == ui.AlignLeft {
			for _, cx := range ui.BuildCellWithXArray(col) {
				k, cell := cx.X, cx.Cell
				if k == columnWidths[j] || colXCoordinate+k == lt.Inner.Max.X {
					cell.Rune = ui.ELLIPSES
					buf.SetCell(cell, image.Pt(colXCoordinate+k-1, yCoordinate))
					break
				} else {
					buf.SetCell(cell, image.Pt(colXCoordinate+k, yCoordinate))
				}
			}
		} else if lt.TextAlignment == ui.AlignCenter {
			xCoordinateOffset := (columnWidths[j] - len(col)) / 2
			stringXCoordinate := xCoordinateOffset + colXCoordinate
			for _, cx := range ui.BuildCellWithXArray(col) {
				k, cell := cx.X, cx.Cell
				buf.SetCell(cell, image.Pt(stringXCoordinate+k, yCoordinate))
			}
		} else if lt.TextAlignment == ui.AlignRight {
			stringXCoordinate := ui.MinInt(colXCoordinate+columnWidths[j], lt.Inner.Max.X) - len(col)
			for _, cx := range ui.BuildCellWithXArray(col) {
				k, cell := cx.X, cx.Cell
				buf.SetCell(cell, image.Pt(stringXCoordinate+k, yCoordinate))
			}
		}
		colXCoordinate += columnWidths[j] + 1
	}

	// draw vertical separators
	separatorStyle := lt.Block.BorderStyle

	separatorXCoordinate := lt.Inner.Min.X
	verticalCell := ui.NewCell(ui.VERTICAL_LINE, separatorStyle)
	for i, width := range columnWidths {
		if lt.FillRow && i < len(columnWidths)-1 {
			verticalCell.Style.Bg = rowStyle.Bg
		} else {
			verticalCell.Style.Bg = lt.Block.BorderStyle.Bg
		}

		separatorXCoordinate += width
		buf.SetCell(verticalCell, image.Pt(separatorXCoordinate, yCoordinate))
		separatorXCoordinate++
	}

	yCoordinate++

	// draw horizontal separator
	horizontalCell := ui.NewCell(ui.HORIZONTAL_LINE, separatorStyle)
	if lt.RowSeparator && yCoordinate < lt.Inner.Max.Y {
		buf.Fill(horizontalCell, image.Rect(lt.Inner.Min.X, yCoordinate, lt.Inner.Max.X, yCoordinate+1))
		yCoordinate++
	}
	return yCoordinate
}

func (lt *ListTable) OnEvent(event *ui.Event) bool {
	switch event.ID {
	case "<Down>":
		if lt.SelectedRow >= len(lt.Rows)-2 {
			return false
		}
		lt.CursorDown()
		return true
	case "<Up>":
		if lt.SelectedRow <= 0 {
			return false
		}
		lt.CursorUp()
		return true
	case "<Enter>":
		if s, ok := lt.handler.(ListTableSelectable); ok {
			row := lt.Rows[lt.SelectedRow+1]
			return s.OnSelect(row)
		}
		return false
	case "<Delete>":
		if d, ok := lt.handler.(ListTableDeletable); ok {
			row := lt.Rows[lt.SelectedRow+1]
			return d.OnDelete(row)
		}
		return false
	}
	if e, ok := lt.handler.(ListTableEventable); ok {
		row := lt.Rows[lt.SelectedRow+1]
		return e.OnEvent(event, row)
	}
	return false
}

func (lt *ListTable) CursorDown() {
	lt.SelectedRow += 1
}

func (lt *ListTable) CursorUp() {
	lt.SelectedRow -= 1
}

func (lt *ListTable) Reload() error {
	data, err := lt.handler.LoadData()
	if err != nil {
		return err
	}
	for _, row := range data {
		lt.Rows = append(lt.Rows, row)
	}
	// If deleting last row
	if lt.SelectedRow >= len(lt.Rows)-1 {
		lt.SelectedRow = len(lt.Rows) - 2
	}
	return nil
}

func (lt *ListTable) OnFocusIn() {
	lt.BorderStyle = theme.Theme["grid"].Active
	lt.TitleStyle = theme.Theme["title"].Active
	lt.RowStyle = theme.Theme["listItem"].Active
	lt.HeaderStyle = theme.Theme["listHeader"].Active
	lt.SelectedRowStyle = theme.Theme["listItemSelected"].Active
}

func (lt *ListTable) OnFocusOut() {
	lt.BorderStyle = theme.Theme["grid"].Inactive
	lt.TitleStyle = theme.Theme["title"].Inactive
	lt.RowStyle = theme.Theme["listItem"].Inactive
	lt.HeaderStyle = theme.Theme["listHeader"].Inactive
	lt.SelectedRowStyle = theme.Theme["listItemSelected"].Inactive
}
