package progressbar

import "time"

type ChannelType = interface{}

type MinimumRequiredBars struct {
	count int
}

type RequiredBarCreation OptionalBarCreation
type OptionalBarCreation struct {
	BarName       string
	TotalEntities int
}

type OptionalSpinnerBarCreation struct {
	BarName string
}

type BarUpdate struct {
	BarName string
	Change  int
}

type SpinnerBarUpdate struct {
	BarUpdate
}

type TimedBarCreation struct {
	BarName string
	End     time.Time
}

type BarClose struct {
	BarName          string
	AllowUncompleted bool
}

// NewMinimumRequiredBars creates a request to set the minimum number of bars.
// It is used to prevent the progress bar from finishing before all bars were created.
func NewMinimumRequiredBars(count int) MinimumRequiredBars {
	return MinimumRequiredBars{
		count: count,
	}
}

// NewRequiredBar creates a request to create a new required bar.
// It is used to create a bar and mark it for the minimum requirement count.
func NewRequiredBar(name string, total int) RequiredBarCreation {
	return RequiredBarCreation{
		BarName:       name,
		TotalEntities: total,
	}
}

// NewOptionalBar creates a request to create a new optional bar.
// It is used to create a bar without marking it for the minimum requirement count.
func NewOptionalBar(name string, total int) OptionalBarCreation {
	return OptionalBarCreation{
		BarName:       name,
		TotalEntities: total,
	}
}

func NewOptionalSpinnerBar(name string) OptionalSpinnerBarCreation {
	return OptionalSpinnerBarCreation{
		BarName: name,
	}
}

// NewUpdate creates a request to update the count for an existing bar.
// The change must be positive.
func NewUpdate(name string, change int) BarUpdate {
	return BarUpdate{
		BarName: name,
		Change:  change,
	}
}

func NewSpinnerBarUpdate(name string, change int) SpinnerBarUpdate {
	return SpinnerBarUpdate{
		BarUpdate: BarUpdate{
			BarName: name,
			Change:  change,
		},
	}
}

// NewTimedBar creates a request to create a time-based bar.
// Timed bars are counted with seconds and removed when finished.
func NewTimedBar(name string, end time.Time) TimedBarCreation {
	return TimedBarCreation{
		BarName: name,
		End:     end,
	}
}

// NewBarClose creates a request to close an existing bar.
// It is used to prevent the program from being stuck if a progress bar is not completed (due to error).
// In case the progress bar already completed, it is just ignored.
func NewBarClose(name string) BarClose {
	return BarClose{
		BarName:          name,
		AllowUncompleted: false,
	}
}

func NewBarCloseAllowUncompleted(name string) BarClose {
	return BarClose{
		BarName:          name,
		AllowUncompleted: true,
	}
}
