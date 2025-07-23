package utils

type Anies []any

func (a Anies) First() any {
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (a Anies) Last() any {
	if len(a) == 0 {
		return nil
	}
	return a[len(a)-1]
}

func (a Anies) Index(i int) any {
	if len(a) == 0 {
		return nil
	}
	return a[i]
}

func Filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
