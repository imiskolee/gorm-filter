package gorm_filter

import (
	"errors"
	"fmt"
	"github.com/imiskolee/form"
	"github.com/jinzhu/gorm"
	"strings"
)

type Filter struct {
	Field string
	OP    string
	Value string
}

func (f *Filter) Run(db *gorm.DB) *gorm.DB {
	handler, ok := handlers[strings.ToLower(f.OP)]
	if !ok {
		panic(fmt.Sprintf("Can not support op `%s`", f.OP))
	}
	return handler(f)(db)
}

func Parse(queryString string, db *gorm.DB) (*gorm.DB, error) {
	form := form.NewForm(queryString)
	form.NeedQueryUnescape(true)
	data, err := form.Decode()
	if err != nil {
		return nil, err
	}
	filters, ok := data["filter"]
	if !ok {
		return db, nil
	}
	filterConverter, ok := filters.(map[string]interface{})
	if !ok {
		return nil, errors.New("can not parse filter DSL")
	}
	var filtersDSL []Filter
	for k, v := range filterConverter {
		var f Filter
		f.Field = k
		val, ok := v.(map[string]interface{})
		if !ok {
			f.OP = "_eq"
			f.Value = fmt.Sprint(v)
			filtersDSL = append(filtersDSL, f)
			continue
		}
		for op, val := range val {
			var f Filter
			f.Field = k
			f.OP = op
			f.Value = fmt.Sprint(val)
			filtersDSL = append(filtersDSL, f)
		}
	}

	for _, filter := range filtersDSL {
		db = filter.Run(db)
	}
	return db, nil
}
