package groupby

func GroupBy[T any](slice []T, isEqual func(a T, b T) bool) [][]T {
	groups := make([][]T, 0)
withNextItem:
	for _, item := range slice {
		for i, group := range groups {
			if isEqual(group[0], item) {
				groups[i] = append(group, item)
				continue withNextItem
			}
		}
		groups = append(groups, []T{item})
	}
	return groups
}
