package pilldispenserprotocol

import (
	"encoding/binary"

	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

// Schedule data for single pill cell.
//
// So there is a defenition of a schedule data memory layout
// single struct is for single cell.
// One packet contains:
//     - Start time for cell. uint32 timestamp (in seconds)
//     - End time for cell.   uint32 timestamp (in seconds)
//     - And single byte for meta data
// SO there is 7 bytes per packet. Total data for all cells:
//
// [ 4 bytes ][ 4 bytes ][ 1 byte ]|[ 4 bytes ][ 4 bytes ][ 1 byte ]|[ 4 bytes ][ 4 bytes ][ 1 byte ]|[ 4
//                                 |                                |                                |       etc..
//   0 indx packet -> cell 1       |  1 indx packet -> cell 2       |  2 indx packet -> cell 3       |
//
// Last 4 bytes is uint32
//
// All data is little endian (esp32 is little endian)

// 4 bytes (uint32) start time timesptamp,
// 4 bytes (uint32) end time timestamp,
// 1 byte for metadata
const cellDataLength = 4 + 4 + 1

const IsOfflineNotificationsAllowedBitN = 0

// EmptySchedule generates empty data for cases when no schedule set for pill dispenser
func EmptySchedule(cellsCount int) []byte {

	// initiate data with proper length with zeros and additional 4 bytes (uint32) for refresh rate duration
	data := make([]byte, cellsCount*cellDataLength, cellsCount*cellDataLength+4)

	// add default refresh rate interval value
	data = binary.LittleEndian.AppendUint32(data, uint32(models.DefaultRefreshRateInterval))

	return data
}

// ScheduleFromScheduleData encodes schedule data struct to byte array
func ScheduleFromScheduleData(s *models.ScheduleData) []byte {

	cellsCount := len(s.Cells)

	// (uint32, uint32, uint8) * cell-count + uint32
	data := make([]byte, 0, cellsCount*cellDataLength+4)

	for _, cell := range s.Cells {

		timestampStart := uint32(cell.StartTime.Time.UTC().Unix())
		data = binary.LittleEndian.AppendUint32(data, timestampStart)

		timestampEnd := uint32(cell.EndTime.Time.UTC().Unix())
		data = binary.LittleEndian.AppendUint32(data, timestampEnd)

		var meta uint8
		if s.Schedule.IsOfflineNotificationsAllowed {
			meta = meta | (uint8(1) << IsOfflineNotificationsAllowedBitN)
		}
		data = append(data, meta)
	}

	data = binary.LittleEndian.AppendUint32(data, uint32(s.Schedule.RefreshRateInterval.Int64))

	return data
}
