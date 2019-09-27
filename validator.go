package tasker

import "fmt"

type validator struct{}

func (v *validator) validateHour(value int64) {
	if value < 24 && value >= 0 {
		return
	}

	panic(fmt.Errorf("Hour value must be between 0 and 23 inclusively"))
}

func (v *validator) validateMinute(value int64) {
	if value < 60 && value >= 0 {
		return
	}

	panic(fmt.Errorf("Minute value must be between 0 and 59 inclusively"))
}
