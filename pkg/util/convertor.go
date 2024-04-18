package util

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// StructToMapIgnoreNil 将 struct 转换成 map[string]interface{}
//
// params: ignoreNil 参数控制是否需要忽略 nil 字段
func StructToMapIgnoreNil(obj interface{}, ignoreNil bool) map[string]interface{} {
	// 获取obj的反射值对象
	objVal := reflect.ValueOf(obj)

	// 处理 obj 是指针的情况，指针需要使用 Elem() 得到其所指向的值
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	// 确保 obj 为结构体
	if objVal.Kind() != reflect.Struct {
		// 如果不是结构体，可以选择抛出一个错误或者返回一个空的map
		panic(fmt.Sprintf("StructToMap error, the obj.Kind must be a struct or ptr type, but go %v.", objVal.Kind()))
	}
	// 创建一个映射，用于存储字段名和值
	result := make(map[string]interface{}, objVal.NumField())
	// 获取结构体类型信息
	objType := objVal.Type()

	// 遍历结构体中的所有字段
	for i := 0; i < objVal.NumField(); i++ {
		// 获取字段的类型信息
		field := objType.Field(i)
		// 获取字段的值
		fieldValue := objVal.Field(i)

		// 如果你想要忽略没有值的字段，可以这么做:
		if ignoreNil && fieldValue.IsZero() {
			continue
		}

		// 将字段名和字段值添加到结果映射中
		result[field.Name] = fieldValue.Interface()
	}

	return result
}

// StructToMap 将 struct 转换成 map[string]interface{}
func StructToMap(obj interface{}) map[string]interface{} {
	return StructToMapIgnoreNil(obj, false)
}

func MapValueToAny(stringMap map[string]string) map[string]interface{} {
	// 创建一个新的 map[string]interface{} 用于存放转换后的键值对
	interfaceMap := make(map[string]interface{}, len(stringMap))

	// 遍历 stringMap 并将值复制到 interfaceMap 中
	for key, value := range stringMap {
		interfaceMap[key] = value
	}

	return interfaceMap
}

// MapToStruct 将 map[string]interface{} 转化为指定类型的对象，其中 result 必须为指针类型并且字段被 mapstructure 标记
func MapToStruct(m map[string]interface{}, result interface{}) error {
	// 设置解码选项，使字符串可以被解码为其他类型
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   result,
		TagName:  "mapstructure",
		// 这里的 WeaklyTypedInput 选项允许弱类型输入
		// 例如，它可以将字符串 "12345" 转换为 int64
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	if err = decoder.Decode(m); err != nil {
		panic(err)
	}

	return nil
}
