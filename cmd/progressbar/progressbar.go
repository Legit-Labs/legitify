package progressbar

import (
	"io"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/screen"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var pb *progressBar

func init() {
	pb = newProgressBar()
}

func Run() group_waiter.Waitable {
	return pb.Run()
}
func Report(msg ChannelType) {
	pb.ReportProgress(msg)
}

type progressBar struct {
	barTotals map[string]int
	progress  *mpb.Progress
	bars      map[string]*mpb.Bar
	waiter    *pbWaiter
	inChannel chan ChannelType
	enabled   bool
}

func newProgressBar() *progressBar {
	enabled := screen.IsTty()

	var outputFile io.Writer
	if enabled {
		outputFile = screen.Writer()
	} else {
		outputFile = io.Discard
	}

	pb := mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}),
		mpb.WithWidth(64),
		mpb.WithOutput(outputFile),
	)

	waiter := newPbWaiter(pb)

	p := &progressBar{
		barTotals: make(map[string]int),
		bars:      make(map[string]*mpb.Bar),
		progress:  pb,
		waiter:    waiter,
		enabled:   enabled,
		inChannel: make(chan ChannelType),
	}

	return p
}

func (pb *progressBar) Run() group_waiter.Waitable {
	if !pb.enabled {
		screen.Printf("Progress bar is disabled because stderr is not a terminal. Starting collection...\n")
	}

	go func() {
		for d := range pb.inChannel {
			switch data := d.(type) {
			case MinimalBars:
				pb.handleMinimalBars(data)
			case RequiredBarCreation:
				pb.handleRequiredBarCreation(data)
			case OptionalBarCreation:
				pb.handleOptionalBarCreation(data)
			case BarUpdate:
				pb.handleBarUpdate(data)
			case TimedBarCreation:
				pb.handleTimedBarCreation(data)
			case BarClose:
				pb.handleBarClose(data)
			default:
				log.Panicf("unexpected progress update type: %t", d)
			}
		}
	}()

	return pb.waiter
}

func (pb *progressBar) ReportProgress(msg ChannelType) {
	pb.inChannel <- msg
}

func (pb *progressBar) handleMinimalBars(data MinimalBars) {
	pb.waiter.SetMinCount(data.count)
}

func (pb *progressBar) handleRequiredBarCreation(data RequiredBarCreation) {
	pb.handleOptionalBarCreation(OptionalBarCreation(data))
	pb.waiter.ReportBarCreation()
}

func (pb *progressBar) handleOptionalBarCreation(data OptionalBarCreation) {
	if data.TotalEntities == 0 {
		return
	}

	displayName := data.BarName

	if _, exists := pb.bars[displayName]; exists {
		log.Panicf("trying to create a bar that already exists: %s (%v)", displayName, data)
	}

	pb.bars[displayName] = pb.progress.AddBar(int64(data.TotalEntities),
		mpb.PrependDecorators(
			decor.Name(displayName, decor.WC{W: len(displayName) + 1, C: decor.DSyncSpaceR}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)
}

func (pb *progressBar) handleBarUpdate(data BarUpdate) {
	displayName := data.BarName

	val, exists := pb.bars[displayName]
	if !exists {
		log.Panicf("trying to update a bar that doesn't exist: %s (%v)", displayName, data)
	}

	if data.Change <= 0 {
		return
	}

	val.IncrBy(data.Change)
	if val.Completed() && !pb.enabled {
		screen.Printf("Finished collecting %s\n", displayName)
	}
}

func (pb *progressBar) handleTimedBarCreation(data TimedBarCreation) {
	total := int64(time.Until(data.End).Seconds())

	displayName := data.BarName
	bar := pb.progress.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name(displayName, decor.WC{W: len(displayName) + 1, C: decor.DSyncSpaceR}),
			decor.CountersNoUnit("%ds / %ds", decor.WCSyncWidth),
		),
	)
	bar.SetPriority(math.MaxInt)

	go func() {
		for i := 0; i < int(total)-1; i++ {
			time.Sleep(time.Second)
			bar.Increment()
		}

		// must not complete to abort - so just do the last one manually
		time.Sleep(time.Second)
		bar.Abort(true)
	}()
}
func (pb *progressBar) handleBarClose(data BarClose) {
	displayName := data.BarName

	val, exists := pb.bars[displayName]
	if !exists {
		log.Panicf("trying to update a bar that doesn't exist: %s (%v)", displayName, data)
	}

	if !val.Completed() {
		log.Printf("BUG: closing bar %s although it is not completed. please report this issue to legitify repository.", displayName)
	}

	val.Abort(false)
}
