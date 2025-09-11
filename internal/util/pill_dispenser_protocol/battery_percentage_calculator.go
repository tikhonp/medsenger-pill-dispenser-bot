package pilldispenserprotocol

const (
	meanCutoff = 3303.4193548387098
	meanDenom  = 370.6629032258064
)

func clip(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// BatteryPercentageFromVoltage calculates battery percentage from voltage in millivolts.
// percent = (value - mean_cutoff) / mean_denom * 100
// If denom is zero or NaN returns np.nan.
func BatteryPercentageFromVoltage(voltageMillivolts int) int {
	pct := (float64(voltageMillivolts) - meanCutoff) / meanDenom * 100
	pct = clip(pct, 0, 100)
	return int(pct)
}

func IsBatteryLow(voltageMillivolts int) bool {
	return BatteryPercentageFromVoltage(voltageMillivolts) < 20
}
