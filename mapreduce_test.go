package gostream

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
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
	s.Combine(s1).
		Sorted().
		Limit(6).
		Map(func(i interface{}) interface{} {
			return nil
		})//==20
	//s.Limit(10)
	fmt.Println("after arr: ", s.Collect())
	//fmt.Println("after value: ", value)
}
func TestCopy(t *testing.T) {
	a:=[]interface{}{4,3,2,1}
	s,err:=BuildStream(IntSlice(a))
	if err != nil {
		fmt.Println(err)
		return
	}
	s1:=s.Copy()
	s1.Map(func(i interface{}) interface{} {
		return i.(int)+1
	}).Sorted()
	fmt.Println(s1.Collect())
	fmt.Println(s.Collect())
}

func BenchmarkCopyRegule(b *testing.B) {
	/*
	goos: windows
	goarch: amd64
	pkg: github.com/felix0080/gostream
	BenchmarkCopy-4   	   33630	     31917 ns/op
	PASS
	*/
	for i := 0; i < b.N; i++ {
		a1:=[]interface{}{6,5,8,7}
		a2:=make([]interface{},len(a1))
		//fmt.Println(fmt.Sprintf("%p %p",a1,a2))
		Clone(&a1,&a2)
		//fmt.Println(a2)
	}
}
func BenchmarkCopyReflect(b *testing.B) {
	/*
		使用反射前
		goos: windows
		goarch: amd64
		pkg: github.com/felix0080/gostream
		BenchmarkCopy1-4   	10915362	       101 ns/op
		PASS
		使用反射后
		goos: windows
		goarch: amd64
		pkg: github.com/felix0080/gostream
		BenchmarkCopy1-4   	 2914296	       436 ns/op
		PASS

		总结有300ns的损耗，还不错，采纳反射
	*/
	for i := 0; i < b.N; i++ {
		a1:=[]interface{}{6,5,"asd",7}
		move:=a1
		//a2:=make([]interface{},len(a1))
		value:=reflect.ValueOf(move)
		a2:=reflect.MakeSlice(value.Type(),value.Len(),value.Cap())
		reflect.Copy(a2,reflect.ValueOf(move))
	}
}
// Clone 完整复制数据
func Clone(a, b interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(a); err != nil {
		return err
	}
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}