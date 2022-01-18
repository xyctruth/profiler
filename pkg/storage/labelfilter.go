package storage

type LabelFilter struct {
	Label
	condition string
}

func (f LabelFilter) Policy(slice1, slice2 []string) []string {
	if f.condition == "AND" {
		return Intersect(slice1, slice2)
	} else if f.condition == "OR" {
		return Union(slice1, slice2)
	}
	return Union(slice1, slice2)
}

//Union 并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}
