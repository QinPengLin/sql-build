package sqlBuild

import (
	"strings"
	"strconv"
	"alleyFunAdmin/component/sql-build/debug"
)

type DeleteBuild struct {
	BuildCore
}

func (d *DeleteBuild) Delete(table string) DeleteInf {
	d.setTabName(table)
	return d
}

func (s *DeleteBuild) WhereFunc(value interface{}, key string, rules ...Rule) DeleteInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	s.where(value, key, rule, GetWhereSetFuncValues)
	return s
}

func (d *DeleteBuild) Where(value interface{}, key string, rules ... Rule) DeleteInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	d.where(value, key, rule,GetWhereSetValues)
	return d
}
func (d *DeleteBuild) Where_(value interface{}, key string, rules ... Rule) DeleteInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	d.where_(value, key, rule,GetWhereSetValues)
	return d
}
func (d *DeleteBuild) Like(value string, key string) DeleteInf {
	d.like(value, key)
	return d
}
func (d *DeleteBuild) In(values interface{}, key string) DeleteInf {
	d.in(values, key)
	return d
}
func (d *DeleteBuild) NotIn(values interface{}, key string) DeleteInf {
	d.notin(values, key)
	return d
}
func (d *DeleteBuild) OrderBy(orderBy string) DeleteInf {
	d.orderBy(orderBy)
	return d
}
func (d *DeleteBuild) Limit(limit int) DeleteInf {
	d.limit(limit)
	return d
}
func (d *DeleteBuild) Offset(offset int) DeleteInf {
	d.offset(offset)
	return d
}
func (d *DeleteBuild) GroupBy(groupBy string) DeleteInf {
	d.groupBy(groupBy)
	return d
}
func (d *DeleteBuild) String() (string, error) {
	if d.err != nil {
		return "", d.err
	}
	//table
	if d.tableName == "" {
		return "", ErrTabName
	}
	//where
	var where string
	if len(d.whereValues) > 0 {
		where = strings.Join(d.whereValues, " and ")
	}
	//in notin
	var in, notin string
	var inTemp, notinTemp []string
	for k, v := range d.inMap {
		inValues := k + " IN (" + strings.Join(v, ",") + ")"
		inTemp = append(inTemp, inValues)
	}
	for k, v := range d.notinMap {
		notinValues := k + " NOT IN (" + strings.Join(v, ",") + ")"
		notinTemp = append(notinTemp, notinValues)
	}
	in = strings.Join(inTemp, " and ")
	notin = strings.Join(notinTemp, " and ")

	//like
	var like string
	if len(d.likeValues) > 0 {
		like = strings.Join(d.likeValues, " and ")
	}

	var groupBy string
	//groupby
	if len(d.groupByValues) > 0 {
		groupBy = " GROUP BY " + strings.Join(d.groupByValues, ",")
	}
	//orderby
	var orderBy string
	if len(d.orderValues) > 0 {
		orderBy = " ORDER BY " + strings.Join(d.orderValues, ",")
	}
	//limit offset
	var limits string
	if d.limitValue > 0 {
		limits = " LIMIT " + strconv.Itoa(d.limitValue)
		if d.offsetValue > 0 {
			limits += " OFFSET " + strconv.Itoa(d.offsetValue)
		}
	} else if d.offsetValue > 0 {
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
	sql := "DELETE FROM " + d.tableName + lastWhere + groupBy + orderBy + limits
	debug.Println("sql:" + sql)
	return sql, nil
}

func Delete(table string) DeleteInf {
	auto := new(DeleteBuild)
	auto.Delete(table)
	return auto
}
