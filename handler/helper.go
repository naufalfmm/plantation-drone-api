package handler

import "sort"

func getPrevNextCoordinate(x, y, length int) (prevX, prevY, nextX, nextY int) {
	if y%2 != 0 {
		prevX = x - 1
		prevY = y

		nextX = x + 1
		nextY = y

		if x == length {
			nextX = x
			nextY = y + 1
		}

		if x == 1 {
			prevX = x
			prevY = y - 1
		}
	}

	if y%2 == 0 {
		prevX = x + 1
		prevY = y

		nextX = x - 1
		nextY = y

		if x == length {
			prevX = x
			prevY = y - 1
		}

		if x == 1 {
			nextX = x
			nextY = y + 1
		}
	}

	return
}

func findMedian(data []int) float64 {
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	if len(data)%2 > 0 {
		return float64(data[int(len(data)/2)])
	}

	mid := int(len(data) / 2)

	return float64(data[mid]+data[mid+1]) / 2
}
