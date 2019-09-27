 [![Fork me on Gitee](https://gitee.com/softbar/gostream/widgets/widget_2.svg)](https://gitee.com/softbar/gostream)

# gostream
go Map Reduce Sorted Filter Limit stream , a simple library  
本库可用于科学计算，数据分析等方向，欢迎贡献或提需求

# 进度  已实现
  * Limit 支持流截断             
  * Combine 支持流之间的拼接     
  * Filter 支持流过滤
  * Sorted 支持流排序
  * Reduce 支持流 Reduce
  * Map 支持流 Map
  * MultipartMap 支持并发流  
  * Copy 流拷贝
  * Distinct 排除重复元素  拟定适用技术 bloom 过滤器 / Cuckoo Filter 过滤器
 # 进度  待实现
  * Manager 流的管理器，管理多个流，做到简单，快捷操作多个流
  * 流之间拼接允许按照优先级进行排列（和优先级队列相似，用于优先级重要的场景）
  
  
  
  
 [![GoDoc](https://godoc.org/github.com/felix0080/gostream?status.png)](https://godoc.org/github.com/felix0080/gostream) 


  Documentation here: https://godoc.org/github.com/felix0080/gostream

# Usage
 ```go
package main
import (
	"fmt"
	"strconv"
	. "github.com/felix0080/gostream"
)
type IntSlice []interface{}

func (s IntSlice) Len() int { return len(s) }
func (s IntSlice) Swap(i, j int){ s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool {
	return s[i].(int) < s[j].(int)
} 
type Item int
func (s Item) HashCode() []byte {
    return []byte(strconv.Itoa(int(s)))
}
func main() {
   //调用BuildStream 方法，传入实现了sort 接口的slice，若不需要排序，则无需实现sort
    //定义一个slice.填充元素
    a:=[]interface{}{4,3,2,1}
    //使用实现了sort的type 包装，并调用调用BuildStream
    s,err:=BuildStream(IntSlice(a))
    if err != nil {
        fmt.Println(err)
        return
    }
    //定义一个slice.填充元素
    a1:=[]interface{}{6,5,8,7}
    //使用实现了sort的type 包装,并调用调用BuildStream
    s1,err:=BuildStream(IntSlice(a1))
    if err != nil {
        fmt.Println(err)
        return
    }
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
   
    s2:=s.Copy()
    s2.Map(func(i interface{}) interface{} {
        return i.(int)+1
    }).Sorted()
    fmt.Println(s2.Collect())
    fmt.Println(s.Collect())
    a3 := []interface{}{Item(1), Item(2), Item(4), Item(4)}
    s3, err := BuildStream(IntSlice(a3))
    if err != nil {
        fmt.Println(err)
        return
    }
    s3.Distinct()
    fmt.Println(s3.Collect())
 }
   
```
