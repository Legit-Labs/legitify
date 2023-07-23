package errlog

import (
	"encoding/json"
	"fmt"
)

type reasonType int

const (
	reasonTypePrerequisite reasonType = iota
	reasonTypePermission   reasonType = iota
)

type SkipReason struct {
	reason     string
	reasonType reasonType
}

func NewPrerequisiteSkipReason(reason string) SkipReason {
	return SkipReason{
		reason:     reason,
		reasonType: reasonTypePrerequisite,
	}
}

func NewPermissionSkipReason(reason string) SkipReason {
	return SkipReason{
		reason:     reason,
		reasonType: reasonTypePermission,
	}
}

func (s SkipReason) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s: %s", s.ReasonPrefix(), s.reason))
}

func (s SkipReason) ReasonPrefix() string {
	switch s.reasonType {
	case reasonTypePrerequisite:
		return "Unmet prerequisite"
	case reasonTypePermission:
		return "Missing permission"
	default:
		return fmt.Sprintf("Unknown reason type: %d", s.reasonType)
	}
}
