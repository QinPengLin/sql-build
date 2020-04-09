package sqlBuild

import (
	"alleyFunAdmin/component/sql-build/debug"
)

type SelectInf interface {
	Select(table string) SelectInf
	Column(columns ... string) SelectInf
	Where(value interface{}, key string, rules ... Rule) SelectInf
	Where_(value interface{}, key string, rules ... Rule) SelectInf
	WhereString(value interface{}, rules ... Rule) SelectInf
	WhereMap(value map[string]interface{}, rules ... Rule) SelectInf
	WhereFunc(value interface{}, key string, rules ...Rule) SelectInf
	Like(value string, key string) SelectInf
	In(values interface{}, key string) SelectInf
	NotIn(values interface{}, key string) SelectInf
	OrderBy(orderBy string) SelectInf
	OrderByArr(orderBy []string) SelectInf
	Limit(limit int) SelectInf
	Offset(offset int) SelectInf
	GroupBy(groupBy string) SelectInf
	String() (string, error)
}

type InsertInf interface {
	Insert(table string) InsertInf
	Option(options ...string) InsertInf
	NoOption(noOptions ...string) InsertInf
	Value(value interface{}, rules ... Rule) InsertInf
	Values(value interface{}, rules ... Rule) InsertInf
	OrUpdate() InsertInf
	String() (string, error)
}

type UpdateInf interface {
	Update(table string) UpdateInf
	Set(value interface{}, key string, rules ... Rule) UpdateInf
	Set_(value interface{}, key string, rules ... Rule) UpdateInf
	Where(value interface{}, key string, rules ... Rule) UpdateInf
	Where_(value interface{}, key string, rules ... Rule) UpdateInf
	WhereFunc(value interface{}, key string, rules ...Rule) UpdateInf
	Like(value string, key string) UpdateInf
	In(values interface{}, key string) UpdateInf
	NotIn(values interface{}, key string) UpdateInf
	OrderBy(orderBy string) UpdateInf
	Limit(limit int) UpdateInf
	Offset(offset int) UpdateInf
	GroupBy(groupBy string) UpdateInf
	String() (string, error)
}

type DeleteInf interface {
	Delete(table string) DeleteInf
	Where(value interface{}, key string, rules ... Rule) DeleteInf
	Where_(value interface{}, key string, rules ... Rule) DeleteInf
	WhereFunc(value interface{}, key string, rules ...Rule) DeleteInf
	Like(value string, key string) DeleteInf
	In(values interface{}, key string) DeleteInf
	NotIn(values interface{}, key string) DeleteInf
	OrderBy(orderBy string) DeleteInf
	Limit(limit int) DeleteInf
	Offset(offset int) DeleteInf
	GroupBy(groupBy string) DeleteInf
	String() (string, error)
}

func Debug(deb ...bool) {
	d := true
	if len(deb) > 0 {
		d = deb[0]
	}
	debug.Debug = d
}
