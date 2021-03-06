package d24

import (
	"fmt"
)

type continuousRange struct {
	min, max int
	step     int
}

var ZeroContinuousRange = continuousRange{0, 0, 1}

var OneContinuousRange = continuousRange{1, 1, 1}

type ContinuousRange interface {
	Range
	Length() int
	Min() int
	Max() int
	Step() int
}

func newContinuousRange(min, max, step int) ContinuousRange {
	if max < min {
		temp := max
		max = min
		min = temp
	}

	if min == 0 && max == 0 && step == 1 {
		return &ZeroContinuousRange
	}

	if step <= 0 {
		panic(fmt.Sprintf("Invalid step: %d", step))
	}

	if (max-min)%step != 0 {
		panic(fmt.Sprintf("not possible to get to %d from %d at step %d", max, min, step))
	}

	return &continuousRange{min, max, step}
}

func (r *continuousRange) Includes(value int) bool {
	insideBounds := value >= r.min && value <= r.max
	if !insideBounds {
		return false
	}

	// But it also needs to be on the step
	isOnStep := (value-r.min)%r.step == 0
	return isOnStep
}

func (r *continuousRange) Intersects(other *continuousRange) bool {
	tests := [][2]*continuousRange{
		{r, other},
		{other, r},
	}
	for i := range tests {
		lhs, rhs := tests[i][0], tests[i][1]

		minInside := lhs.min >= rhs.min && lhs.min <= rhs.max
		if minInside {
			return (lhs.min-rhs.min)%lhs.step == 0
		}

		maxInside := lhs.max >= rhs.min && lhs.max <= rhs.max
		if maxInside {
			return (lhs.max-rhs.min)%lhs.step == 0
		}
	}

	return false
}

func (r *continuousRange) Length() int {
	return (r.max - r.min) / r.step
}

func (r *continuousRange) Min() int {
	return r.min
}

func (r *continuousRange) Max() int {
	return r.max
}

func (r *continuousRange) Split() []ContinuousRange {
	return []ContinuousRange{r}
}

func (r *continuousRange) Step() int {
	return r.step
}

func (r *continuousRange) String() string {
	if r.min == r.max {
		return fmt.Sprintf("%d", r.min)
	} else if r.min+r.step == r.max {
		return fmt.Sprintf("<%d,%d>", r.min, r.max)
	} else if r.min+(r.step*2) == r.max {
		return fmt.Sprintf("<%d,%d,%d>", r.min, r.min+r.step, r.max)
	} else if r.step != 1 {
		return fmt.Sprintf("<%d..%d step %d>", r.min, r.max, r.step)
	} else {
		return fmt.Sprintf("<%d..%d>", r.min, r.max)
	}
}

func (r *continuousRange) Values(context string) func() (int, bool) {
	pos := 0

	return func() (int, bool) {
		if pos > r.max-r.min {
			return 0, false
		}
		value := r.min + pos
		pos += r.step
		return value, true
	}
}
