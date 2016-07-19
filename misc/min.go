package misc

func Min(params ...int) int {

	var r int

	r = params[0]
	for _, v := range params {
		if v < r {
			r = v
		}
	}
	return r
}
