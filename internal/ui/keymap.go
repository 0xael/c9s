package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit       key.Binding
	Help       key.Binding
	Tab1       key.Binding
	Tab2       key.Binding
	Tab3       key.Binding
	Tab4       key.Binding
	NextTab    key.Binding
	PrevTab    key.Binding
	Up         key.Binding
	Down       key.Binding
	Enter      key.Binding
	Back       key.Binding
	Filter     key.Binding
	Refresh    key.Binding
	Transcript key.Binding
	Dispatch   key.Binding
	Stop       key.Binding
	Remove     key.Binding
	Respawn    key.Binding
	Attach     key.Binding
	Follow     key.Binding
	Top        key.Binding
	Bottom     key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Yank       key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Tab1: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "sessions"),
	),
	Tab2: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "jobs"),
	),
	Tab3: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "dashboard"),
	),
	Tab4: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "history"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Transcript: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "transcript"),
	),
	Dispatch: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "dispatch"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove"),
	),
	Respawn: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "respawn"),
	),
	Attach: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "attach"),
	),
	Follow: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "follow"),
	),
	Top: key.NewBinding(
		key.WithKeys("g", "home"),
		key.WithHelp("g", "top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G", "end"),
		key.WithHelp("G", "bottom"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+u", "pgup"),
		key.WithHelp("ctrl+u", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+d", "pgdown"),
		key.WithHelp("ctrl+d", "page down"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy id"),
	),
}
