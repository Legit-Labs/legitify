package progressbar

import "time"

type ChannelType = interface{}

type MinimalBars struct {
	count int
}

type RequiredBarCreation OptionalBarCreation
type OptionalBarCreation struct {
	BarName       string
	TotalEntities int
}

type BarUpdate struct {
	BarName string
	Change  int
}

type TimedBarCreation struct {
	BarName string
	End     time.Time
}

func NewMinimalBars(count int) MinimalBars {
	return MinimalBars{
		count: count,
	}
}

func NewRequiredBar(name string, total int) RequiredBarCreation {
	return RequiredBarCreation{
		BarName:       name,
		TotalEntities: total,
	}
}

func NewOptionalBar(name string, total int) OptionalBarCreation {
	return OptionalBarCreation{
		BarName:       name,
		TotalEntities: total,
	}
}

func NewUpdate(name string, change int) BarUpdate {
	return BarUpdate{
		BarName: name,
		Change:  change,
	}
}

func NewTimedBar(name string, end time.Time) TimedBarCreation {
	return TimedBarCreation{
		BarName: name,
		End:     end,
	}
}
