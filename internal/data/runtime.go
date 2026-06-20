package data

import (
	"fmt"
	"strconv"
)

type Runtime int

// We're deliberately using a `value rereiver` for the method.
func (r Runtime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(fmt.Sprintf("% mins", r))), nil
}
