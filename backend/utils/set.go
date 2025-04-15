package utils

type Set[K comparable] map[K]struct{}

func SetFromSlice[K comparable](src []K) Set[K] {
	set := make(Set[K], len(src))
	for i := range src {
		set[src[i]] = struct{}{}
		i++
	}
	return set
}
