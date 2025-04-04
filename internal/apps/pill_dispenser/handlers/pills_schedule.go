package handlers

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const ContentTypeOctetStream = "application/octet-stream"

func (pdh *PillDispenserHandler) GetSchedule(c echo.Context) error {
	serialNumber := c.QueryParam("serial_number")
	if serialNumber != "" {
		fmt.Println("Scheduler request from:", serialNumber)
	}

	var cellsCount = 4
	cellsCountStr := c.QueryParam("cells_count")
	if cellsCountP, err := strconv.Atoi(cellsCountStr); err == nil {
		cellsCount = cellsCountP
	}

	// (uint32, uint32, uint8) * cellcount  + uint32
	out := make([]byte, 0, cellsCount*(4+4+1)+4)
	now := time.Now()
	for i := range cellsCount {
		// make array of timestamps from now and every minute after
		timestamp_start := uint32(now.Add(time.Minute * time.Duration(i)).Unix())
		timestamp_end := uint32(now.Add(time.Minute * time.Duration(i)).Add(time.Second * 30).Unix())
		meta := uint8(1)
		out = binary.LittleEndian.AppendUint32(out, timestamp_start)
		out = binary.LittleEndian.AppendUint32(out, timestamp_end)
		out = append(out, meta)
	}

	var refreshRateSec uint32 = 60 * 2
	out = binary.LittleEndian.AppendUint32(out, refreshRateSec)

	return c.Blob(http.StatusOK, ContentTypeOctetStream, out)
}
