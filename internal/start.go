package internal

import (
	"bufio"
	"context"
	"os"
	"time"

	"github.com/speed2exe/siever/internal/filter"
	"golang.org/x/sync/errgroup"
)

var FLUSH_INTERVAL = 200 * time.Millisecond

func Start(ctx context.Context) error {

	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()

	filterFuncs := filter.MakeFilters()

	// filtering stages + output stage
	stages := make([]chan []byte, len(filterFuncs) + 1)
	for i := range stages {
		stages[i] = make(chan []byte, 1) // leave a buffer for other steps
	}

	eg, ctx := errgroup.WithContext(ctx)

	// read from input and stream to first stage
	eg.Go(func() error {
		for {
			line, err := input.ReadBytes('\n')
			if err != nil {
				return err
			}
			stages[0] <- line
		}
	})

	// apply filterFuncs for each stages, except final stage
	for i := range filterFuncs {
		filter := filterFuncs[i]
		currentStage := stages[i]
		nextStage := stages[i+1]
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case line := <-currentStage:
					ok, err := filter(line)
					if err != nil {
						return err
					}
					if ok {
						nextStage <- line
					}
				}
			}
		})
	}

	eg.Go(func() error {
		outputStage := stages[len(stages)-1]
		ticker := time.NewTicker(FLUSH_INTERVAL)
		for {
			select {
			case <- ticker.C:
				if err := output.Flush(); err != nil {
					return err
				}
			case <- ctx.Done():
				return nil
			case line := <-outputStage:
				if _, err := output.Write(line); err != nil {
					return err
				}
			}
		}
	})

	return eg.Wait()
}
