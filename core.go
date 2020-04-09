package sqlBuild

import (
	"errors"
	"strings"
	"alleyFunAdmin/component/sql-build/debug"
	"reflect"
	"sync"
)

var (
	ErrTabName      = errors.New("The tabName can not be empty")
	ErrValueType    = errors.New("The Value not have need type")
	ErrInjection    = errors.New("Injection err")
	ErrNoUpdate     = errors.New("Not Found Update Data")
	ErrCondition    = errors.New("Fail to meet the condition")
	ErrSet          = errors.New("Fail to meet the set")
	ErrNoLimit      = errors.New("Need 'Offset' and 'Limit' are used together")
	ErrInsertColumn = errors.New("Not found Insert Column")
	ErrInsertValue  = errors.New("Not found Insert Data")
)

type BuildCore struct {
	tableName     string
	columnValues  []string //查询列
	whereValues   []string //条件值
	setValues     []string //修改值
	orderValues   []string
	limitValue    int
	offsetValue   int
	inMap         map[string][]string
	notinMap      map[string][]string
	groupByValues []string
	likeValues    []string

	//insert
	insertOptions     []string
	insertNoOptions   []string
	insertColumns     []string
	insertAutoIndex   int    //自增列的field下标
	insertMyCatValue  string //mycat id列的函数名称
	insertTags        []int
	insertValues      []string
	isOrUpdate        bool
	insertValuesMutex sync.Mutex

	err error
}
type Rule struct {
	IntValue     int
	Int8Value    int8
	Int16Value   int16
	Int32Value   int32
	Int64Value   int64
	UintValue    uint
	Uint8Value   uint8
	Uint16Value  uint16
	Uint32Value  uint32
	Uint64Value  uint64
	Float32Value float32
	Float64Value float64
	StringValue  string
}

func (b *BuildCore) setTabName(tabName string) {
	if CheckInjection(tabName) {
		b.err = ErrInjection
		return
	}
	if tabName == "" {
		b.err = ErrTabName
		debug.Warning("The tabName can not be empty")
	}
	b.tableName = tabName
}

func (b *BuildCore) orderBy(orderByValue string) {
	if b.err != nil {
		return
	}
	if orderByValue == "" {
		debug.Println("The orderByValue is nil")
		return
	}
	if CheckInjection(orderByValue) {
		b.err = ErrInjection
		return
	}
	b.orderValues = append(b.orderValues, orderByValue)
}

func (b *BuildCore) orderByArr(orderByValue []string) {
	if b.err != nil {
		return
	}
	if len(orderByValue) == 0 {
		debug.Println("The orderByValue is nil")
		return
	}

	for _, v := range orderByValue{
		if CheckInjection(v) {
			b.err = ErrInjection
			return
		}
		b.orderValues = append(b.orderValues, OrderByArrString(v))
	}
}

func (b *BuildCore) groupBy(groupByValue string) {
	if b.err != nil {
		return
	}
	if groupByValue == "" {
		debug.Println("The groupByValue is nil")
		return
	}
	if CheckInjection(groupByValue) {
		b.err = ErrInjection
		return
	}
	b.groupByValues = append(b.groupByValues, groupByValue)
}
func (b *BuildCore) like(likeValue, key string) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The LikeKey can not be empty")
		return
	}
	if likeValue == "" || strings.Count(likeValue, "%") == len(likeValue) {
		debug.Println("The likeValue is nil")
		return
	}
	if CheckInjection(likeValue) {
		b.err = ErrInjection
		return
	}
	if strings.Contains(likeValue, "%") {
		b.likeValues = append(b.likeValues, key+" like '"+likeValue+"'")
	} else {
		b.likeValues = append(b.likeValues, key+" like "+strings.Join([]string{"'%", "%'"}, likeValue))
	}
}

func (b *BuildCore) in(inValues interface{}, key string) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The InKey can not be empty")
		return
	}
	if b.inMap == nil {
		b.inMap = make(map[string][]string)
	}
	result, err := GetInValues(inValues)
	if err != nil {
		b.err = err
		return
	}
	if len(result) > 0 {
		b.inMap[key] = result
	}
}

func (b *BuildCore) notin(notinValues interface{}, key string) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The NotInKey can not be empty")
		return
	}
	if b.notinMap == nil {
		b.notinMap = make(map[string][]string)
	}
	result, err := GetInValues(notinValues)
	if err != nil {
		b.err = err
		return
	}
	if len(result) > 0 {
		b.notinMap[key] = result
	}
}

func (b *BuildCore) where(whereValue interface{}, key string, rule Rule, f func(values interface{}, rule Rule) (value string,
	err error)) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The 'WhereKey' can not be empty")
		return
	}
	value, err := f(whereValue, rule)
	if err != nil {
		b.err = err
		return
	}
	if value != rule.StringValue && value != strings.Join([]string{"'", "'"}, rule.StringValue) {
		if !strings.ContainsAny(key, ">=<") {
			key += " = "
		}
		b.whereValues = append(b.whereValues, key+value)
	}
}

func (b *BuildCore) where_(whereValue interface{}, key string, rule Rule, f func(values interface{}, rule Rule) (value string,
	err error)) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The WhereKey can not be empty")
		return
	}
	value, err := f(whereValue, rule)
	if err != nil {
		b.err = err
		return
	}
	if value != rule.StringValue && value != strings.Join([]string{"'", "'"}, rule.StringValue) {
		if !strings.ContainsAny(key, ">=<") {
			key += " = "
		}
		b.whereValues = append(b.whereValues, key+value)
	} else {
		b.err = ErrCondition
	}
}

