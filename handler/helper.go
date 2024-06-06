package handler

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
