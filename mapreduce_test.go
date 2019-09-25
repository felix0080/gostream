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
	/*var a  []interface{}
	for i:=0;i<100000000 ;i++  {
		a=append(a,i)
	}*/
	a:=[]interface{}{4,3,2,1}
	s,err:=BuildStream(IntSlice(a))
	if err != nil {
		fmt.Println(err)
		return
	}
	a1:=[]interface{}{6,5,8,7}
	s1,err:=BuildStream(IntSlice(a1))
	if err != nil {
		fmt.Println(err)
		return
	}
	value:=s.Combine(s1).
		Sorted().
		Limit(6).
		Map(func(i interface{}) interface{} {
			return i.(int)+1
		}).Filter(func(i interface{}) bool {
			return i.(int)==7
		}).Reduce(func(i interface{}, i2 interface{}) interface{} {
			return i.(int)+i2.(int)
		})//==20
	//s.Limit(10)
	fmt.Println("after arr: ", s.Collect())
	fmt.Println("after value: ", value)
}
