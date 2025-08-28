// Package handlers provides HTTP handlers for the Medsenger Pill Dispenser Bot service.
package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/diagnostics/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (diagn *ProvideDiagnosticsHandler) Get(c echo.Context) error {
	statuses, err := diagn.DB.BatteryStatuses().GetAll()
	if err != nil {
		return err
	}

	// Prepare slices expected by the template
	var voltageData [][]float64
	var timeLabels [][]string
	var seriesNames []string

	if len(statuses) == 0 {
		return util.TemplRender(c, views.ChargePage(voltageData, timeLabels, seriesNames, ""))
	}

	// statuses ordered by serial_nu, created_at ASC by GetAll
	currSerial := statuses[0].SerialNumber
	var currVoltages []float64
	var currTimes []string
	// var currIDs []string

	for _, s := range statuses {
		if s.SerialNumber != currSerial {
			// flush current series
			voltageData = append(voltageData, currVoltages)
			timeLabels = append(timeLabels, currTimes)
			seriesNames = append(seriesNames, currSerial)

			// reset for next series
			currSerial = s.SerialNumber
			currVoltages = nil
			currTimes = nil
			// currIDs = nil
		}

		currVoltages = append(currVoltages, float64(s.Voltage))
		currTimes = append(currTimes, s.CreatedAt.Format("2006-01-02 15:04:05"))
		// currIDs = append(currIDs, fmt.Sprintf("%d", s.ID))
	}

	// flush last series
	if len(currVoltages) > 0 || len(currTimes) > 0 {
		voltageData = append(voltageData, currVoltages)
		timeLabels = append(timeLabels, currTimes)
		seriesNames = append(seriesNames, currSerial)
	}

	return util.TemplRender(c, views.ChargePage(voltageData, timeLabels, seriesNames, ""))
}

