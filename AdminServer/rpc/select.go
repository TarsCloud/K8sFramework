package rpc

import (
	"database/sql"
	"encoding/json"
	"github.com/elgris/sqrl"
	"strconv"
	"strings"
)

type SqlColumn struct {
	ColumnName string
	ColumnType string
}

type RequestColumnSqlColumnMap map[string]SqlColumn

func replaceRequestColumn(requestColumnSqlColumnMap RequestColumnSqlColumnMap, requestColumns SelectRequestColumns) ([]string, []string) {
	resultColumns := make([]string, 0, len(requestColumns))
	selectColumns := make([]string, 0, len(requestColumns))

	for _, v := range requestColumns {
		if sqlColumn, ok := requestColumnSqlColumnMap[v]; ok == true {
			resultColumns = append(resultColumns, v)
			selectColumns = append(selectColumns, sqlColumn.ColumnName)
		}
	}
	return resultColumns, selectColumns
}

func replaceFilterColumns(requestColumnSqlColumnMap RequestColumnSqlColumnMap, filter *SelectRequestFilter) SelectRequestFilter {

	selectFilter := SelectRequestFilter{}

	replaceMapInterface := func(source map[string]interface{}, destination *map[string]interface{}) {
		if source == nil {
			return
		}

		if *destination == nil {
			*destination = make(map[string]interface{})
		}

		for k, v := range source {
			if v == nil {
				continue
			}
			replaceKey, ok := (requestColumnSqlColumnMap)[k]
			if ok {
				(*destination)[replaceKey.ColumnName] = v
			}
		}
	}

	replaceMapString := func(source map[string]string, destination *map[string]string) {
		if source == nil {
			return
		}

		if *destination == nil {
			*destination = make(map[string]string)
		}

		for k, v := range source {
			if v == "" {
				continue
			}
			replaceKey, ok := (requestColumnSqlColumnMap)[k]
			if ok {
				(*destination)[replaceKey.ColumnName] = v
			}
		}
	}

	replaceMapInterface(filter.EQ, &selectFilter.EQ)
	replaceMapInterface(filter.NE, &selectFilter.NE)
	replaceMapInterface(filter.GT, &selectFilter.GT)
	replaceMapInterface(filter.GE, &selectFilter.GE)
	replaceMapInterface(filter.LT, &selectFilter.LT)
	replaceMapInterface(filter.LE, &selectFilter.LE)
	replaceMapInterface(filter.IN, &selectFilter.IN)
	replaceMapString(filter.LIKE, &selectFilter.LIKE)

	return selectFilter
}

func replaceOrderColumns(requestColumnSqlColumnMap RequestColumnSqlColumnMap, order *SelectRequestOrder) SelectRequestOrder {
	selectOrder := SelectRequestOrder{}
	for i := range *order {
		if sqlColumn, ok := requestColumnSqlColumnMap[(*order)[i].Column]; ok {
			if !strings.EqualFold((*order)[i].Order, "desc") {
				selectOrder = append(selectOrder, SelectRequestOrderElem{Column: sqlColumn.ColumnName})
			} else {
				selectOrder = append(selectOrder, SelectRequestOrderElem{Column: sqlColumn.ColumnName, Order: "DESC"})
			}
		}
	}
	return selectOrder
}

func selectSQLBuilderAppendFilter(filter *SelectRequestFilter, builder *sqrl.SelectBuilder) {
	builder.Where("1=1")
	builder.Where(sqrl.Eq(filter.EQ))
	builder.Where(sqrl.NotEq(filter.NE))
	builder.Where(sqrl.Gt(filter.GT))
	builder.Where(sqrl.GtOrEq(filter.GE))
	builder.Where(sqrl.Lt(filter.LT))
	builder.Where(sqrl.LtOrEq(filter.LE))

	if filter.IN != nil {
		for k, v := range filter.IN {
			builder.Where(sqrl.Eq{k: v})
		}
	}

	if filter.LIKE != nil {
		for k, v := range filter.LIKE {
			builder.Where(k + " like " + "\"" + v + "\"")
		}
	}
}

func selectSQLBuilderAppendOrder(order SelectRequestOrder, builder *sqrl.SelectBuilder) {
	for i := range order {
		if strings.EqualFold(order[i].Order, "DESC") {
			builder.OrderBy(order[i].Column + " DESC")
		} else {
			builder.OrderBy(order[i].Column)
		}
	}
}

func selectSQLBuilderAppendLimiter(limiter *SelectRequestLimiter, builder *sqrl.SelectBuilder) {
	if limiter.Offset != 0 {
		builder.Offset(limiter.Offset)
	}
	if limiter.Rows != 0 {
		builder.Limit(limiter.Rows)
	}
}

