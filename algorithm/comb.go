package algorithm

import (
	"strings"
)

func Comb(data []interface{}, k int) [][]interface{} {
	result := make([][]interface{}, 0)
	n := len(data)
	if n == 0 {
		return result
	}
	flagArrs := CombFlag(n, k)
	for _, item := range flagArrs {
		one := make([]interface{}, 0)
		for i := 0; i < n; i++ {
			if item[i:i+1] == "1" {
				one = append(one, data[i])
			}
		}
		result = append(result, one)
	}
	return result
}

/**
 * 获得从 m 中取 n 的所有组合
 * 思路如下：
 * 生成一个长度为 m 的数组，
 * 数组元素的值为 1 表示其下标代表的数被选中，为 0 则没选中。
 *
 * 1. 初始化数组，前 n 个元素置 1，表示第一个组合为前 n 个数；
 * 2. 从左到右扫描数组元素值的 “10” 组合，找到第一个 “10” 组合后将其变为 “01” 组合；
 * 3. 将其左边的所有 “1” 全部移动到数组的最左端
 * 4. 当 n 个 “1” 全部移动到最右端时（没有 “10” 组合了），得到了最后一个组合。
 */
func CombFlag(m, n int) []string {
	resultArrs := make([]string, 0)
	// 先生成一个长度为 m 字符串，开头为 n 个 1， 例如“11100”
	str := strings.Repeat("1", n) + strings.Repeat("0", m-n)
	resultArrs = append(resultArrs, str)

	pos := 0
	keyStr := "10"
	for {
		index := strings.Index(str, keyStr)
		if index == -1 {
			break
		}
		pos = index
		str = strings.Replace(str, keyStr, "01", 1)
		count := strings.Count(str[0:pos], "1")
		str = strings.Repeat("1", count) + strings.Repeat("0", pos+1-count) + str[pos+1:]
		resultArrs = append(resultArrs, str)
	}
	return resultArrs
}
