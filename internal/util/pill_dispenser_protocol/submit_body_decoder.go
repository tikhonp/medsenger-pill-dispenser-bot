// Package pilldispenserprotocol provides functionality to decode the body of a submit message from a pill dispenser.
package pilldispenserprotocol

import (
	"encoding/binary"
	"fmt"
	"time"
)

type PillDispenserSubmitBody []byte

// Decode decodes data sent from submit from pill dispenser
func (data PillDispenserSubmitBody) Decode() (submitTime time.Time, cellIndex int, serialNumber string, err error) {
	if len(data) < 5 {
		return submitTime, cellIndex, serialNumber, fmt.Errorf("SUBMIT: recieved body less than 5 bytes (%d)", len(data))
	}

	var timestamp uint32
	_, err = binary.Decode(data[:4], binary.BigEndian, &timestamp)
	if err != nil {
		return submitTime, cellIndex, serialNumber, err
	}

	submitTime = time.Unix(int64(timestamp), 0)
	cellIndex = int(data[4])
	serialNumber = string(data[5:])

	return submitTime, cellIndex, serialNumber, err
}
