package gostream

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)
/*
定义一个流，内部为 实现了Sort 接口的 Slice(interfalce{}).
若没有实现Sort接口，sort方法将不可使用
使用方法类似于java8 stream
欢迎贡献，大家一起完善此库
 */
type Stream struct {
	array interface{}
}

// Example:
//
//	type IntSlice []interface{}
//
//	func (s IntSlice) Len() int { return len(s) }
//	func (s IntSlice) Swap(i, j int){ s[i], s[j] = s[j], s[i] }
//	func (s IntSlice) Less(i, j int) bool {
//		return s[i].(int) < s[j].(int)
//	}
//	a:=[]interface{}{4,3,2,1}
//	s,err:=BuildStream(IntSlice(a))
//
func BuildStream(array interface{})(*Stream,error){
	/*
	需要检查array 是否是slice
	 */
	if reflect.ValueOf(array).Kind() != reflect.Slice {
		return nil,fmt.Errorf("%s","must be a slice")
	}
	return &Stream{array},nil
}
//用于映射每个元素到对应的结果
func (stream *Stream) Map(f func(interface{})interface{})*Stream {
	v:=reflect.ValueOf(stream.array)
	lens:=v.Len()
	for i := 0; i < lens ; i++  {
		newItem:=f(v.Index(i).Interface())
		//newItem can't be nil
		if newItem !=nil {
			v.Index(i).Set(reflect.ValueOf(newItem))
		}
	}
	return stream
}
//用于映射每个元素到对应的结果
func (stream *Stream) MultipartMap(worknum int,f func(interface{})interface{})*Stream {
	v:=reflect.ValueOf(stream.array)
	lens:=v.Len()
	workitemnum:=lens / worknum
	lastnum:=lens % worknum
	var wg sync.WaitGroup
	for i:=0;i<=worknum ;i++  {
		wg.Add(1)
		start := i * workitemnum
		if i != worknum {
			go func(start,end int) {
				defer wg.Done()
				for i := start; i < end; i++  {
					newItem:=f(v.Index(i).Interface())
					v.Index(i).Set(reflect.ValueOf(newItem))
				}
			}( start , start + workitemnum )
			continue
		}
		go func(start,end int) {
			defer wg.Done()
			for i := start; i < end; i++  {
				newItem:=f(v.Index(i).Interface())
				v.Index(i).Set(reflect.ValueOf(newItem))
			}
		}( start , start + lastnum )
	}
	wg.Wait()
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
//需要自行实现sort接口才可使用此方法
func (stream *Stream) Sorted() *Stream{
	arr,ok:=stream.array.(sort.Interface)
	if ok {
		sort.Sort(arr)
	}else {
		fmt.Println("没实现该接口")
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
/*
  组合流，将两个stream 简单组合在一起
 */
func (stream *Stream) Combine(anotherStream *Stream) *Stream {
	v:=reflect.ValueOf(stream.array)
	anov:=reflect.ValueOf(anotherStream.array)
	value:=reflect.AppendSlice(v.Slice(0,v.Len()),anov.Slice(0,anov.Len()))
	stream.array=value.Interface()
	return stream
}
/*
	收集流，将流的内部数组返回。
 */
func (stream *Stream) Collect() interface{} {
	return stream.array
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