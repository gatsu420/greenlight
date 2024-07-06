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
