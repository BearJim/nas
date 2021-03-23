package nasConvert

import (
	"errors"
)

// TS 24.008 10.5.7.4, TS 24.501 9.11.2.4
// the unit of timerValue is second
func GPRSTimer2ToNas(timerValue int) (uint8, error) {
	timerValueNas := uint8(0)

	if timerValue <= 64 {
		if timerValue%2 != 0 {
			return 0, errors.New("timer Value is not multiples of 2 seconds")
		}
		timerValueNas = uint8(timerValue / 2)
	} else {
		t := uint8(timerValue / 60) // t is multiples of 1 min
		if t <= 31 {
			timerValueNas = (timerValueNas | 0x20) + t
		} else {
			if t%6 != 0 {
				return 0, errors.New("timer Value is not multiples of decihours")
			}
			t = t / 6
			timerValueNas = (timerValueNas | 0x40) + t
		}
	}
	return timerValueNas, nil
}
