package misc

func Min(params ...int64) int64 {

	r := params[0]
	for _, v := range params {
		if v < r {
			r = v
		}
	}
	return r
}
