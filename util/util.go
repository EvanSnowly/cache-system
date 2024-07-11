package util

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const (
	BYTE = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

func ParseSize(size string) (int64, string, error) {
	var n, u string
	for _, value := range size {
		if unicode.IsDigit(value) {
			n += string(value)
		} else if unicode.IsLetter(value) {
			u += string(value)
		}
	}
	u = strings.ToUpper(u)
	dn, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		log.Println("The unit is wrong, and the default setting is 100 megabytes")
		return 100, "100MB", err
	}
	var bs int64 = 0
	switch u {
	case "B":
		bs = dn
	case "KB":
		bs = dn * KB
	case "MB":
		bs = dn * MB
	case "GB":
		bs = dn * GB
	case "TB":
		bs = dn * TB
	default:
		log.Println("The unit is wrong, and the default setting is 100 megabytes")
		return 100, "100MB", err
	}

	return bs, n + u, nil
}

// SizeOfVariable 返回一个变量的内存大小（以字节为单位）
func SizeOfVariable(v any) uintptr {
	val := reflect.ValueOf(v)
	return calculateSize(val)
}

// calculateSize 递归计算变量的内存大小
func calculateSize(val reflect.Value) uintptr {
	switch val.Kind() {
	case reflect.String:
		return uintptr(val.Len())
	case reflect.Slice:
		var size uintptr
		for i := 0; i < val.Len(); i++ {
			size += calculateSize(val.Index(i))
		}
		return size
	case reflect.Array:
		var size uintptr
		for i := 0; i < val.Len(); i++ {
			size += calculateSize(val.Index(i))
		}
		return size
	case reflect.Map:
		var size uintptr
		for _, key := range val.MapKeys() {
			size += calculateSize(key)
			size += calculateSize(val.MapIndex(key))
		}
		return size
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return 0
		}
		return calculateSize(val.Elem())
	case reflect.Struct:
		var size uintptr
		for i := 0; i < val.NumField(); i++ {
			size += calculateSize(val.Field(i))
		}
		return size
	default:
		return uintptr(unsafe.Sizeof(val.Interface()))
	}
}
