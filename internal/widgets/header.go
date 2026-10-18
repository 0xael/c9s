package widgets

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/elalwany/c9s/internal/ui"
)

type Header struct {
	tabs      []string
	activeTab int
	width     int
}

func NewHeader(tabs []string) Header {
	return Header{tabs: tabs}
}

func (h Header) SetActive(idx int) Header {
	h.activeTab = idx
	return h
}

func (h Header) SetWidth(w int) Header {
	h.width = w
	return h
}

func (h Header) View() string {
	logo := ui.LogoStyle.Render(" c9s ")

	var tabs []string
	for i, t := range h.tabs {
		label := fmt.Sprintf("%d:%s", i+1, t)
		if i == h.activeTab {
			tabs = append(tabs, ui.TabActive.Render(label))
		} else {
			tabs = append(tabs, ui.TabInactive.Render(label))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)

	hint := ui.HelpDesc.Render("?:help q:quit")

	left := lipgloss.JoinHorizontal(lipgloss.Bottom, logo, " ", tabBar)

	gap := h.width - lipgloss.Width(left) - lipgloss.Width(hint)
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + hint
}
