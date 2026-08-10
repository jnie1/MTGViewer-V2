package cards

import (
	"fmt"
	"strings"

	"github.com/mtgjson/mtgjson-sdk-go/db"
)

type CardIdQuery struct {
	MultiverseIds []int
	SetNumbers    []SetCollectorNumber
	NameSets      []NameSet
}

func (query CardIdQuery) IsEmpty() bool {
	return len(query.MultiverseIds) == 0 && len(query.SetNumbers) == 0 && len(query.NameSets) == 0
}

func WhereIds(b *db.SQLBuilder, ids CardIdQuery, multiverseId, name, setCode, collectorNumber string) *db.SQLBuilder {
	conds := []string{}

	if len(ids.MultiverseIds) > 0 {
		vals := make([]string, len(ids.MultiverseIds))
		for i, id := range ids.MultiverseIds {
			vals[i] = fmt.Sprintf("%d", id)
		}
		conds = append(conds, inValues(b, multiverseId, vals))
	}

	if len(ids.SetNumbers) > 0 {
		columns := []string{setCode, collectorNumber}
		vals := toTuples(ids.SetNumbers, SetCollectorNumber.tuple)
		conds = append(conds, inTuples(b, columns, vals))
	}

	if len(ids.NameSets) > 0 {
		columns := []string{name, setCode}
		vals := toTuples(ids.NameSets, NameSet.tuple)
		conds = append(conds, inTuples(b, columns, vals))
	}

	if len(conds) == 1 {
		b.AddWhere(conds[0])
	} else if len(conds) > 1 {
		b.AddWhere(fmt.Sprintf("(%s)", strings.Join(conds, " OR ")))
	}

	return b
}

func (setNumber SetCollectorNumber) tuple() []string {
	return []string{setNumber.Set, setNumber.CollectorNumber}
}

func (nameSet NameSet) tuple() []string {
	return []string{nameSet.Name, nameSet.Set}
}

func toTuples[V any, T any](values []V, mapper func(V) []T) [][]T {
	tuples := make([][]T, len(values))
	for i, val := range values {
		tuples[i] = mapper(val)
	}
	return tuples
}

func inValues[T any](b *db.SQLBuilder, column string, values []T) string {
	placeholders := make([]string, len(values))
	for i, val := range values {
		param := b.AddParam(val)
		placeholders[i] = fmt.Sprintf("$%d", param)
	}
	return fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ","))
}

func inTuples[T any](b *db.SQLBuilder, columns []string, tuples [][]T) string {
	tupleLen := len(columns)
	for _, t := range tuples {
		if len(t) != tupleLen {
			panic("tuples have unexpected len")
		}
	}

	params := make([]string, len(tuples))
	for i, tup := range tuples {
		placeholders := make([]string, tupleLen)
		for j, val := range tup {
			param := b.AddParam(val)
			placeholders[j] = fmt.Sprintf("$%d", param)
		}
		params[i] = fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
	}

	cond := fmt.Sprintf("(%s) IN (%s)", strings.Join(columns, ","), strings.Join(params, ", "))
	return cond
}
