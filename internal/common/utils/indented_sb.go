package utils

import (
	"fmt"
	"strings"
)

type PrependedStringBuilder struct {
	prepend string
	sb      strings.Builder
}

func NewPrependedStringBuilder(prepend string) *PrependedStringBuilder {
	return &PrependedStringBuilder{
		prepend: prepend,
	}
}

func (isb *PrependedStringBuilder) WriteStringf(format string, args ...interface{}) {
	isb.WriteString(fmt.Sprintf(format, args...))
}

func (isb *PrependedStringBuilder) WriteString(str string) {
	isb.sb.WriteString(isb.prepend)
	isb.sb.WriteString(str)
}

func (isb *PrependedStringBuilder) String() string {
	return isb.sb.String()
}
