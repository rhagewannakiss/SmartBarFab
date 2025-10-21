package recommendation

import (
	"github.com/a-palonskaa/SmartBar/tg_bot/internal/coctail"
)

func intersectDrinks(a, b []coctail.Coctail) []coctail.Coctail {
	set := make(map[coctail.Coctail]bool)
	for _, x := range a {
		set[x] = true
	}

	var result []coctail.Coctail
	for _, y := range b {
		if set[y] {
			result = append(result, y)
		}
	}
	return result
}

func RecommendByFilters(filters ...[]coctail.Coctail) []coctail.Coctail {
	if len(filters) == 0 {
		return nil
	}

	result := filters[0]
	for _, f := range filters[1:] {
		result = intersectDrinks(result, f)
	}

	return result
}
