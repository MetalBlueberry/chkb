// Code generated by "stringer -type=KeyActions -trimprefix KeyAction"; DO NOT EDIT.

package chkb

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KeyActionNil - -1]
	_ = x[KeyActionMap-0]
	_ = x[KeyActionDown-1]
	_ = x[KeyActionUp-2]
	_ = x[KeyActionTap-3]
	_ = x[KeyActionDoubleTap-4]
	_ = x[KeyActionHold-5]
}

const _KeyActions_name = "NilMapDownUpTapDoubleTapHold"

var _KeyActions_index = [...]uint8{0, 3, 6, 10, 12, 15, 24, 28}

func (i KeyActions) String() string {
	i -= -1
	if i < 0 || i >= KeyActions(len(_KeyActions_index)-1) {
		return "KeyActions(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _KeyActions_name[_KeyActions_index[i]:_KeyActions_index[i+1]]
}