package sqlBuild

import (
	"reflect"
	"fmt"
	"strings"
	"alleyFunAdmin/component/sql-build/debug"
	"sync"
)

type InsertBuild struct {
	BuildCore
}

func (i *InsertBuild) Insert(table string) InsertInf {
	i.setTabName(table)
	return i
}
func (i *InsertBuild) Value(value interface{}, rules ... Rule) InsertInf {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("The struct paramer must be use ptr"))
	}
	ind := val.Elem()
	tag := "insert"
	i.setValueColumns(ind.Type(), tag)
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	i.value(ind, rule, )
	return i
}
func (i *InsertBuild) NoOption(noOptions ...string) InsertInf {
	if len(noOptions) > 0 {
		i.setNoOptions(noOptions)
	}
	return i
}

func (i *InsertBuild) Option(options ...string) InsertInf {
	if len(options) > 0 {
		i.setOptions(options)
	}
	return i
}

func (i *InsertBuild) Values(value interface{}, rules ... Rule) InsertInf {
	vals := reflect.Indirect(reflect.ValueOf(value))
	if vals.Kind() != reflect.Slice {
		panic(fmt.Errorf("The struct paramer must be use slice"))
	}
	if vals.Len() <= 0 {
		i.err = ErrInsertValue
		return i
	}
	structVal := reflect.Indirect(vals.Index(0))
	tag := "insert"
	i.setValueColumns(structVal.Type(), tag)
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	var wg sync.WaitGroup
	wg.Add(vals.Len())
	for index := 0; index < vals.Len(); index++ {
		go func(value reflect.Value) {
			i.value(value, rule, &wg)
		}(reflect.Indirect(vals.Index(index)))
	}
	wg.Wait()
	return i
}

func (i *InsertBuild) OrUpdate() InsertInf {
	i.isOrUpdate = true
	return i
}

func (i *InsertBuild) String() (string, error) {
	if i.err != nil {
		return "", i.err
	}
	//table
	if i.tableName == "" {
		return "", ErrTabName
	}
	//insertColumn
	var insertColumn string
	if len(i.insertColumns) > 0 {
		insertColumn = strings.Join(i.insertColumns, ",")
	}
	//insertValues
	var insertValue string
	if len(i.insertValues) > 1000 {
		debug.Warning("insert datas >1000")
	}
	insertValue = strings.Join(i.insertValues, "),(")

	if insertColumn == "" {
		return "", ErrInsertColumn
	}
	if insertValue == "" {
		return "", ErrInsertValue
	}
	//orupdate
	var orUpdate string
	if i.isOrUpdate {
		var orUpdates []string
		for _, v := range i.insertColumns {
			temp := strings.Join([]string{v, v}, " = values(")
			orUpdates = append(orUpdates, temp)
		}
		orUpdate = " ON DUPLICATE KEY UPDATE " + strings.Join(orUpdates, "),") + ")"
	}

	sql := "INSERT INTO " + i.tableName + "(" + insertColumn + ") VALUES (" + insertValue + ") " + orUpdate
	debug.Println("sql:" + sql)
	return sql, nil
}
func Insert(table string) InsertInf {
	auto := new(InsertBuild)
	auto.Insert(table)
	return auto
}
