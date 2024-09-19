package tools

import "fmt"

const (
	DefaultPageLimit  = 10
	DefaultPageNumber = 0
)

func Paginate(limit, num, length int) (startIndex, endIndex int, err error) {
	if limit < 0 || num < 0 {
		return 0, 0, fmt.Errorf("invalid params, limit, num: %d %d", limit, num)
	}
	startIndex = num * limit
	if startIndex >= length {
		startIndex = length
	}
	endIndex = startIndex + limit
	if endIndex > length {
		endIndex = length
	}

	return startIndex, endIndex, nil
}