//直接支持字符串条件
func (b *BuildCore) whereString(whereValue interface{}, rule Rule, f func(values interface{}, rule Rule)(value string,
	err error)){
	if b.err != nil {
		return
	}
	value, err := f(whereValue, rule)
	if err != nil {
		b.err = err
		return
	}
	if value != rule.StringValue && value != strings.Join([]string{"'", "'"}, rule.StringValue) {
		b.whereValues = append(b.whereValues, strings.Trim(value, "'"))
	} else {
		b.err = ErrCondition
	}
}

func (b *BuildCore) set(setValue interface{}, key string, rule Rule) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The 'SetKey' can not be empty")
		return
	}
	value, err := GetWhereSetValues(setValue, rule)
	if err != nil {
		b.err = err
		return
	}
	if value != rule.StringValue && value != strings.Join([]string{"'", "'"}, rule.StringValue) {
		if !strings.ContainsAny(key, ">=<") {
			key += " = "
		}
		b.setValues = append(b.setValues, key+value)
	}
}
func (b *BuildCore) set_(setValue interface{}, key string, rule Rule) {
	if b.err != nil {
		return
	}
	if key == "" {
		debug.Println("The 'SetKey' can not be empty")
		return
	}
	value, err := GetWhereSetValues(setValue, rule)
	if err != nil {
		b.err = err
		return
	}
	if value != rule.StringValue && value != strings.Join([]string{"'", "'"}, rule.StringValue) {
		if !strings.ContainsAny(key, ">=<") {
			key += " = "
		}
		b.setValues = append(b.setValues, key+value)
	} else {
		b.err = ErrSet
	}
}

func (b *BuildCore) limit(limitValue int) {
	if b.err != nil {
		return
	}
	if limitValue > 0 {
		b.limitValue = limitValue
	} else {
		debug.Warning("limit can not < 1")
	}
}
func (b *BuildCore) offset(offsetValue int) {
	if b.err != nil {
		return
	}
	if offsetValue > 0 {
		b.offsetValue = offsetValue
	} else {
		debug.Warning("offset can not < 1")
	}
}

func (b *BuildCore) column(column string) {
	if b.err != nil {
		return
	}
	if CheckInjection(column) {
		b.err = ErrInjection
		return
	}
	if column != "" {
		b.columnValues = append(b.columnValues, column)
	} else {
		debug.Println("column is nil")
	}
}

func (b *BuildCore) setValueColumns(ty reflect.Type, tag string) {
	b.insertAutoIndex = -1
	for i := 0; i < ty.NumField(); i++ {
		name := ty.Field(i).Tag.Get(tag)
		if name != "" && b.isOptions(name) && b.isNoOptions(name) {
			columnTags := strings.Split(name, ";")
			var myCat, columnTag string
			var auto bool
			for _, tempTag := range columnTags {
				if len(tempTag) > 6 && strings.HasPrefix(tempTag, "mycat:") {
					myCat = strings.Replace(tempTag, "mycat:", "", 1)
				} else if tempTag == "auto" {
					auto = true
				} else {
					columnTag = tempTag
				}
			}
			b.insertColumns = append(b.insertColumns, columnTag) //追加上与数据库对应的tag
			if auto {
				b.insertAutoIndex = len(b.insertColumns) - 1
				if myCat != "" {
					b.insertMyCatValue = myCat
				}
			}
			b.insertTags = append(b.insertTags, i)
		}
	}
}
func (b *BuildCore) setNoOptions(noOptions []string) {
	b.insertNoOptions = noOptions
}

func (b *BuildCore) isNoOptions(column string) bool {
	for i := 0; i < len(b.insertNoOptions); i++ {
		if b.insertNoOptions[i] == column {
			return false
		}
	}
	return true
}

func (b *BuildCore) setOptions(options []string) {
	b.insertOptions = options
}

func (b *BuildCore) isOptions(column string) bool {
	if len(b.insertOptions) == 0 {
		return true
	}
	for i := 0; i < len(b.insertOptions); i++ {
		if b.insertOptions[i] == column {
			return true
		}
	}
	return false
}

func (b *BuildCore) value(ind reflect.Value, rule Rule, wg ... *sync.WaitGroup) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}
	if b.err != nil {
		return
	}
	var values []string
	for _, v := range b.insertTags {
		value, err := GetValue(ind.Field(v), rule)
		if err != nil {
			b.err = err
			return
		}
		values = append(values, value)
	}
	if len(values) > 0 {
		//判断自增列为default
		if b.insertAutoIndex >= 0 && len(values) > b.insertAutoIndex && values[b.insertAutoIndex] == "DEFAULT" {
			if b.insertMyCatValue == "" {
				values[b.insertAutoIndex] = "NULL"
			} else {
				values[b.insertAutoIndex] = b.insertMyCatValue
			}
		}
		row := strings.Join(values, ",")
		func() {
			b.insertValuesMutex.Lock()
			defer b.insertValuesMutex.Unlock()
			b.insertValues = append(b.insertValues, row)
		}()
	} else {
		debug.Println("Insert no found data")
	}
}
