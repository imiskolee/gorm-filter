package gorm_filter

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type FilterHandler func(*Filter) func(db *gorm.DB) *gorm.DB

var handlers map[string]FilterHandler

func init() {
	handlers = map[string]FilterHandler{}
	handlers["_eq"] = eqHandler
	handlers["_neq"] = neqHandler
	handlers["_gt"] = gtHandler
	handlers["_gte"] = gteHandler
	handlers["_lt"] = ltHandler
	handlers["_lte"] = lteHandler
	handlers["_in"] = inHandler
	handlers["_not_in"] = notInHandler
	handlers["_null"] = nullHandler
	handlers["_contains"] = containsHandler
	handlers["_not_contains"] = notContainsHandler
}

func eqHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` = ?", f.TableName, f.Field), f.Value)
	}
}

func neqHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` != ?", f.TableName, f.Field), f.Value)
	}
}

func gtHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` > ?", f.TableName, f.Field), f.Value)
	}
}

func gteHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` >= ?", f.TableName, f.Field), f.Value)
	}
}

func ltHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` < ?", f.TableName, f.Field), f.Value)
	}
}

func lteHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` <= ?", f.TableName, f.Field), f.Value)
	}
}

func inHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	vals := strings.Split(f.Value, ",")
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` IN (?)", f.TableName, f.Field), vals)
	}
}

func notInHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	vals := strings.Split(f.Value, ",")
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` NOT IN (?)", f.TableName, f.Field), vals)
	}
}

func nullHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	isNull, err := strconv.ParseBool(f.Value)
	if err != nil {
		panic(fmt.Sprintf("Can not parse `%s` to bool", f.Value))
	}
	return func(db *gorm.DB) *gorm.DB {
		if isNull {
			return db.Where(fmt.Sprintf("`%s` IS NULL", f.Field))
		}
		return db.Where(fmt.Sprintf("`%s`.`%s` IS NOT NULL", f.TableName, f.Field))
	}
}

func containsHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` LIKE ?", f.TableName, f.Field), fmt.Sprintf("%%%s%%", f.Value))
	}
}

func notContainsHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s`.`%s` not LIKE ?", f.TableName, f.Field), fmt.Sprintf("%%%s%%", f.Value))
	}
}
