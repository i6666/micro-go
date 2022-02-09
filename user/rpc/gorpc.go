package rpc

import "errors"

type StringRequest struct {
	A string
	B string
}

var ErrMaxSize = errors.New("err max size")

type Service interface {
	// 拼接 A + B
	Concat(req StringRequest, ret *string) error
}
type StringService struct {
}

func (s StringService) Concat(req StringRequest, ret *string) error {
	if len(req.A)+len(req.B) > 10 {
		*ret = ""
		return ErrMaxSize
	}

	*ret = req.A + req.B
	return nil
}
