package gostream

import (
	"fmt"
	"testing"
)

type IntSlice []interface{}

func (s IntSlice) Len() int { return len(s) }
func (s IntSlice) Swap(i, j int){ s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool {
	return s[i].(int) < s[j].(int)
}

func TestName(t *testing.T) {
	a := []interface{}{1,2,3,4,5}
	s:=BuildStream(a)
	//s.Map(func(i interface{}) interface{} {
	//	return i.(int)*2
	//})
	s.Limit(10)
	fmt.Println("After sorted: ", s)
}
