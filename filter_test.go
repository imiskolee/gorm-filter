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
}

func (t *Table) TableName() string {
	return "table"
}

func TestFilter(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&Table{})

	cases := map[string]string{
		"filter[a][_eq]=a":     "SELECT * FROM `table` WHERE `a` = \"a\"",
		"filter[b][_neq]=b":    "SELECT * FROM `table` WHERE `b` != \"b\"",
		"filter[c]=c":          "SELECT * FROM `table` WHERE `c` = \"c\"",
		"filter[d][_in]=1,2,3": "SELECT * FROM `table` WHERE `d` IN (\"1\",\"2\",\"3\")",
		"filter[e][_null]=1":   "SELECT * FROM `table` WHERE `e` IS NULL",
		"filter[a][_eq]=a&filter[b][_eq]=b&filter[c][_gt]=1": "SELECT * FROM `table` WHERE `a` = \"a\" AND `b` = \"b\" AND `c` > \"1\"",
	}

	for c, _ := range cases {
		parsedDB, err := Parse(c, db.Debug().Table("table"))
		if err != nil {
			t.Fatal(err)
		}
		var lst []Table
		parsedDB.Find(&lst)
	}
}

func TestCase(t *testing.T) {
	db, err := gorm.Open("sqlite", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	filterQuery := "filter[a][_eq]=1"
	var table Table
	db = db.Where("store_id = 1")
	db, err = Parse(filterQuery, db)
	db.Last(&table)
}
