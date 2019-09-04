package objects

func translateToGameX(x int32, gw int32, sw int) int {
	xRatio := float64(sw) / float64(gw)
	xScale := float64(x) * xRatio

	if xScale > 0 && xScale < 1 {
		xScale = 1
	}
	return int(xScale)
}

func translateToGameY(y int32, gh int32, sh int) int {
	yRatio := float64(sh) / float64(gh)
	yScale := float64(y) * yRatio
	if yScale > 0 && yScale < 1 {
		yScale = 1
	}
	return int(yScale)
}
