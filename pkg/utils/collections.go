package utils

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func KeysValues[K comparable, V any](m map[K]V) ([]K, []V) {
	keys, values := make([]K, 0, len(m)), make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// GroupBy 接受一个切片和一个分类函数，返回一个映射
// 其中 K 是键的类型，T 是切片中元素的类型
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	grouped := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		grouped[key] = append(grouped[key], item)
	}
	return grouped
}
