package scan

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
}

type statusType int8

type Student struct {
	Id        int64
	Name      string
	Cities    []string
	FriendIds pq.Int64Array
	Scores    []interface{}
	Money     decimal.Decimal
	Status    statusType
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func getTestStudents() *sql.Rows {
	return getTestRows(`
select * from (values
(1, '李雷',   '{成都,上海}'::text[], '{1001,1002}'::int[], '["语文",99,"数学",100]'::JSON,
'25.04', 0, '2001-09-01 12:25:48+08'::timestamptz, NULL),
(2, '韩梅梅', '{广州,北京}'::text[], '{1001,1003}'::int[], '["语文",98,"数学",95]'::JSON,
'95.90', 0, '2001-09-01 10:25:48+08'::timestamptz, '2001-09-02 10:25:58+08'::timestamptz)
) as tmp(id, name, cities, friend_ids, scores, money, status, created_at, updated_at)
`)
}

func getTestIntValues() *sql.Rows {
	return getTestRows(`select * from (values (9), (99), (999)) as tmp`)
}

func getTestNull() *sql.Rows {
	return getTestRows(`select null`)
}

func ExampleScan_struct() {
	var row Student
	if err := Scan(getTestStudents(), &row); err != nil {
		log.Panic(err)
	}
	fmt.Printf("{%d %s %v %v %v %v %d\n  %v %v}\n",
		row.Id, row.Name, row.Cities, row.FriendIds, row.Scores, row.Money, row.Status,
		row.CreatedAt, row.UpdatedAt,
	)
	// Output:
	// {1 李雷 [成都 上海] [1001 1002] [语文 99 数学 100] 25.04 0
	//   2001-09-01 12:25:48 +0800 CST <nil>}
}

func ExampleScan_structSlice() {
	var rows []Student
	if err := Scan(getTestStudents(), &rows); err != nil {
		log.Panic(err)
	}
	for _, row := range rows {
		fmt.Printf("{%d %s %v %v %v %v %d\n  %v %v}\n",
			row.Id, row.Name, row.Cities, row.FriendIds, row.Scores, row.Money, row.Status,
			row.CreatedAt, row.UpdatedAt,
		)
	}
	// Output:
	// {1 李雷 [成都 上海] [1001 1002] [语文 99 数学 100] 25.04 0
	//   2001-09-01 12:25:48 +0800 CST <nil>}
	// {2 韩梅梅 [广州 北京] [1001 1003] [语文 98 数学 95] 95.9 0
	//   2001-09-01 10:25:48 +0800 CST 2001-09-02 10:25:58 +0800 CST}
}

func ExampleScan_structPointerSlice() {
	var rows []*Student
	if err := Scan(getTestStudents(), &rows); err != nil {
		log.Panic(err)
	}
	for _, row := range rows {
		fmt.Printf("{%d %s %v %v %v %v %d\n  %v %v}\n",
			row.Id, row.Name, row.Cities, row.FriendIds, row.Scores, row.Money, row.Status,
			row.CreatedAt, row.UpdatedAt,
		)
	}
	// Output:
	// {1 李雷 [成都 上海] [1001 1002] [语文 99 数学 100] 25.04 0
	//   2001-09-01 12:25:48 +0800 CST <nil>}
	// {2 韩梅梅 [广州 北京] [1001 1003] [语文 98 数学 95] 95.9 0
	//   2001-09-01 10:25:48 +0800 CST 2001-09-02 10:25:58 +0800 CST}
}

func ExampleScan_int() {
	var value int
	if err := Scan(getTestIntValues(), &value); err != nil {
		log.Panic(err)
	}
	fmt.Println(value)

	if err := Scan(getTestNull(), &value); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 9
	// sql: Scan error on column index 0: converting driver.Value type <nil> ("<nil>") to a int: invalid syntax
}

func ExampleScan_intPointer() {
	var pointer *int
	if err := Scan(getTestIntValues(), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(*pointer)

	if err := Scan(getTestNull(), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(pointer)

	var p **int
	if err := Scan(getTestIntValues(), &p); err != nil {
		log.Panic(err)
	}
	fmt.Println(**p)

	// Output:
	// 9
	// <nil>
	// 9
}

func ExampleScan_intValueOutOfRange() {
	var value int8
	if err := Scan(getTestRows(`select 128`), &value); err != nil {
		fmt.Println(err)
	}
	// Output:
	// sql: Scan error on column index 0: converting driver.Value type int64 ("128") to a int8: value out of range
}

func ExampleScan_intSlice() {
	var values []int
	if err := Scan(getTestIntValues(), &values); err != nil {
		log.Panic(err)
	}
	fmt.Println(values)
	// Output: [9 99 999]
}

func ExampleScan_pqInt64Array() {
	var values pq.Int64Array
	if err := Scan(getTestRows(`select '{9,99,999}'::int[]`), &values); err != nil {
		log.Panic(err)
	}
	fmt.Println(values)

	if err := Scan(getTestNull(), &values); err != nil {
		log.Panic(err)
	}
	fmt.Println(values)
	// Output:
	// [9 99 999]
	// []
}

func ExampleScan_float() {
	var f float32
	if err := Scan(getTestRows(`select 1.23`), &f); err != nil {
		log.Panic(err)
	}
	fmt.Println(f)
	// Output: 1.23
}

func getTestRows(sql string) *sql.Rows {
	rows, err := db.Query(sql)
	if err != nil {
		log.Panic(err)
	}
	return rows
}
