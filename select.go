package sqlBuild

import (
	"alleyFunAdmin/component/sql-build/debug"
	"strconv"
	"strings"
)

type SelectBuild struct {
	BuildCore
}

func (s *SelectBuild) WhereFunc(value interface{}, key string, rules ...Rule) SelectInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	s.where(value, key, rule, GetWhereSetFuncValues)
	return s
}


func (s *SelectBuild) Select(table string) SelectInf {
	s.setTabName(table)
	return s
}

func (s *SelectBuild) Column(columns ... string) SelectInf {
	for _, column := range columns {
		s.column(column)
	}
	return s
}
func (s *SelectBuild) Where(value interface{}, key string, rules ... Rule) SelectInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	s.where(value, key, rule, GetWhereSetValues)
	return s
}
func (s *SelectBuild) Where_(value interface{}, key string, rules ... Rule) SelectInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	s.where_(value, key, rule, GetWhereSetValues)
	return s
}
func (s *SelectBuild) WhereString(value interface{}, rules ... Rule) SelectInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	if value=="" {
		return s
	}
	s.whereString(value, rule, GetWhereSetValues)
	return s
}
func (s *SelectBuild) WhereMap(value map[string]interface{}, rules ... Rule) SelectInf {
	var rule Rule
	if len(rules) > 0 {
		rule = rules[0]
	}
	//处理map
	valueString:=MapToString(value)
	if valueString == "" {
		return s
	}
	s.whereString(valueString, rule, GetWhereSetValues)
	return s
}
func (s *SelectBuild) Like(value string, key string) SelectInf {
	s.like(value, key)
	return s
}
func (s *SelectBuild) In(values interface{}, key string) SelectInf {
	s.in(values, key)
	return s
}
func (s *SelectBuild) NotIn(values interface{}, key string) SelectInf {
	s.notin(values, key)
	return s
}
func (s *SelectBuild) OrderBy(orderBy string) SelectInf {
	s.orderBy(orderBy)
	return s
}
func (s *SelectBuild) OrderByArr(orderBy []string) SelectInf {
	s.orderByArr(orderBy)
	return s
}
func (s *SelectBuild) Limit(limit int) SelectInf {
	s.limit(limit)
	return s
}
func (s *SelectBuild) Offset(offset int) SelectInf {
	s.offset(offset)
	return s
}
func (s *SelectBuild) GroupBy(groupBy string) SelectInf {
	s.groupBy(groupBy)
	return s
}
func (s *SelectBuild) String() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	//table
	if s.tableName == "" {
		return "", ErrTabName
	}
	//column
	var column string
	if len(s.columnValues) == 0 {
		column = " *"
	} else {
		column = " " + strings.Join(s.columnValues, ",")
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
		if strings.Contains(strings.ToLower(column), "count(") {
			debug.Warning("'count' and 'limit'  cannot exist at the same time")
		}
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
	sql := "SELECT" + column + " FROM " + s.tableName + lastWhere + groupBy + orderBy + limits
	debug.Println("sql:" + sql)
	return sql, nil
}

func Select(table string) SelectInf {
	auto := new(SelectBuild)
	auto.Select(table)
	return auto
}
