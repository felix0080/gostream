/*
Package gostream implements a stream slice.

Usage

调用BuildStream 方法，传入实现了sort 接口的slice，若不需要排序，则无需实现sort
	//定义一个slice.填充元素
	a:=[]interface{}{4,3,2,1}
	//使用实现了sort的type 包装，并调用调用BuildStream
	s:=BuildStream(IntSlice(a))
	//定义一个slice.填充元素
	a1:=[]interface{}{6,5,8,7}
	//使用实现了sort的type 包装,并调用调用BuildStream
	s1:=BuildStream(IntSlice(a1))
	//将两个流合并为1个，并且以s为主流
	value:=s.Combine(s1).
	//排序，对流的排序，若没有实现sort接口，将不会进行排序
		Sorted().
	//只保留前6个元素
		Limit(6).
	//map 将函数作用到每一个流元素上
		Map(func(i interface{}) interface{} {
			return i.(int)+1
		}).
	//Filter 过滤函数返回值为true的元素，形成新流
		Filter(func(i interface{}) bool {
			return i.(int)==7
		}).
	//Reduce 将前两个元素计算后的结果和第三个元素计算，以此类推
	Reduce(func(i interface{}, i2 interface{}) interface{} {
			return i.(int)+i2.(int)
		})//==20
	//s.Limit(10)
	fmt.Println("after arr: ", s.Collect())
	fmt.Println("after value: ", value)

可使用copy,将流深度拷贝成新的Stream
	a:=[]interface{}{4,3,2,1}
	s:=BuildStream(IntSlice(a))
	s1:=s.Copy()
	s1.Map(func(i interface{}) interface{} {
		return i.(int)+1
	}).Sorted()
	fmt.Println(s1.Collect())
	fmt.Println(s.Collect())

去除重复元素样例
	func (s Item) HashCode() []byte {
		return []byte(strconv.Itoa(int(s)))
	}
	func (s IntSlice) Len() int      { return len(s) }
	func (s IntSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
	func (s IntSlice) Less(i, j int) bool {
		return s[i].(int) < s[j].(int)
	}
	func TestStream_Distinct(t *testing.T) {
		a := []interface{}{Item(1), Item(2), Item(4), Item(4)}
		s := BuildStream(IntSlice(a))
		s.Distinct()
		fmt.Println(s.Collect())
	}
*/
package gostream
