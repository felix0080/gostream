package gostream

import (
	"reflect"
	"sort"
)

type Stream struct {
	array interface{}
}

func BuildStream(array interface{})*Stream{
	return &Stream{array}
}
//用于映射每个元素到对应的结果
func (stream *Stream) Map(f func(interface{})interface{})*Stream {
	v:=reflect.ValueOf(stream.array)
	for i := 0; i < v.Len() ; i++  {
		newItem:=f(v.Index(i).Interface())
		v.Index(i).Set(reflect.ValueOf(newItem))
	}
	return stream
}

/*
当数组小于2时，返回Nil
 */
func (stream *Stream) Reduce(f func(interface{},interface{})interface{}) interface{}{
	v:=reflect.ValueOf(stream.array)
	if v.Len() == 0 ||v.Len() == 1 {
		return nil
	}
	if v.Len() == 2 {
		return f(v.Index(0).Interface(),v.Index(1).Interface())
	}
	tmpValue:=f(v.Index(0).Interface(),v.Index(1).Interface())
	for i := 2; i < v.Len() ; i++  {
		tmpValue=f(tmpValue,v.Index(i).Interface())
	}
	return tmpValue
}

//用户对流进行排序
func (stream *Stream) Sorted() *Stream{
	arr,ok:=stream.array.(sort.Interface)
	if ok {
		sort.Sort(arr)
	}
	return stream
}
//设置条件满足的被过滤
func (stream *Stream) Filter(f func(interface{})bool)  *Stream{
	v:=reflect.ValueOf(stream.array)
	len:=v.Len()
	for i := 0; i <  len; i++  {
		if f(v.Index(i).Interface()) {
			if i == len-1 {
				v=v.Slice(0,i)
			}else{
				reflect.AppendSlice(v.Slice(0,i),v.Slice(i+1,len))

			}
			i--
			len--
		}
	}
	v=v.Slice(0,len)
	stream.array=v.Interface()
	return stream
}

//用户获取指定数量的流
func (stream *Stream) Limit(len int) *Stream{
	v:=reflect.ValueOf(stream.array)
	if len > v.Len() {
		return stream
	}
	v=v.Slice(0,len)
	stream.array=v.Interface()
	return stream
}