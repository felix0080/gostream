package gostream

import (
	"fmt"
	"github.com/seiflotfy/cuckoofilter"
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
func BuildStream(array interface{}) (*Stream, error) {
	/*
		需要检查array 是否是slice
	*/
	if reflect.ValueOf(array).Kind() != reflect.Slice {
		return nil, fmt.Errorf("%s", "must be a slice")
	}
	return &Stream{array}, nil
}

//流拷贝，拷贝一个新流，返回与以前的流相同大小，容量的新流，且将内容拷贝进去
//	s2:=s.Copy()
//	s2.Map(func(i interface{}) interface{} {
//		return i.(int)+1
//	}).Sorted()
//	fmt.Println(s2.Collect())
//	fmt.Println(s.Collect())
func (stream *Stream) Copy() *Stream {
	value := reflect.ValueOf(stream.array)
	a2 := reflect.MakeSlice(value.Type(), value.Len(), value.Cap())
	reflect.Copy(a2, value)
	return &Stream{
		array: a2.Interface(),
	}
}

//用于映射每个元素到对应的结果
//返回的数组的类型需要指定，传入第一个参数reflect.Type
func (stream *Stream) Map(types interface{},f func(interface{}) interface{}) *Stream {
	v := reflect.ValueOf(stream.array)
	lens := v.Len()
	nslice := reflect.MakeSlice(reflect.TypeOf(types), lens, lens)
	for i := 0; i < lens; i++ {
		newItem := f(v.Index(i).Interface())
		//newItem can't be nil
		if newItem != nil {
			//新建一个新流存储数据，防止类型错误
			nslice.Index(i).Set(reflect.ValueOf(newItem))
		}
	}
	stream.array = nslice.Interface()
	return stream
}

/*func (stream *Stream) Copy()*Stream {
	a1:=[]interface{}{6,5,"asd",7}
	a2:=make([]interface{},len(a1))
	fmt.Println(fmt.Sprintf("%p %p",a1,a2))
	copy(a2,a1)
	fmt.Println(a2)
}*/
//用于映射每个元素到对应的结果
//返回的数组的类型需要指定，传入第一个参数reflect.Type
func (stream *Stream) MultipartMap(worknum int,types interface{}, f func(interface{}) interface{}) *Stream {
	v := reflect.ValueOf(stream.array)
	lens := v.Len()
	nslice:=reflect.MakeSlice(reflect.TypeOf(types),lens,lens)
	workitemnum := lens / worknum
	lastnum := lens % worknum
	var wg sync.WaitGroup
	for i := 0; i <= worknum; i++ {
		wg.Add(1)
		start := i * workitemnum
		if i != worknum {
			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					newItem := f(v.Index(i).Interface())
					nslice.Index(i).Set(reflect.ValueOf(newItem))
				}
			}(start, start+workitemnum)
			continue
		}
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				newItem := f(v.Index(i).Interface())
				nslice.Index(i).Set(reflect.ValueOf(newItem))
			}
		}(start, start+lastnum)
	}
	wg.Wait()
	stream.array = nslice.Interface()
	return stream
}

/*
当数组小于2时，返回Nil
*/
func (stream *Stream) Reduce(f func(interface{}, interface{}) interface{}) interface{} {
	v := reflect.ValueOf(stream.array)
	if v.Len() == 0 || v.Len() == 1 {
		return nil
	}
	if v.Len() == 2 {
		return f(v.Index(0).Interface(), v.Index(1).Interface())
	}
	tmpValue := f(v.Index(0).Interface(), v.Index(1).Interface())
	for i := 2; i < v.Len(); i++ {
		tmpValue = f(tmpValue, v.Index(i).Interface())
	}
	return tmpValue
}

//用户对流进行排序
//需要自行实现sort接口才可使用此方法
func (stream *Stream) Sorted() *Stream {
	arr, ok := stream.array.(sort.Interface)
	if ok {
		sort.Sort(arr)
	} else {
		panic("The sort interface is not implemented")
	}
	return stream
}

//设置条件满足的被过滤
func (stream *Stream) Filter(f func(interface{}) bool) *Stream {
	v := reflect.ValueOf(stream.array)
	len := v.Len()
	for i := 0; i < len; i++ {
		if f(v.Index(i).Interface()) {
			if i == len-1 {
				v = v.Slice(0, i)
			} else {
				reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, len))
			}
			i--
			len--
		}
	}
	v = v.Slice(0, len)
	stream.array = v.Interface()
	return stream
}

/*
  组合流，将两个stream 简单组合在一起
*/
func (stream *Stream) Combine(anotherStream *Stream) *Stream {
	v := reflect.ValueOf(stream.array)
	anov := reflect.ValueOf(anotherStream.array)
	value := reflect.AppendSlice(v.Slice(0, v.Len()), anov.Slice(0, anov.Len()))
	stream.array = value.Interface()
	return stream
}

/*
	收集流，将流的内部数组返回。
*/
func (stream *Stream) Collect() interface{} {
	return stream.array
}

//用户获取指定数量的流
func (stream *Stream) Limit(len int) *Stream {
	v := reflect.ValueOf(stream.array)
	if len > v.Len() {
		return stream
	}
	v = v.Slice(0, len)
	stream.array = v.Interface()
	return stream
}

//此方法将需要流中的元素实现 unique接口 hashCode（）
//算法采用布谷鸟算法
//	Example
//	func (s Item) HashCode() []byte {
//		return []byte(strconv.Itoa(int(s)))
//	}
//	func (s IntSlice) Len() int      { return len(s) }
//	func (s IntSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
//	func (s IntSlice) Less(i, j int) bool {
//		return s[i].(int) < s[j].(int)
//	}
//	func TestStream_Distinct(t *testing.T) {
//		a := []interface{}{Item(1), Item(2), Item(4), Item(4)}
//		s, err := BuildStream(IntSlice(a))
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		s.Distinct()
//		fmt.Println(s.Collect())
//	}
func (stream *Stream) Distinct() *Stream {
	v := reflect.ValueOf(stream.array)
	lens := v.Len()
	cf := cuckoo.NewFilter(1000)
	for i := 0; i < lens; i++ {
		newItem, ok := v.Index(i).Interface().(Unique)
		if ok && newItem != nil {
			value := newItem.HashCode()
			//存在，将删除此元素
			if !cf.InsertUnique(value) {
				if i == lens-1 {
					v = v.Slice(0, i)
				} else {
					reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, lens))
				}
				i--
				lens--
			}
		}
	}
	v = v.Slice(0, lens)
	stream.array = v.Interface()
	cf.Reset()
	return stream
}

type Unique interface {
	HashCode() []byte
}
