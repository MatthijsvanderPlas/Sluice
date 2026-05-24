package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/MatthijsvanderPlas/sluice/ui"
)

const ApplicationVerion = "v0.1.0"

func main() {
	versionFlag := flag.Bool("version", false, "print application version")
	shortVersionFlag := flag.Bool("v", false, "print application version (shorthand)")

	flag.Parse()

	if *versionFlag || *shortVersionFlag {
		fmt.Printf("Sluice %s\n", ApplicationVerion)
		os.Exit(0)
	}

	stat, err := os.Stdin.Stat()
	if err != nil || (stat.Mode()&os.ModeCharDevice) != 0 {
		fmt.Println("No input detected. Pipe a stream of data into Sluice, e.g.: some-command | sluice")
		os.Exit(1)
	}

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
