package sqlBuild

import (
	"strings"
	"strconv"
	"alleyFunAdmin/component/sql-build/debug"
)

type UpdateBuild struct {
	BuildCore
}

func (u *UpdateBuild) Update(table string) UpdateInf {
	u.setTabName(table)
	return u
}
func (u *UpdateBuild) Set(value interface{}, key string, rules ... Rule) UpdateInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	u.set(value, key, rule)
	return u
}
func (u *UpdateBuild) Set_(value interface{}, key string, rules ... Rule) UpdateInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	u.set_(value, key, rule)
	return u
}

func (s *UpdateBuild) WhereFunc(value interface{}, key string, rules ...Rule) UpdateInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	s.where(value, key, rule, GetWhereSetFuncValues)
	return s
}

func (u *UpdateBuild) Where(value interface{}, key string, rules ... Rule) UpdateInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	u.where(value, key, rule, GetWhereSetValues)
	return u
}
func (u *UpdateBuild) Where_(value interface{}, key string, rules ... Rule) UpdateInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	u.where_(value, key, rule, GetWhereSetValues)
	return u
}
func (u *UpdateBuild) Like(value string, key string) UpdateInf {
	u.like(value, key)
	return u
}
func (u *UpdateBuild) In(values interface{}, key string) UpdateInf {
	u.in(values, key)
	return u
}
func (u *UpdateBuild) NotIn(values interface{}, key string) UpdateInf {
	u.notin(values, key)
	return u
}
func (u *UpdateBuild) OrderBy(orderBy string) UpdateInf {
	u.orderBy(orderBy)
	return u
}
func (u *UpdateBuild) Limit(limit int) UpdateInf {
	u.limit(limit)
	return u
}
func (u *UpdateBuild) Offset(offset int) UpdateInf {
	u.offset(offset)
	return u
}
func (u *UpdateBuild) GroupBy(groupBy string) UpdateInf {
	u.groupBy(groupBy)
	return u
}
func (s *UpdateBuild) String() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	//table
	if s.tableName == "" {
		return "", ErrTabName
	}
	//set
	var set string
	if len(s.setValues) > 0 {
		set = " SET " + strings.Join(s.setValues, ",")
	} else {
		return "", ErrNoUpdate
	}

	//where
	var where string
	if len(s.whereValues) > 0 {
		where = strings.Join(s.whereValues, " and ")
	}
	//in notin
	var in, notin string
	var inTemp, notinTemp []string
	for k, v := range s.inMap {
		inValues := k + " IN (" + strings.Join(v, ",") + ")"
		inTemp = append(inTemp, inValues)
	}
	for k, v := range s.notinMap {
		notinValues := k + " NOT IN (" + strings.Join(v, ",") + ")"
		notinTemp = append(notinTemp, notinValues)
	}
	in = strings.Join(inTemp, " and ")
	notin = strings.Join(notinTemp, " and ")

	//like
	var like string
	if len(s.likeValues) > 0 {
		like = strings.Join(s.likeValues, " and ")
	}
	//groupby
	var groupBy string
	if len(s.groupByValues) > 0 {
		groupBy = " GROUP BY " + strings.Join(s.groupByValues, ",")
	}
	//orderby
	var orderBy string
	if len(s.orderValues) > 0 {
		orderBy = " ORDER BY " + strings.Join(s.orderValues, ",")
	}
	//limit offset
	var limits string
	if s.limitValue > 0 {
		limits = " LIMIT " + strconv.Itoa(s.limitValue)
		if s.offsetValue > 0 {
			limits += " OFFSET " + strconv.Itoa(s.offsetValue)
		}
	} else if s.offsetValue > 0 {
		return "", ErrNoLimit
	}

	var wheres []string
	if where != "" {
		wheres = append(wheres, where)
	}
	if in != "" {
		wheres = append(wheres, in)
	}
	if notin != "" {
		wheres = append(wheres, notin)
	}
	if like != "" {
		wheres = append(wheres, like)
	}
	var lastWhere string
	if len(wheres) > 0 {
		lastWhere = " WHERE " + strings.Join(wheres, " and ")
	}
	sql := "UPDATE " + s.tableName + set + lastWhere + groupBy + orderBy + limits
	debug.Println("sql:" + sql)
	return sql, nil
}

func Update(table string) UpdateInf {
	auto := new(UpdateBuild)
	auto.Update(table)
	return auto
}
