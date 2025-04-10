package models

import "fmt"

type HardwareType string

const (
	HardwareType2x2 HardwareType = "HW_2X2_V1"
	HardwareType2x7 HardwareType = "HW_2X7_V1"
)

func (ht *HardwareType) Set(value string) error {
	switch HardwareType(value) {
	case HardwareType2x2, HardwareType2x7:
		*ht = HardwareType(value)
		return nil
	default:
		return fmt.Errorf("invalid hardware type %s", value)
	}
}

func (ht *HardwareType) String() string {
	return string(*ht)
}
