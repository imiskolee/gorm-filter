package gorm_filter

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var handlers map[string]func(*Filter) func(db *gorm.DB) *gorm.DB

func init() {
	handlers = map[string]func(*Filter) func(db *gorm.DB) *gorm.DB{}
	handlers["_eq"] = eqHandler
	handlers["_neq"] = neqHandler
	handlers["_gt"] = gtHandler
	handlers["_gte"] = gteHandler
	handlers["_lt"] = ltHandler
	handlers["_lte"] = lteHandler
	handlers["_in"] = inHandler
	handlers["_not_in"] = notInHandler
	handlers["_null"] = nullHandler
	handlers["_contians"] = containsHandler
}

func eqHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` = ?", f.Field), f.Value)
	}
}

func neqHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` != ?", f.Field), f.Value)
	}
}

func gtHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` > ?", f.Field), f.Value)
	}
}

func gteHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` >= ?", f.Field), f.Value)
	}
}

func ltHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` < ?", f.Field), f.Value)
	}
}

func lteHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` <= ?", f.Field), f.Value)
	}
}

func inHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	vals := strings.Split(f.Value, ",")
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` IN (?)", f.Field), vals)
	}
}

func notInHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	vals := strings.Split(f.Value, ",")
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` NOT IN (?)", f.Field), vals)
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
		return db.Where(fmt.Sprintf("`%s` IS NOT NULL", f.Field))
	}
}

func containsHandler(f *Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", f.Field), fmt.Sprintf("%%%s%%", f.Value))
	}
}
