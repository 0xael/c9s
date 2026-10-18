package widgets

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/0xael/c9s/internal/ui"
)

type Confirm struct {
	action ui.ConfirmAction
	id     string
	name   string
	active bool
	cursor int
	width  int
}

func NewConfirm() Confirm {
	return Confirm{}
}

func (c Confirm) Show(action ui.ConfirmAction, id, name string) Confirm {
	c.action = action
	c.id = id
	c.name = name
	c.active = true
	c.cursor = 1
	return c
}

func (c Confirm) Active() bool {
	return c.active
}

func (c Confirm) SetWidth(w int) Confirm {
	c.width = w
	return c
}

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	if !c.active {
		return c, nil
	}

	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return c, nil
	}

	switch km.String() {
	case "y":
		c.active = false
		return c, func() tea.Msg {
			return ui.ConfirmResultMsg{Action: c.action, ID: c.id, Confirmed: true}
		}
	case "n", "esc":
		c.active = false
		return c, func() tea.Msg {
			return ui.ConfirmResultMsg{Action: c.action, ID: c.id, Confirmed: false}
		}
	case "left", "h":
		c.cursor = 0
	case "right", "l":
		c.cursor = 1
	case "enter":
		c.active = false
		return c, func() tea.Msg {
			return ui.ConfirmResultMsg{Action: c.action, ID: c.id, Confirmed: c.cursor == 0}
		}
	}

	return c, nil
}

func (c Confirm) View() string {
	if !c.active {
		return ""
	}

	var verb string
	switch c.action {
	case ui.ConfirmStop:
		verb = "Stop"
	case ui.ConfirmRemove:
		verb = "Remove"
	default:
		verb = "Confirm"
	}

	question := fmt.Sprintf("%s '%s'?", verb, c.name)

	yesStyle := lipgloss.NewStyle().Padding(0, 2)
	noStyle := lipgloss.NewStyle().Padding(0, 2)

	if c.cursor == 0 {
		yesStyle = yesStyle.Bold(true).Foreground(ui.T.Success)
		noStyle = noStyle.Foreground(ui.T.DimText)
	} else {
		yesStyle = yesStyle.Foreground(ui.T.DimText)
		noStyle = noStyle.Bold(true).Foreground(ui.T.Error)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center,
		yesStyle.Render("[Yes]"),
		noStyle.Render("[No]"),
	)

	content := lipgloss.JoinVertical(lipgloss.Center, question, "", buttons)

	box := ui.ConfirmBox.Render(content)

	return lipgloss.Place(c.width, 3+lipgloss.Height(box),
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
