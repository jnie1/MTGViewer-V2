package cards

import (
	"fmt"
	"strings"

	"github.com/mtgjson/mtgjson-sdk-go/db"
)

// adds a where in clause with exactly size 2 tuples
// golang doesn't allow generic typing around different sized arrays
// so the sql is basically always: WHERE (col1, col2) IN ((v1, v2), (v3, v4), ...)
func whereTuples2[C ~[2]string, T ~[2]V, V any](b *db.SQLBuilder, columns C, tuples []T) *db.SQLBuilder {
	placeholders := make([]string, len(tuples))
	for i, tup := range tuples {
		group := [2]string{}
		for j, val := range tup {
			param := b.AddParam(val)
			group[j] = fmt.Sprintf("$%d", param)
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(group[:], ","))
	}
	cond := fmt.Sprintf("(%s) IN (%s)", strings.Join(columns[:], ","), strings.Join(placeholders, ", "))
	b.AddWhere(cond)
	return b
}
