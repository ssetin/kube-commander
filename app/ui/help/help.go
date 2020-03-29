package help

import (
	"github.com/AnatolyRugalev/kube-commander/app/focus"
	"github.com/AnatolyRugalev/kube-commander/app/ui/theme"
	"github.com/AnatolyRugalev/kube-commander/commander"
	"github.com/gdamore/tcell/views"
)

type widget struct {
	*views.Text
	*focus.Focusable
}

func (w widget) MaxSize() (int, int) {
	return w.Text.Size()
}

var text = `kube-commander - browse your Kubernetes cluster in a casual way!

Global:
 ?: Shows help dialog             D: Describe selected resource
 Q: Quit                          E: Edit selected resource
 Ctrl+N or F2: Change namespace   Delete: Delete resource (with confirmation)

Navigation:
 ↑↓→←: List navigation
 Enter: Select menu item          Esc: Go back

Resource types navigation:
 Ctrl+P: Pods                      Ctrl+C: Config Maps
 Ctrl+D: Deployments               Ctrl+I: Ingresses

Pods:
 L: Show logs                      Shift+L: Show previous logs
 F: Forward port
 X: eXec /bin/sh or /bin/bash inside container
`

func NewHelpWidget() *widget {
	widget := widget{
		Text:      views.NewText(),
		Focusable: focus.NewFocusable(),
	}
	widget.Text.SetText(text)
	widget.Text.SetStyle(theme.Default)
	return &widget
}

func ShowHelpPopup(workspace commander.Workspace) {
	help := NewHelpWidget()
	workspace.ShowPopup("Help", help)
}
