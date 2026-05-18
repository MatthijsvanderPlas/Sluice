package main

import (
	"bufio"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/MatthijsvanderPlas/sluice/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			p.Send(ui.NewLineMsg(scanner.Text()))
		}
		err := scanner.Err()
		if err != nil {
			p.Send(ui.NewLineMsg("ERROR reading stdin: " + err.Error()))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there was an error: %v", err)
	}
}
