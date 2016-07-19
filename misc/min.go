package misc

func Min(params ...float64) float64 {

	var r float64

	r = params[0]
	for _, v := range params {
		if v < r {
			r = v
		}
	}
	return r
}
