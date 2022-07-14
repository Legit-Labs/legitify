package severity

type Severity = string

const (
	Critical Severity = "CRITICAL"
	High     Severity = "HIGH"
	Medium   Severity = "MEDIUM"
	Low      Severity = "LOW"
	Unknown  Severity = "UNKNOWN"
)

var all = map[string]int{
	Critical: 0,
	High:     1,
	Medium:   2,
	Low:      3,
	Unknown:  4,
}

func IsValid(severity Severity) bool {
	if severity == Unknown {
		return false
	}

	_, ok := all[severity]
	return ok
}

func Less(first, second Severity) bool {
	return all[first] < all[second]
}
