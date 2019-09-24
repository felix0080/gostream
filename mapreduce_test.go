package gostream

import (
	"testing"
)

type IntSlice []interface{}

func (s IntSlice) Len() int { return len(s) }
func (s IntSlice) Swap(i, j int){ s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool {
	return s[i].(int) < s[j].(int)
}

func TestName(t *testing.T) {
	var a  []interface{}
	for i:=0;i<100000000 ;i++  {
		a=append(a,i)
	}
	s:=BuildStream(a)
	s.Map(func(i interface{}) interface{} {
		return i.(int)*2
	})
	//s.Limit(10)
	//fmt.Println("After sorted: ", s)
}
