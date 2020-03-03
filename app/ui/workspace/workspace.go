package workspace

import (
	"github.com/AnatolyRugalev/kube-commander/app/focus"
	"github.com/AnatolyRugalev/kube-commander/app/ui/help"
	"github.com/AnatolyRugalev/kube-commander/app/ui/resourceMenu"
	"github.com/AnatolyRugalev/kube-commander/app/ui/resources/namespace"
	"github.com/AnatolyRugalev/kube-commander/app/ui/widgets/listTable"
	"github.com/AnatolyRugalev/kube-commander/app/ui/widgets/popup"
	"github.com/AnatolyRugalev/kube-commander/commander"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type workspace struct {
	*views.BoxLayout
	focus.Focusable

	container commander.Container
	focus     commander.FocusManager

	popup  commander.Popup
	menu   commander.MenuListView
	widget commander.Widget

	namespace         string
	namespaceResource *commander.Resource

	selectedWidgetId int
}

func (w *workspace) ResourceProvider() commander.ResourceProvider {
	return w.container.ResourceProvider()
}

func (w *workspace) CommandBuilder() commander.CommandBuilder {
	return w.container.CommandBuilder()
}

func (w *workspace) CommandExecutor() commander.CommandExecutor {
	return w.container.CommandExecutor()
}

func (w *workspace) Client() commander.Client {
	return w.container.Client()
}

func (w *workspace) CurrentNamespace() string {
	return w.namespace
}

func (w *workspace) SwitchNamespace(namespace string) {
	w.namespace = namespace
	if r, ok := w.widget.(reloadable); ok {
		r.Reload()
		w.Update()
	}
}

func NewWorkspace(container commander.Container, namespace string) *workspace {
	return &workspace{
		BoxLayout:        views.NewBoxLayout(views.Horizontal),
		container:        container,
		selectedWidgetId: -1,
		namespace:        namespace,
	}
}

func (w *workspace) FocusManager() commander.FocusManager {
	return w.focus
}

func (w *workspace) ShowPopup(widget commander.MaxSizeWidget) {
	if r, ok := widget.(reloadable); ok {
		r.Reload()
	}
	w.popup = popup.NewPopup(w.container.Screen().View(), widget, func() {
		w.popup = nil
		w.Update()
	})
	w.focus.Focus(w.popup)
	w.Update()
}

func (w *workspace) Update() {
	w.container.Screen().Update()
}

func (w *workspace) HandleError(err error) {
	panic(err)
}

func (w workspace) Draw() {
	w.BoxLayout.Draw()
	if w.popup != nil {
		w.popup.Draw()
	}
}

func (w workspace) Resize() {
	w.BoxLayout.Resize()
	if w.popup != nil {
		w.popup.Reposition(w.container.Screen().View())
		w.popup.Resize()
	}
}

func (w *workspace) HandleEvent(e tcell.Event) bool {
	if w.focus.HandleEvent(e, w.popup == nil) {
		return true
	}
	if w.popup != nil {
		return false
	}
	switch ev := e.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlN, tcell.KeyF2:
			namespace.PickNamespace(w, w.namespaceResource, w.SwitchNamespace)
		default:
			if ev.Rune() == '?' {
				help.ShowHelpPopup(w)
				return true
			}
		}
	}
	return false
}

func (w *workspace) Init() error {
	resMap, err := w.ResourceProvider().PreferredResources()
	if err != nil {
		return err
	}
	w.namespaceResource = resMap["Namespace"]

	resMenu, err := resourceMenu.NewResourcesMenu(w, w.onMenuSelect)
	if err != nil {
		return err
	}

	resMenu.SetStyler(w.styler)
	w.menu = resMenu
	w.widget = w.menu.SelectedItem().Widget()
	w.BoxLayout.AddWidget(w.menu, 0.1)
	w.BoxLayout.AddWidget(w.widget, 0.9)
	w.focus = focus.NewFocusManager(w.menu)

	return nil
}

func (w *workspace) styler(list commander.ListView, rowId int, row commander.Row) tcell.Style {
	style := listTable.DefaultStyler(list, rowId, row)
	if rowId != w.menu.SelectedRowId() && rowId == w.selectedWidgetId {
		style = style.Background(tcell.ColorBlueViolet)
	}
	return style
}

type reloadable interface {
	Reload()
}

func (w *workspace) onMenuSelect(itemId int, item commander.MenuItem) bool {
	if item.Widget() != w.widget {
		w.BoxLayout.RemoveWidget(w.widget)
		w.widget = item.Widget()
		w.BoxLayout.AddWidget(w.widget, 0.9)
		w.selectedWidgetId = itemId
	}
	w.focus.Focus(w.widget)

	if r, ok := w.widget.(reloadable); ok {
		r.Reload()
	}

	return true
}