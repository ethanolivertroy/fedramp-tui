package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Help     key.Binding
	Home     key.Binding
	Reqs     key.Binding
	Defs     key.Binding
	Indicators key.Binding
	Filter   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("k/up", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("j/down", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/view details"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "go back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Home: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "documents"),
		),
		Reqs: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "requirements"),
		),
		Defs: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "definitions"),
		),
		Indicators: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "indicators"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
	}
}

// ShortHelp returns brief key help
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Quit}
}

// FullHelp returns full key help
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Home, k.Reqs, k.Defs, k.Indicators},
		{k.Filter, k.Help, k.Quit},
	}
}
