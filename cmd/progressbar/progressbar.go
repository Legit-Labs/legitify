package progressbar

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/Legit-Labs/legitify/cmd/tty"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type ProgressBar struct {
	metadata map[namespace.Namespace]collectors.Metadata
	disabled bool
}

func NewProgressBar(md map[namespace.Namespace]collectors.Metadata) *ProgressBar {
	return &ProgressBar{
		metadata: md,
		disabled: !tty.IsStderrTty(),
	}
}

func createBar(progress *mpb.Progress, totalCount int, displayName string) *mpb.Bar {
	return progress.AddBar(int64(totalCount),
		mpb.PrependDecorators(
			decor.Name(displayName, decor.WC{W: len(displayName) + 1, C: decor.DSyncSpaceR}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)
}

func (pb *ProgressBar) createBars() (*mpb.Progress, map[string]*mpb.Bar) {
	var wg sync.WaitGroup
	var outputFile io.Writer
	if pb.disabled {
		outputFile = io.Discard
	} else {
		outputFile = os.Stderr
	}

	bars := make(map[string]*mpb.Bar)
	p := mpb.New(mpb.WithWaitGroup(&wg),
		mpb.WithWidth(64),
		mpb.WithOutput(outputFile))

	for ns, md := range pb.metadata {
		if md.TotalEntities > 0 {
			bars[ns] = createBar(p, md.TotalEntities, ns)
		}
	}

	return p, bars
}

func (pb *ProgressBar) Run(progress <-chan collectors.CollectionMetric) group_waiter.Waitable {
	if pb.disabled {
		fmt.Fprintf(os.Stderr, "Progress bar is disabled because stderr is not a terminal. Starting collection...\n")
	}

	p, bars := pb.createBars()
	go func() {
		for data := range progress {
			displayName := data.Namespace
			val, ok := bars[displayName]

			if ok {
				if data.CollectionChange != 0 {
					val.IncrBy(data.CollectionChange)
				}

				if data.Finished {
					val.SetTotal(int64(pb.metadata[displayName].TotalEntities), true)
					if pb.disabled {
						fmt.Fprintf(os.Stderr, "Finished collecting %s\n", displayName)
					}
				}
			} else {
				log.Printf("Failed to find bar with name: %s (%v)", displayName, data)
			}
		}
	}()

	return p
}
