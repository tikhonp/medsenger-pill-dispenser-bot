package main

import (
	"encoding/binary"
	"fmt"
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
    e.POST("/status", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]interface{}{"is_tracking_data": true, "supported_scenarios": []int{}, "tracked_contracts": []int{}})
    })

	e.Logger.Fatal(e.Start(":3054"))
}
