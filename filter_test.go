package gorm_filter

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
)

type Table struct {
	ID int
	A  int
	B  int
	C  int
	D  int
	E  int
}

type TableB struct {
	ID int
	A  int
	B  int
	C  int
}

func (t *TableB) TableName() string {
	return "table_b"
}

func (t *Table) TableName() string {
	return "table"
}

func TestFilter(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test1.db")
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&Table{}, &TableB{})

	cases := map[string]string{
		"filter[a][_eq]=a":     "SELECT * FROM `table` WHERE `a` = \"a\"",
		"filter[b][_neq]=b":    "SELECT * FROM `table` WHERE `b` != \"b\"",
		"filter[c]=c":          "SELECT * FROM `table` WHERE `c` = \"c\"",
		"filter[d][_in]=1,2,3": "SELECT * FROM `table` WHERE `d` IN (\"1\",\"2\",\"3\")",
		"filter[e][_null]=1":   "SELECT * FROM `table` WHERE `e` IS NULL",
		"filter[a][_eq]=a&filter[b][_eq]=b&filter[c][_gt]=1": "SELECT * FROM `table` WHERE `a` = \"a\" AND `b` = \"b\" AND `c` > \"1\"",
	}

	for c, _ := range cases {
		runner, err := NewFilterDSL(c, "table", db)
		if err != nil {
			t.Fatal(err)
		}
		parsedDB, err := runner.Run()
		if err != nil {
			t.Fatal(err)
		}
		var lst []Table
		if err := parsedDB.Debug().Find(&lst).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func TestJoinCase(t *testing.T) {
	dsl := "filter[a][_eq]=a&filter[bc][_eq]=1"
	db, err := gorm.Open("sqlite3", "test1.db")
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&Table{}, &TableB{})
	runner, err := NewFilterDSL(dsl, "table", db)
	runner.Register("bc", func(f *Filter) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			newF := f.Copy()
			newF.TableName = "table_b"
			newF.Field = "c"
			db = db.Joins("LEFT JOIN table_b on table_b.a = `table`.id")
			return newF.Run(db)
		}
	})
	parsedDB, err := runner.Run()
	if err != nil {
		t.Fatal(err)
	}
	var lst []Table
	if err := parsedDB.Debug().Find(&lst).Error; err != nil {
		t.Fatal(err)
	}
}
