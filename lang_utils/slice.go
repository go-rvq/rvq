package lang_utils

func GroupBySlice[T comparable, Kt comparable](slices []T, f func(T) Kt) (v [][]T) {
	indexMap := make(map[Kt]int)
	for _, s := range slices {
		key := f(s)
		index, ok := indexMap[key]
		if !ok {
			v = append(v, []T{})
			index = len(v) - 1
			indexMap[key] = index
		}
		v[index] = append(v[index], s)
	}
	return
}

func Map[T any, R any](s []T, f func(T) R) (ret []R) {
	ret = make([]R, len(s))
	for i, v := range s {
		ret[i] = f(v)
	}
	return
}
