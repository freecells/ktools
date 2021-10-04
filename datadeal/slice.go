package datadeal

import "fmt"

/*
 * 删除Slice中的元素。
 * params:
 *   s: slice对象的指针，如*[]string, *[]int, ...
 *   index: 要删除元素的索引
 * return:
 *   true: 删除成功
 *   false: 删除失败（不支持的数据类型）
 * 说明：直接操作传入的Slice对象，不需要转换为[]interface{}类型。
 */
func SliceRemove(s interface{}, index int) bool {
	if ps, ok := s.(*[]string); ok {
		*ps = append((*ps)[:index], (*ps)[index+1:]...)
	} else if ps, ok := s.(*[]int); ok {
		*ps = append((*ps)[:index], (*ps)[index+1:]...)
	} else if ps, ok := s.(*[]float64); ok {
		*ps = append((*ps)[:index], (*ps)[index+1:]...)
	} else {
		fmt.Printf("<SliceRemove3> Unsupported type: %T\n", s)
		return false
	}

	return true
}

/*
 * 清空Slice，传入的slice对象地址不变。
 * params:
 *   s: slice对象的指针，如*[]string, *[]int, ...
 * return:
 *   true: 清空成功
 *   false: 清空失败（不支持的数据类型）
 */
func SliceClear(s interface{}) bool {
	if ps, ok := s.(*[]string); ok {
		*ps = (*ps)[0:0]
		//*ps = append([]string{})
	} else if ps, ok := s.(*[]int); ok {
		*ps = (*ps)[0:0]
		//*ps = append([]int{})
	} else if ps, ok := s.(*[]float64); ok {
		*ps = (*ps)[0:0]
		//*ps = append([]float64{})
	} else {
		fmt.Printf("<SliceClear3> Unsupported type: %T\n", s)
		return false
	}

	return true
}
