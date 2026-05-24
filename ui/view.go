package ui

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	faintStyle = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#ECECEC"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	matchStyle = lipgloss.NewStyle().Background(lipgloss.Color("#F1FA8C")).Foreground(lipgloss.Color("#282A36"))
)

var lastStyleUser lipgloss.Style

func (m ViewModel) View() tea.View {
	snapshot := m.filteredSnapshot()
	if len(snapshot) == 0 {
		return tea.NewView("")
	}

	height := m.screenHeight
	if height == 0 {
		height = len(snapshot)
	}

	start := max(len(snapshot)-m.scrollingOffset-height, 0)
	end := len(snapshot) - m.scrollingOffset

	visible := snapshot[start:end]
	styled := make([]string, len(visible))
	needle := m.filter.Value()
	for i, line := range visible {
		var prevLine string
		if i == 0 {
			prevLine = ""
		} else {
			prevLine = visible[i-1]
		}
		lineStart := "[" + strconv.Itoa(start+i) + "] "
		styled[i] = lineStart + styledLine(line, needle, prevLine)
	}
	s := strings.Join(styled, "\n")

	var c *tea.Cursor
	if m.editMode && !m.filter.VirtualCursor() {
		c = m.filter.Cursor()
		if c != nil {
			c.Y += lipgloss.Height(s) + lipgloss.Height(m.headerView())
		}
	}

	str := lipgloss.JoinVertical(lipgloss.Top, s, m.headerView(), m.filter.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}
	v := tea.NewView(str)
	if m.editMode {
		v.Cursor = c
	} else {
		v.Cursor = nil
	}

	return v
}

func styledLine(line string, needle string, prevLine string) string {
	var styledLine string

	if len(strings.TrimLeft(line, " ")) != len(line) || len(strings.TrimLeft(line, " \t")) != len(line) || len(strings.TrimLeft(line, "-")) != len(line) {
		// Indented lines get the style of the line that came before it
		styledLine = lastStyleUser.Render(line)
		if needle != "" {
			styledLine = highlightMatches(styledLine, needle)
		}
		return styledLine
	}

	if errorLine(prevLine) || warnLine(prevLine) {
		// usually errors and warnings are multiline so we assume the next line is always apart of it
		// this catches cases for us where the initial error is too long and spills into the next styledLine
		// And therefore we do not get an indent
		styledLine = lastStyleUser.Render(line)
		if needle != "" {
			styledLine = highlightMatches(styledLine, needle)
		}
		return styledLine
	}

	if errorLine(line) {
		styledLine = errorStyle.Render(line)
		lastStyleUser = errorStyle
	} else if warnLine(line) {
		styledLine = warnStyle.Render(line)
		lastStyleUser = warnStyle
	} else if debugLine(line) {
		styledLine = faintStyle.Render(line)
		lastStyleUser = faintStyle
	} else {
		styledLine = line
		lastStyleUser = lipgloss.Style{}
	}

	if needle != "" {
		styledLine = highlightMatches(styledLine, needle)
	}

	return styledLine
}

func highlightMatches(s, needle string) string {
	lower := strings.ToLower(s)
	lowerNeedle := strings.ToLower(needle)
	var b strings.Builder
	for {
		idx := strings.Index(lower, lowerNeedle)
		if idx == -1 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:idx])
		b.WriteString(matchStyle.Render(s[idx : idx+len(needle)]))
		s = s[idx+len(needle):]
		lower = lower[idx+len(needle):]
	}
	return b.String()
}

func errorLine(l string) bool {
	return strings.Contains(strings.ToLower(l), "error") ||
		strings.Contains(strings.ToLower(l), "err")
}

func warnLine(l string) bool {
	return strings.Contains(strings.ToLower(l), "warning") ||
		strings.Contains(strings.ToLower(l), "wrn")
}

func debugLine(l string) bool {
	return strings.Contains(strings.ToLower(l), "debug") ||
		strings.Contains(strings.ToLower(l), "dbg")
}

func (m ViewModel) headerView() string {
	return lipgloss.NewStyle().Bold(true).Underline(true).Render("Search logs:")
}

func (m ViewModel) footerView() string {
	return faintStyle.Render("ctrl+c to quit")
}
