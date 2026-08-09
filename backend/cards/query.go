package cards

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mtgjson/mtgjson-sdk-go/db"
)

var ErrTupleMismatch = errors.New("tuples have unexpected len")

func WhereInTuples(b *db.SQLBuilder, columns []string, tuples ...[]any) *db.SQLBuilder {
	tupleLen := len(columns)
	for _, t := range tuples {
		if len(t) != tupleLen {
			panic(ErrTupleMismatch)
		}
	}

	if len(tuples) == 0 {
		b.AddWhere("FALSE")
		return b
	}

	params := make([]string, len(tuples))
	for i, tup := range tuples {
		placeholders := make([]string, tupleLen)
		for j, val := range tup {
			param := b.AddParam(val)
			placeholders[j] = fmt.Sprintf("$%d", param)
		}
		params[i] = strings.Join(placeholders, ",")
	}

	cond := fmt.Sprintf("(%s) IN (%s)", strings.Join(columns, ","), strings.Join(params, ","))
	b.AddWhere(cond)

	return b
}
