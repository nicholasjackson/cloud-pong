package objects

func translateToGameX(x int32, gw int32, sw int) int {
	xRatio := float64(sw) / float64(gw)
	return int(float64(x) * xRatio)
}

func translateToGameY(y int32, gh int32, sh int) int {
	yRatio := float64(sh) / float64(gh)
	return int(float64(y) * yRatio)
}
