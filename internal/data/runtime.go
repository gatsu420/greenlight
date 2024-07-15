package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := strconv.Quote(fmt.Sprintf("%v mins", r))
	return []byte(jsonVal), nil
}

type Genre string

func (g Genre) MarshalJSON() ([]byte, error) {
	jsonVal := strconv.Quote(fmt.Sprintf("kategori %v", g))
	return []byte(jsonVal), nil
}

type Star string

func (s Star) MarshalJSON() ([]byte, error) {
	jv := strconv.Quote(fmt.Sprintf("pemain %v", s))
	return []byte(jv), nil
}
