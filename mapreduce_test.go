package gostream

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
)

type IntSlice []interface{}
type IntSlice2 []interface{}
type Item int

func (s Item) HashCode() []byte {
	return []byte(strconv.Itoa(int(s)))
}
func (s IntSlice) Len() int      { return len(s) }
func (s IntSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool {
	return s[i].(int) < s[j].(int)
}
func (s IntSlice2) Len() int      { return len(s) }
func (s IntSlice2) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s IntSlice2) Less(i, j int) bool {
	return s[i].(int) > s[j].(int)
}
func TestStream_Distinct(t *testing.T) {
	a := []interface{}{Item(1), Item(2), Item(4), Item(4)}
	s:= BuildStream(IntSlice(a))
	s.Distinct()
	fmt.Println(s.Collect())
}
func TestName(t *testing.T) {
	/*var a  []interface{}
	for i:=0;i<100000000 ;i++  {
		a=append(a,i)
	}*/
	a := []interface{}{4, 3, 2, 1}
	s := BuildStream(IntSlice(a))
	a1 := []interface{}{6, 5, 8, 7}
	s1 := BuildStream(IntSlice(a1))
	var types []int
	s.Combine(s1).
		Sorted().
		Limit(6).
		Map(types,func(i interface{}) interface{} {
			return 1
		}) //==20
	//s.Limit(10)
	fmt.Println("after arr: ", s.Collect())
	//fmt.Println("after value: ", value)
}
func TestCopy(t *testing.T) {
	a := []interface{}{4, 3, 2, 1}
	s := BuildStream(IntSlice(a))
	s1 := s.Copy()
	var types []int
	s1.Map(types,func(i interface{}) interface{} {
		return i.(int) + 1
	}).Sorted()
	fmt.Println(s1.Collect())
	fmt.Println(s.Collect())
}

type ForFoo struct {
	Item int
}
func TestStream_Map(t *testing.T) {
	a := []ForFoo{ForFoo{4}, ForFoo{3}, ForFoo{2}, ForFoo{1}}
	s := BuildStream((a))
	//s1 := s.Copy()
	var types []int
	s.Map(types,func(i interface{}) interface{} {
		return i.(ForFoo).Item + 1
	}).Sorted()
	//fmt.Println(s1.Collect())
	//_,ok:=s.Collect().([]interface{})
	log.Println(reflect.ValueOf(s.Collect()).Type())
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
		a1 := []interface{}{6, 5, 8, 7}
		a2 := make([]interface{}, len(a1))
		//fmt.Println(fmt.Sprintf("%p %p",a1,a2))
		Clone(&a1, &a2)
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
		a1 := []interface{}{6, 5, "asd", 7}
		move := a1
		//a2:=make([]interface{},len(a1))
		value := reflect.ValueOf(move)
		a2 := reflect.MakeSlice(value.Type(), value.Len(), value.Cap())
		reflect.Copy(a2, reflect.ValueOf(move))
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
func TestStream_GroupBy(t *testing.T) {
/*	var i []int
	v := reflect.MakeSlice(reflect.TypeOf(i), 5, 5)

	v=reflect.Append(v, reflect.ValueOf(1))
	fmt.Println(v.Index(1))
	fmt.Println(v.Len())*/
	/*s := Student{
		Name:"fanxing",
	}
	v:=reflect.ValueOf(s)
	v1:=v.FieldByName("Name")
	fmt.Println(v1.Kind()==reflect.String)*/
	var s []Student
	for i := 0; i < 5; i++ {
		s=append(s,Student{strconv.Itoa(i)})
	}
	s=append(s, Student{strconv.Itoa(1)})
	ss:=BuildStream(s)
	sss:=ss.GroupBy("Name")
	fmt.Println(sss)
}

type Student struct {
	Name string
}