package transcript

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/0xael/c9s/internal/ui"
	"github.com/0xael/c9s/internal/data"
)

type styledLine struct {
	text string
}

func parseEntries(entries []data.TranscriptEntry, width int) []styledLine {
	if width < 20 {
		width = 80
	}
	contentWidth := width - 4

	var lines []styledLine
	for _, e := range entries {
		ls := parseEntry(e, contentWidth)
		lines = append(lines, ls...)
	}
	return lines
}

func parseEntry(e data.TranscriptEntry, width int) []styledLine {
	switch e.Type {
	case "user":
		return parseUserEntry(e, width)
	case "assistant":
		return parseAssistantEntry(e, width)
	case "system":
		return parseSystemEntry(e, width)
	default:
		return nil
	}
}

type messageEnvelope struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type contentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	Content   json.RawMessage `json:"content"`
	ToolUseID string          `json:"tool_use_id"`
}

func parseUserEntry(e data.TranscriptEntry, width int) []styledLine {
	var msg messageEnvelope
	if json.Unmarshal(e.Message, &msg) != nil {
		return nil
	}

	text := extractContent(msg.Content)
	if text != "" {
		header := lipgloss.NewStyle().Bold(true).Foreground(ui.T.Primary).Render("USER")
		var lines []styledLine
		lines = append(lines, styledLine{text: ""})
		lines = append(lines, styledLine{text: header})
		for _, l := range wrapText(text, width) {
			lines = append(lines, styledLine{
				text: lipgloss.NewStyle().Foreground(ui.T.Text).Render("  " + l),
			})
		}
		return lines
	}

	var blocks []contentBlock
	if json.Unmarshal(msg.Content, &blocks) == nil {
		for _, b := range blocks {
			if b.Type == "tool_result" {
				resultStyle := lipgloss.NewStyle().Foreground(ui.T.DimText)
				result := extractToolResultText(b)
				if result != "" {
					summary := result
					maxLen := width - 10
					if maxLen > 3 && len(summary) > maxLen {
						summary = summary[:maxLen-3] + "..."
					}
					return []styledLine{{text: resultStyle.Render("  " + summary)}}
				}
			}
		}
	}

	return nil
}

func parseAssistantEntry(e data.TranscriptEntry, width int) []styledLine {
	var msg messageEnvelope
	if json.Unmarshal(e.Message, &msg) != nil {
		return nil
	}

	header := lipgloss.NewStyle().Bold(true).Foreground(ui.T.Success).Render("ASSISTANT")
	toolStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.T.Warning)
	textStyle := lipgloss.NewStyle().Foreground(ui.T.Text)

	var blocks []contentBlock
	if json.Unmarshal(msg.Content, &blocks) != nil {
		text := extractContent(msg.Content)
		if text == "" {
			return nil
		}
		var lines []styledLine
		lines = append(lines, styledLine{text: ""})
		lines = append(lines, styledLine{text: header})
		for _, l := range wrapText(text, width) {
			lines = append(lines, styledLine{text: textStyle.Render("  " + l)})
		}
		return lines
	}

	var lines []styledLine
	lines = append(lines, styledLine{text: ""})
	lines = append(lines, styledLine{text: header})

	for _, b := range blocks {
		switch b.Type {
		case "text":
			if b.Text == "" {
				continue
			}
			for _, l := range wrapText(b.Text, width) {
				lines = append(lines, styledLine{text: textStyle.Render("  " + l)})
			}
		case "tool_use", "server_tool_use":
			label := toolStyle.Render(fmt.Sprintf("  TOOL: %s", b.Name))
			lines = append(lines, styledLine{text: label})
		}
	}
	return lines
}

func parseSystemEntry(e data.TranscriptEntry, width int) []styledLine {
	var msg messageEnvelope
	if json.Unmarshal(e.Message, &msg) != nil {
		return nil
	}

	text := extractContent(msg.Content)
	if text == "" {
		return nil
	}

	header := lipgloss.NewStyle().Bold(true).Foreground(ui.T.Accent).Render("SYSTEM")
	style := lipgloss.NewStyle().Foreground(ui.T.DimText)

	var lines []styledLine
	lines = append(lines, styledLine{text: ""})
	lines = append(lines, styledLine{text: header})
	for _, l := range wrapText(text, width) {
		lines = append(lines, styledLine{text: style.Render("  " + l)})
	}
	return lines
}

func extractContent(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}

	var blocks []contentBlock
	if json.Unmarshal(raw, &blocks) == nil {
		var parts []string
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				parts = append(parts, b.Text)
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

func extractToolResultText(b contentBlock) string {
	if b.Content == nil {
		return ""
	}
	var s string
	if json.Unmarshal(b.Content, &s) == nil {
		first, _, _ := strings.Cut(s, "\n")
		return first
	}
	var blocks []contentBlock
	if json.Unmarshal(b.Content, &blocks) == nil {
		for _, inner := range blocks {
			if inner.Type == "text" && inner.Text != "" {
				first, _, _ := strings.Cut(inner.Text, "\n")
				return first
			}
		}
	}
	return ""
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 76
	}
	var result []string
	for _, paragraph := range strings.Split(text, "\n") {
		if paragraph == "" {
			result = append(result, "")
			continue
		}
		for len(paragraph) > width {
			cut := width
			if i := strings.LastIndex(paragraph[:cut], " "); i > 0 {
				cut = i
			}
			result = append(result, paragraph[:cut])
			paragraph = strings.TrimLeft(paragraph[cut:], " ")
		}
		if paragraph != "" {
			result = append(result, paragraph)
		}
	}
	return result
}
