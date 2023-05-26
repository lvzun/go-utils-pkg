package array

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	globalArray *Array
)

type Array struct {
	data []interface{}
	lock *sync.RWMutex
}

func init() {
	globalArray = New()
}

func New() *Array {
	return &Array{
		data: make([]interface{}, 0),
		lock: &sync.RWMutex{},
	}
}

func (a *Array) Add(value interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.data = append(a.data, value)
}
func (a *Array) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return len(a.data)
}

func (a *Array) Index(value interface{}) int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if value == nil {
		return -1
	}
	for i, v := range a.data {
		if reflect.DeepEqual(v, value) {
			return i
		}
	}
	return -1
}

func (a *Array) RemoveFastByValue(value interface{}) error {
	index := a.Index(value)
	if index == -1 {
		return fmt.Errorf("移除失败，未找到%v值", value)
	}
	a.lock.Lock()
	defer a.lock.Unlock()
	a.RemoveFast(index)
	return nil
}

/*
	此移除方法需要拷贝数据元素，性能相对会慢一些。但不会更改顺序
*/
func (a *Array) RemoveSlow(index int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	size := len(a.data)
	copy(a.data[index:], a.data[index+1:]) // 转移 a[i+1:] 后面的数据.
	a.data[size-1] = ""                    // 清空最一个
	a.data = a.data[:size-1]               // 截断数组
}

/*
	此移除方法会改变数据的顺序，如对顺序有要求，请使用RemoveSlow
*/
func (a *Array) RemoveFast(index int) {
	a.lock.Lock()
	defer a.lock.Unlock()
	size := len(a.data)
	a.data[index] = a.data[size-1] // 将最后一个放i处
	a.data[size-1] = ""            // 清空最一个
	a.data = a.data[:size-1]       // 截断数组
}

func (a *Array) Get(index int) interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.data[index]
}

func (a *Array) GetString(index int) string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.data[index].(string)
}

func (a *Array) GetInt(index int) int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.data[index].(int)
}

func (a *Array) GetFloat64(index int) float64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.data[index].(float64)
}

func Add(value interface{}) {
	globalArray.lock.Lock()
	defer globalArray.lock.Unlock()
	globalArray.data = append(globalArray.data, value)
}

func Len() int {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	return len(globalArray.data)
}

func Index(value interface{}) int {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	if value == nil {
		return -1
	}
	for i, v := range globalArray.data {
		if reflect.DeepEqual(v, value) {
			return i
		}
	}
	return -1
}

func RemoveFastByValue(value interface{}) error {
	index := globalArray.Index(value)
	if index == -1 {
		return fmt.Errorf("移除失败，未找到%v值", value)
	}
	globalArray.lock.Lock()
	defer globalArray.lock.Unlock()
	globalArray.RemoveFast(index)
	return nil
}

/*
	此移除方法需要拷贝数据元素，性能相对会慢一些。但不会更改顺序
*/
func RemoveSlow(index int) {
	globalArray.lock.Lock()
	defer globalArray.lock.Unlock()
	size := len(globalArray.data)
	copy(globalArray.data[index:], globalArray.data[index+1:]) // 转移 a[i+1:] 后面的数据.
	globalArray.data[size-1] = ""                              // 清空最一个
	globalArray.data = globalArray.data[:size-1]               // 截断数组
}

/*
	此移除方法会改变数据的顺序，如对顺序有要求，请使用RemoveSlow
*/
func RemoveFast(index int) {
	globalArray.lock.Lock()
	defer globalArray.lock.Unlock()
	size := len(globalArray.data)
	globalArray.data[index] = globalArray.data[size-1] // 将最后一个放i处
	globalArray.data[size-1] = ""                      // 清空最一个
	globalArray.data = globalArray.data[:size-1]       // 截断数组
}

func Get(index int) interface{} {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	return globalArray.data[index]
}

func GetString(index int) string {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	return globalArray.data[index].(string)
}

func GetInt(index int) int {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	return globalArray.data[index].(int)
}

func GetFloat64(index int) float64 {
	globalArray.lock.RLock()
	defer globalArray.lock.RUnlock()
	return globalArray.data[index].(float64)
}
