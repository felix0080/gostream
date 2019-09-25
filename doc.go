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

欢迎提Issues
*/
package gostream
