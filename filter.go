package gorm_filter

import (
	"errors"
	"fmt"
	"github.com/imiskolee/form"
	"gorm.io/gorm"
	"strings"
)

type Filter struct {
	TableName string
	Field     string
	OP        string
	Value     string
	db        *gorm.DB
	Handler   FilterHandler
}

func (f *Filter) Copy() Filter {
	return Filter{
		Field:     f.Field,
		TableName: f.TableName,
		OP:        f.OP,
		Value:     f.Value,
	}
}
func (f *Filter) Run(db *gorm.DB) *gorm.DB {
	f.db = db
	if f.Handler != nil {
		return f.Handler(f)(db)
	}
	handler, ok := handlers[strings.ToLower(f.OP)]
	if !ok {
		panic(fmt.Sprintf("Can not support op `%s`", f.OP))
	}
	return handler(f)(db)
}

type FiltersHandler func(*GroupRunner) func(db *gorm.DB) *gorm.DB

type FilterRunner struct {
	filters      []*Filter
	TableName    string
	queryDSL     string
	db           *gorm.DB
	groupHandles map[string]*GroupRunner
}

type GroupRunner struct {
	filters []*Filter
	handel  FiltersHandler
}

func NewFilterDSL(queryDSL string, tableName string, db *gorm.DB) (*FilterRunner, error) {
	s := &FilterRunner{
		queryDSL:     queryDSL,
		db:           db.Table(tableName),
		TableName:    tableName,
		groupHandles: make(map[string]*GroupRunner),
	}
	if err := s.parse(); err != nil {
		return nil, err
	}
	return s, nil
}

type FieldGroupRule func(string) bool

func (s *FilterRunner) Register(field string, handler FilterHandler) *FilterRunner {
	for k, f := range s.filters {
		if f.Field == field {
			f.Handler = handler
			s.filters[k] = f
		}
	}
	return s
}

func (s *FilterRunner) RegisterGroup(groupField string, handler FiltersHandler, groupFilter ...FieldGroupRule) *FilterRunner {
	for k, f := range s.filters {
		if f != nil && (len(groupFilter) > 0 && groupFilter[0](f.Field) || groupField == f.Field) {
			groupHandle, ok := s.groupHandles[groupField]
			if !ok {
				groupHandle = &GroupRunner{
					filters: []*Filter{f},
					handel:  handler,
				}
				s.groupHandles[groupField] = groupHandle
			} else {
				s.groupHandles[groupField].filters = append(groupHandle.filters, f)
			}
			s.filters[k] = nil
		}
	}
	return s
}
func (s *FilterRunner) Run() (*gorm.DB, error) {
	for _, filter := range s.filters {
		if filter != nil {
			s.db = filter.Run(s.db)
		}
	}
	for _, group := range s.groupHandles {
		s.db = group.handel(group)(s.db)
	}
	return s.db, nil
}

func (s *FilterRunner) Get(field string) (Filter, bool) {
	for _, filter := range s.filters {
		if filter.Field == field {
			return *filter, true
		}
	}
	return Filter{}, false
}

func (s *FilterRunner) parse() error {
	form := form.NewForm(s.queryDSL)
	form.NeedQueryUnescape(true)
	data, err := form.Decode()
	if err != nil {
		return err
	}
	filters, ok := data["filter"]
	if !ok {
		return nil
	}
	filterConverter, ok := filters.(map[string]interface{})
	if !ok {
		return errors.New("can not parse filter DSL")
	}
	var filtersDSL []*Filter
	for k, v := range filterConverter {
		f := &Filter{}
		f.Field = k
		val, ok := v.(map[string]interface{})
		if !ok {
			f.OP = "_eq"
			f.Value = fmt.Sprint(v)
			f.TableName = s.TableName
			filtersDSL = append(filtersDSL, f)
			continue
		}
		for op, val := range val {
			f := &Filter{}
			f.Field = k
			f.OP = op
			f.TableName = s.TableName
			f.Value = fmt.Sprint(val)
			filtersDSL = append(filtersDSL, f)
		}
	}
	s.filters = filtersDSL
	return nil
}
