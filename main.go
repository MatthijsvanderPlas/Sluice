package main

import (
	"bufio"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"example.com/GoLogView/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			p.Send(ui.NewLineMsg(scanner.Text()))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there was an error: %v", err)
	}
}
