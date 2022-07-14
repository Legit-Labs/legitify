package progressbar

import (
	"log"
	"os"
	"sync"

	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type ProgressBar struct {
	metadata map[namespace.Namespace]collectors.Metadata
}

func NewProgressBar(md map[namespace.Namespace]collectors.Metadata) *ProgressBar {
	return &ProgressBar{
		metadata: md,
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
	bars := make(map[string]*mpb.Bar)
	p := mpb.New(mpb.WithWaitGroup(&wg),
		mpb.WithWidth(64),
		mpb.WithOutput(os.Stderr))

	for ns, md := range pb.metadata {
		bars[ns] = createBar(p, md.TotalEntities, ns)
	}

	return p, bars
}

func (pb *ProgressBar) Run(progress <-chan collectors.CollectionMetric) group_waiter.Waitable {
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
				}
			} else {
				log.Printf("Failed to find bar with name: %s", displayName)
			}
		}
	}()

	return p
}
