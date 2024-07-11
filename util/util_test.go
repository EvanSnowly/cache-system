package util

import (
	"fmt"
	"testing"
)

func TestParseSize(t *testing.T) {
	fmt.Println(uintptr(4))
	// 示例变量
	str := "Hello, world!"
	arr := [3]int{1, 2, 3}
	slice := []int{1, 2, 3, 4, 5}
	m := map[string]int{"one": 1, "two": 2, "three": 3}

	// 打印变量的内存大小
	fmt.Printf("Size of string: %d bytes\n", SizeOfVariable(str))
	fmt.Printf("Size of array: %d bytes\n", SizeOfVariable(arr))
	fmt.Printf("Size of slice: %d bytes\n", SizeOfVariable(slice))
	fmt.Printf("Size of map: %d bytes\n", SizeOfVariable(m))
}
