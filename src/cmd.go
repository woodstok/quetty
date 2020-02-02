package quetty

import (
	"io"
	"log"
	"os"
	"sync"
)

const bufferFilePath = "/tmp/quetty-tmux-content"

var logger *log.Logger

func Run(opts *Options) {
	f, err := os.OpenFile("/tmp/quetty.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	logger.Printf("Starting quetty")
	defer f.Close()

	var inputStream io.ReadCloser
	var combinedWg sync.WaitGroup

	if opts.stdin {
		logger.Printf("Reading from stdin")
		inputStream = os.Stdin
	} else if opts.init {
		logger.Printf("Reading from tmux panes")
		tmuxClient := NewTmuxClient()
		curWindow, err := tmuxClient.CurrentWindow()
		if err != nil {
			log.Fatal(err)
		}

		paneList, err := tmuxClient.ListWindowPanes(curWindow)
		if err != nil {
			log.Fatal(err)
		}
		logger.Printf("panelist %v", paneList)

		pipeList := []io.Reader{}
		for _, paneId := range paneList {
			panePipeR, panePipeW := io.Pipe()
			pipeList = append(pipeList, panePipeR)
			go tmuxClient.WriteFromPane(paneId, panePipeW)
		}
		combinedReader, combinedWriter := io.Pipe()

		os.Remove(bufferFilePath)
		logger.Printf("Creating temp pane file")
		tmpFile, err := os.OpenFile(bufferFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			combinedWg.Add(1)
			defer combinedWg.Done()
			io.Copy(io.MultiWriter(combinedWriter, tmpFile), io.MultiReader(pipeList...))
			combinedWriter.Close()
			tmpFile.Close()
			logger.Printf("combined writer done copying")
		}()
		inputStream = combinedReader

	} else {
		logger.Printf("reading from tmux buffer")
		inputReader, err := os.OpenFile(bufferFilePath, os.O_RDONLY, 0644)
		inputStream = inputReader
		if err != nil {
			log.Fatal(err)
		}

	}

	logger.Printf("running with opts %+v", opts)
	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(inputStream)
	if err != nil {
		log.Fatal(err)
	}
	logger.Printf("got %d tokens", len(tokens))
	combinedWg.Wait()

	printer, err := NewPrinter(opts)
	if err != nil {
		log.Fatal(err)
	}

	printer.Print(tokens)
	os.Exit(0)
}
