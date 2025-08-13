package handlers

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

type CellData struct {
	CellID   int       `json:"cell_id"`
	PillName string    `json:"cell_name"`
	Time     time.Time `json:"time"`
}

type LayoutData struct {
	HardwareType models.HardwareType `json:"hardware_type"`
	Cells        []CellData          `json:"cells"`
}

func (clh *ConfirmLayoutHandler) Large(c echo.Context) error {
	layoutData := LayoutData{
		HardwareType: models.HardwareType4x7,
		// up to 28 cells for 4x7 hardware type
		Cells: []CellData{
			{CellID: 1, PillName: "Пилюля 1", Time: time.Now()},
			{CellID: 2, PillName: "Пилюля 2", Time: time.Now().Add(1 * time.Hour)},
			{CellID: 3, PillName: "Пилюля 3", Time: time.Now().Add(2 * time.Hour)},
			{CellID: 4, PillName: "Пилюля 4", Time: time.Now().Add(3 * time.Hour)},
			{CellID: 5, PillName: "Пилюля 5", Time: time.Now().Add(4 * time.Hour)},
			{CellID: 6, PillName: "Пилюля 6", Time: time.Now().Add(5 * time.Hour)},
			{CellID: 7, PillName: "Пилюля 7", Time: time.Now().Add(6 * time.Hour)},
			{CellID: 8, PillName: "Пилюля 8", Time: time.Now().Add(7 * time.Hour)},
			{CellID: 9, PillName: "Пилюля 9", Time: time.Now().Add(8 * time.Hour)},
			{CellID: 10, PillName: "Пилюля 10", Time: time.Now().Add(9 * time.Hour)},
			{CellID: 11, PillName: "Пилюля 11", Time: time.Now().Add(10 * time.Hour)},
			{CellID: 12, PillName: "Пилюля 12", Time: time.Now().Add(11 * time.Hour)},
			{CellID: 13, PillName: "Пилюля 13", Time: time.Now().Add(12 * time.Hour)},
			{CellID: 14, PillName: "Пилюля 14", Time: time.Now().Add(13 * time.Hour)},
			{CellID: 15, PillName: "Пилюля 15", Time: time.Now().Add(14 * time.Hour)},
			{CellID: 16, PillName: "Пилюля 16", Time: time.Now().Add(15 * time.Hour)},
			{CellID: 17, PillName: "Пилюля 17", Time: time.Now().Add(16 * time.Hour)},
			{CellID: 18, PillName: "Пилюля 18", Time: time.Now().Add(17 * time.Hour)},
			{CellID: 19, PillName: "Пилюля 19", Time: time.Now().Add(18 * time.Hour)},
			{CellID: 20, PillName: "Пилюля 20", Time: time.Now().Add(19 * time.Hour)},
			{CellID: 21, PillName: "Пилюля 21", Time: time.Now().Add(20 * time.Hour)},
			{CellID: 22, PillName: "Пилюля 22", Time: time.Now().Add(21 * time.Hour)},
			{CellID: 23, PillName: "Пилюля 23", Time: time.Now().Add(22 * time.Hour)},
			{CellID: 24, PillName: "Пилюля 24", Time: time.Now().Add(23 * time.Hour)},
			{CellID: 25, PillName: "Пилюля 25", Time: time.Now().Add(24 * time.Hour)},
			{CellID: 26, PillName: "Пилюля 26", Time: time.Now().Add(25 * time.Hour)},
			{CellID: 27, PillName: "Пилюля 27", Time: time.Now().Add(26 * time.Hour)},
			{CellID: 28, PillName: "Пилюля 28", Time: time.Now().Add(27 * time.Hour)},
		},
	}

	return c.JSON(200, layoutData)
}

func (clh *ConfirmLayoutHandler) Small(c echo.Context) error {
	layoutData := LayoutData{
		HardwareType: models.HardwareType2x2,
		Cells: []CellData{
			{CellID: 1, PillName: "Пилюля 1", Time: time.Now()},
			{CellID: 2, PillName: "Пилюля 2", Time: time.Now().Add(1 * time.Hour)},
			{CellID: 3, PillName: "Пилюля 3", Time: time.Now().Add(2 * time.Hour)},
			{CellID: 4, PillName: "Пилюля 4", Time: time.Now().Add(3 * time.Hour)},
		},
	}

	return c.JSON(200, layoutData)
}
