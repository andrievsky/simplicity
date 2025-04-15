package utils

type ValueConstraint interface {
	struct{} | bool | string | int64 | float64
}

func MapKeys[K comparable, V ValueConstraint](src map[K]V) []K {
	list := make([]K, len(src))
	i := 0
	for key := range src {
		list[i] = key
		i++
	}
	return list
}

func MapValues[K comparable, V ValueConstraint](src map[K]V) []V {
	list := make([]V, len(src))
	i := 0
	for _, value := range src {
		list[i] = value
		i++
	}
	return list
}
