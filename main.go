package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const ContentTypeOctetStream = "application/octet-stream"

func pillsSchedulerHandler(c echo.Context) error {
	serialNumber := c.QueryParam("serial_number")
	if serialNumber != "" {
		fmt.Println("Scheduler request from:", serialNumber)
	}

	var cellsCount = 4
	cellsCountStr := c.QueryParam("cells_count")
	if cellsCountP, err := strconv.Atoi(cellsCountStr); err == nil {
		cellsCount = cellsCountP
	}

	out := make([]byte, 0, cellsCount*4)
	now := time.Now()
	for i := range cellsCount {
		// make array of timestamps from now and every minute after
		timestamp := uint32(now.Add(time.Minute * time.Duration(i)).Unix())
		out = binary.BigEndian.AppendUint32(out, timestamp)
	}

	return c.Blob(http.StatusOK, ContentTypeOctetStream, out)
}

func pillsSchedulerHandlerV2(c echo.Context) error {
	serialNumber := c.QueryParam("serial_number")
	if serialNumber != "" {
		fmt.Println("Scheduler request from:", serialNumber)
	}

	var cellsCount = 4
	cellsCountStr := c.QueryParam("cells_count")
	if cellsCountP, err := strconv.Atoi(cellsCountStr); err == nil {
		cellsCount = cellsCountP
	}

	out := make([]byte, 0, cellsCount*(4+4+1))
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

	return c.Blob(http.StatusOK, ContentTypeOctetStream, out)
}

func submitPill(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if len(body) < 5 {
		fmt.Printf("SUBMIT: recieved body less than 5 bytes (%d)", len(body))
		return echo.NewHTTPError(http.StatusBadRequest, "data length must be at least 5 bytes")
	}

	var timestamp uint32
	_, err = binary.Decode(body[:4], binary.BigEndian, &timestamp)
	if err != nil {
		fmt.Printf("SUBMIT: err: %s", err.Error())
		return err
	}
	submitTime := time.Unix(int64(timestamp), 0)
	cellIndx := body[4]
	serialNumber := string(body[5:])

	fmt.Printf("Submit pill request: [time > %v cell > %d serialn > %s]\n", submitTime, cellIndx, serialNumber)

	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Купил мужик шляпу, а она ему как раз!")
	})
	e.GET("/schedule", pillsSchedulerHandler)
	e.GET("/schedule/v2", pillsSchedulerHandlerV2)
	e.POST("/submit", submitPill)
	e.POST("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"is_tracking_data": true, "supported_scenarios": []int{}, "tracked_contracts": []int{}})
	})

	e.Logger.Fatal(e.Start(":3054"))
}
