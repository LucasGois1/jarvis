// Code generated by "stringer -type=Status"; DO NOT EDIT.

package entities

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Ended-1]
	_ = x[Active-2]
}

const (
	_Status_name_0 = "Ended"
	_Status_name_1 = "Active"
)

func (i Status) String() string {
	switch {
	case i == 1:
		return _Status_name_0
	case i == 2:
		return _Status_name_1
	default:
		return "Status(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
