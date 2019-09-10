package quetty

import (
	"fmt"
	"log"
	"os"
)

func Run(opts *Options) {

	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tokens are %s\n", tokens)

}
