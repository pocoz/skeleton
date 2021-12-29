package memlogger

import (
	"fmt"
	"strconv"
)

func BytesToHumanReadableForm(bytesNumber, reminderLength int) string {
	const (
		tb = 1099511627776
		gb = 1073741824
		mb = 1048576
		kb = 1024
	)

	var (
		valueOfAmount                   string
		currentNumber, remainder, width int
	)

	switch {
	case bytesNumber > tb:
		valueOfAmount = "tb"
		currentNumber = bytesNumber / tb
		remainder = bytesNumber - (currentNumber * tb)
	case bytesNumber > gb:
		valueOfAmount = "gb"
		currentNumber = bytesNumber / gb
		remainder = bytesNumber - (currentNumber * gb)
	case bytesNumber > mb:
		valueOfAmount = "mb"
		currentNumber = bytesNumber / mb
		remainder = bytesNumber - (currentNumber * mb)
	case bytesNumber > kb:
		valueOfAmount = "kb"
		currentNumber = bytesNumber / kb
		remainder = bytesNumber - (currentNumber * kb)
	default:
		return strconv.Itoa(bytesNumber) + " B"
	}

	if reminderLength == 0 {
		return strconv.Itoa(currentNumber) + " " + valueOfAmount
	}

	// Calculate missing leading zeroes
	switch {
	case remainder > tb:
		width = 15
	case remainder > gb:
		width = 12
	case remainder > mb:
		width = 9
	case remainder > kb:
		width = 6
	default:
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}

	if reminderLength > len(remainderString) {
		reminderLength = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", currentNumber, remainderString[:reminderLength], valueOfAmount)
}
