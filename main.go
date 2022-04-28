package main

import (
	"context"
	"io"
	"os"

	"github.com/speed2exe/siever/internal"
)

func main() {
	if err := internal.Start(context.Background()); err != nil {
		if err != io.EOF {
			// TODO: print in red
			os.Stderr.WriteString("error start: " + err.Error())
		}
	}
}
