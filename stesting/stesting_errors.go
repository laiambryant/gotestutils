package stesting

import "fmt"

type StressTestingError struct {
	Index uint32
	Err   error
}

func (s StressTestingError) Error() string {
	return "Error while running stress test at step " + fmt.Sprint(s.Index) + " of testing: " + s.Err.Error()
}
