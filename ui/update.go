package ui

import (
	tea "charm.land/bubbletea/v2"
)

const (
	footerReserve = 3 // Reserver 3 lines for the input and footer info
	headerReserve = 1 // Reserve one line for the header
)

type NewLineMsg string

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenHeight = msg.Height - headerReserve - footerReserve

	case NewLineMsg:
		m.buffer.Add(string(msg))
		if m.scrollingOffset > 0 {
			m.scrollingOffset++
		}

	case tea.KeyPressMsg:
		if m.editMode {
			switch msg.String() {
			case "esc":
				m.editMode = false
				m.filter.Blur()
				return m, nil

			case "ctrl+c":
				m.filter.Blur()
				m.quitting = true
				return m, tea.Quit
			}

			// Pass keys to the filter (edit mode)
			var cmd tea.Cmd
			m.filter, cmd = m.filter.Update(msg)
			m = m.clampScrollOffset()
			return m, cmd

		} else {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "up", "k":
				m = m.scrollUp(1)
				return m, nil

			case "down", "j":
				m = m.scrollDown(1)
				return m, nil

			case "ctrl+d":
				m = m.scrollDown(10)
				return m, nil

			case "ctrl+u":
				m = m.scrollUp(10)
				return m, nil

			case "i":
				m.editMode = true
				m.filter.Focus()
				return m, nil

			case "g":
				// go to top
				m.scrollingOffset = max(len(m.filteredSnapshot()) - m.screenHeight)
				return m, nil

			case "G":
				m.scrollingOffset = 0
				return m, nil

			case "/":
				m.editMode = true
				m.scrollingOffset = 0
				m.filter.SetValue("")
				m.filter.Focus()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	m = m.clampScrollOffset()
	return m, cmd
}

func (m ViewModel) clampScrollOffset() ViewModel {
	maxOffset := max(len(m.filteredSnapshot())-m.screenHeight, 0)
	if m.scrollingOffset > maxOffset {
		m.scrollingOffset = maxOffset
	}
	if m.scrollingOffset < 0 {
		m.scrollingOffset = 0
	}

	return m
}

func (m ViewModel) scrollDown(n int) ViewModel {
	if m.scrollingOffset > 0 {
		m.scrollingOffset = max(m.scrollingOffset-n, 0)
	}

	return m
}

func (m ViewModel) scrollUp(n int) ViewModel {
	snapshot := m.filteredSnapshot()
	if m.scrollingOffset < (len(snapshot) - m.screenHeight) {
		m.scrollingOffset = min(m.scrollingOffset+n, len(snapshot)-m.screenHeight)
	}

	return m
}
