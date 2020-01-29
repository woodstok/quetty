package quetty

import (
	"io"
	"log"
	"os"
)

func Run(opts *Options) {
	f, err := os.OpenFile("quetty.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	tmuxClient := NewTmuxClient()
	curWindow, err := tmuxClient.CurrentWindow()
	if err != nil {
		log.Fatal(err)
	}

	paneList, err := tmuxClient.ListWindowPanes(curWindow)
	if err != nil {
		log.Fatal(err)
	}
	pipeList := []io.Reader{}
	for _, paneId := range paneList {
		panePipeR, panePipeW := io.Pipe()
		pipeList = append(pipeList, panePipeR)
		go tmuxClient.WriteFromPane(paneId, panePipeW)
	}
	combinedReader, combinedWriter := io.Pipe()

	go func() {
		io.Copy(combinedWriter, io.MultiReader(pipeList...))
		combinedWriter.Close()
	}()

	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(combinedReader)
	if err != nil {
		log.Fatal(err)
	}

	printer, err := NewPrinter(opts)
	if err != nil {
		log.Fatal(err)
	}
	printer.Print(os.Stdout, tokens)
}
