package tui

import (
	"github.com/AnatolyRugalev/kube-commander/internal/theme"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type MenuList struct {
	*widgets.List
	items        []menuItemFunc
	selectedItem Pane
}

type menuItemFunc func() Pane

type menuItem struct {
	name     string
	itemFunc menuItemFunc
}

var items = []menuItem{
	{"Namespaces", func() Pane {
		return NewNamespacesTable()
	}},
	{"Nodes", func() Pane {
		return NewNodesTable()
	}},
	{"PVs", func() Pane {
		return NewPVsTable()
	}},
}

func NewMenuList() *MenuList {
	ml := &MenuList{
		List: widgets.NewList(),
	}
	ml.Title = "Cluster"
	ml.TitleStyle = theme.Theme["title"].Inactive
	ml.BorderStyle = theme.Theme["grid"].Inactive
	ml.TextStyle = theme.Theme["listItem"].Inactive
	ml.SelectedRowStyle = theme.Theme["listItemSelected"].Inactive
	ml.WrapText = false
	for _, item := range items {
		ml.Rows = append(ml.Rows, item.name)
		ml.items = append(ml.items, item.itemFunc)
	}
	ml.SelectedRow = 0
	return ml
}

func (ml *MenuList) OnEvent(event *ui.Event) bool {
	switch event.ID {
	case "<Down>":
		if ml.SelectedRow >= len(ml.Rows)-1 {
			return false
		}
		ml.CursorDown()
		return true
	case "<Up>":
		if ml.SelectedRow <= 0 {
			return false
		}
		ml.CursorUp()
		return true
	case "<Right>", "<Enter>":
		ml.activateCurrent()
		return true
	}
	return false
}

func (ml *MenuList) CursorDown() {
	ml.SelectedRow += 1
	ml.onCursorMove()
}

func (ml *MenuList) CursorUp() {
	ml.SelectedRow -= 1
	ml.onCursorMove()
}

func (ml *MenuList) onCursorMove() {
	ml.selectedItem = ml.items[ml.SelectedRow]()
	screen.ReplaceRightPane(ml.selectedItem)
}

func (ml *MenuList) activateCurrent() {
	if ml.selectedItem != nil {
		screen.Focus(ml.selectedItem)
	}
}

func (ml *MenuList) OnFocusIn() {
	ml.TitleStyle = theme.Theme["title"].Active
	ml.BorderStyle = theme.Theme["grid"].Active
	ml.TextStyle = theme.Theme["listItem"].Active
	ml.SelectedRowStyle = theme.Theme["listItemSelected"].Active
}

func (ml *MenuList) OnFocusOut() {
	ml.TitleStyle = theme.Theme["title"].Inactive
	ml.BorderStyle = theme.Theme["grid"].Inactive
	ml.TextStyle = theme.Theme["listItem"].Inactive
	ml.SelectedRowStyle = theme.Theme["listItemSelected"].Inactive
}
