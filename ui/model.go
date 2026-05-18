package ui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/MatthijsvanderPlas/GoLogView/buffer"
)

type ViewModel struct {
	buffer          buffer.Buffer
	scrollingOffset int
	screenHeight    int
	filter          textinput.Model
	quitting        bool
	editMode        bool
}

func InitialModel() ViewModel {
	ti := textinput.New()
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 256
	ti.SetWidth(20)

	return ViewModel{
		buffer:   buffer.NewRingBuffer(10_000),
		filter:   ti,
		editMode: true,
	}
}

func (m ViewModel) filteredSnapshot() []string {
	snapshot := m.buffer.Snapshot()
	filterText := strings.ToLower(m.filter.Value())
	if filterText == "" {
		return snapshot
	}

	var result []string
	for _, line := range snapshot {
		if strings.Contains(strings.ToLower(line), filterText) {
			result = append(result, line)
		}
	}
	if len(result) == 0 {
		result = append(result, "No results for filter\n")
	}
	return result
}

func (m ViewModel) Init() tea.Cmd {
	return nil
}
