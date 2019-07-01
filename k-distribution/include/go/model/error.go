// File provided by the K Framework Go backend. Timestamp: 2019-06-30 21:44:04.091

package impmodel

import "fmt"

type parseIntError struct {
	parseVal string
}

func (e *parseIntError) Error() string {
	return fmt.Sprintf("Could not parse int from value: %s", e.parseVal)
}