func execFilterCount(db sqrl.BaseRunner, from string, filter *SelectRequestFilter, fixedFilter *SelectRequestFilter) (string, int64, error) {
	const countName = "FilterCount"
	selectBuilder := sqrl.Select("count(*)").From(from)
	if filter != nil {
		selectSQLBuilderAppendFilter(filter, selectBuilder)
	}
	if fixedFilter != nil {
		selectSQLBuilderAppendFilter(fixedFilter, selectBuilder)
	}
	row := selectBuilder.RunWith(db).QueryRow()
	var countValue int64
	if err := row.Scan(&countValue); err != nil {
		return countName, -1, err
	}
	return countName, countValue, nil
}

func execAllCount(db sqrl.BaseRunner, from string, fixedFilter *SelectRequestFilter) (string, int64, error) {

	const countName = "AllCount"

	selectBuilder := sqrl.Select("count(*)").From(from)

	if fixedFilter != nil {
		selectSQLBuilderAppendFilter(fixedFilter, selectBuilder)
	}

	row := selectBuilder.RunWith(db).QueryRow()

	var countValue int64
	if err := row.Scan(&countValue); err != nil {
		return countName, -1, err
	}
	return countName, countValue, nil
}

func execSelectSql(db sqrl.BaseRunner, from string, requestParams *RequestParams, requestColumnSqlColumnMap RequestColumnSqlColumnMap, fixedFilter *SelectRequestFilter) ([]map[string]interface{}, map[string]int64, error) {
	var err error
	var rows *sql.Rows
	var data []map[string]interface{}
	var count map[string]int64

	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	var selectFilter SelectRequestFilter
	if requestParams.Filter != nil {
		selectFilter = replaceFilterColumns(requestColumnSqlColumnMap, requestParams.Filter)
	}

	for {

		if requestParams.Columns == nil {
			break
		}

		resultColumns, selectColumns := replaceRequestColumn(requestColumnSqlColumnMap, *requestParams.Columns)
		if len(resultColumns) == 0 {
			data = make([]map[string]interface{}, 0)
			break
		}

		selectBuilder := sqrl.Select(selectColumns...).From(from)
		selectSQLBuilderAppendFilter(&selectFilter, selectBuilder)

		if fixedFilter != nil {
			selectSQLBuilderAppendFilter(fixedFilter, selectBuilder)
		}

		if requestParams.Order != nil {
			selectOrder := replaceOrderColumns(requestColumnSqlColumnMap, requestParams.Order)
			selectSQLBuilderAppendOrder(selectOrder, selectBuilder)
		}

		if requestParams.Limiter != nil {
			selectSQLBuilderAppendLimiter(requestParams.Limiter, selectBuilder)
		} else {
			const FixedSelectOffset = 0
			const FixedSelectLimit = 50
			selectBuilder.Offset(FixedSelectOffset).Limit(FixedSelectLimit)
		}

		if rows, err = selectBuilder.RunWith(db).Query(); err != nil {
			break
		}

		columns := make([]interface{}, len(selectColumns))
		columnPointers := make([]interface{}, len(selectColumns))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		data = make([]map[string]interface{}, 0, 30)
		for rows.Next() {
			m := make(map[string]interface{}, len(selectColumns))
			if err := rows.Scan(columnPointers...); err != nil {
				break
			}
			for i, columnName := range resultColumns {
				value := columns[i]
				dataType := requestColumnSqlColumnMap[columnName].ColumnType
				switch dataType {
				case "json":
					if value == nil {
						m[columnName] = nil
					} else {
						m[columnName] = json.RawMessage(value.([]byte))
					}
				case "string":
					if value == nil {
						m[columnName] = ""
					} else {
						m[columnName] = string(value.([]byte))
					}
				case "bool":
					if value == nil {
						m[columnName] = false
					} else {
						switch value.(type) {
						case []uint8:
							v := value.([]byte)[0]
							m[columnName] = v == '1'
						case int64:
							v := value.(int64)
							m[columnName] = v != 0
						}
					}
				case "int":
					if value == nil {
						m[columnName] = nil
					} else {
						switch value.(type) {
						case int, int64, int32:
							m[columnName] = value
						case []uint8:
							vStr := string(value.([]byte))
							vInt, _ := strconv.Atoi(vStr)
							m[columnName] = vInt
						}
					}

				default:
					m[columnName] = string(value.([]byte))
				}
			}
			data = append(data, m)
		}

		break
	}

	if err != nil {
		return nil, nil, err
	}

	if requestParams.Count == nil {
		return data, nil, err
	}

	count = make(map[string]int64, len(*requestParams.Count))

	var countName string
	var countValue int64

	for _, v := range *requestParams.Count {
		switch v {
		case "AllCount":
			if countName, countValue, err = execAllCount(db, from, fixedFilter); err != nil {
				return nil, nil, err
			}
			count[countName] = countValue
		case "FilterCount":
			if countName, countValue, err = execFilterCount(db, from, &selectFilter, fixedFilter); err != nil {
				return nil, nil, err
			}
			count[countName] = countValue
		}
	}

	return data, count, nil
}
