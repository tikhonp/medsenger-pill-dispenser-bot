package models

import "fmt"

type HardwareType string

const (
	HardwareType2x2 HardwareType = "HW_2X2_V1"
	HardwareType4x7 HardwareType = "HW_4X7_V1"
)

func (ht *HardwareType) Set(value string) error {
	switch HardwareType(value) {
	case HardwareType2x2, HardwareType4x7:
		*ht = HardwareType(value)
		return nil
	default:
		return fmt.Errorf("invalid hardware type %s", value)
	}
}

func (ht *HardwareType) String() string {
	return string(*ht)
}

func (ht HardwareType) GetCellsCount() int {
	switch ht {
	case HardwareType2x2:
		return 4
	case HardwareType4x7:
		return 28
	}
	return 0
}
