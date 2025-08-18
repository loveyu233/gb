package gb

import "github.com/go-resty/resty/v2"

func R() *resty.Request {
	return resty.New().R()
}
