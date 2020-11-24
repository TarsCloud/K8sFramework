package mysql

import (
	"database/sql"
	"encoding/json"
	"github.com/elgris/sqrl"
	"strconv"
	"strings"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
)


type SqlColumn struct {
	ColumnName string
	ColumnType string
}

type RequestColumnSqlColumnMap map[string]SqlColumn

var TafDb *sql.DB

func SelectQueryResult(from string, Filter, Limiter, Order *string, SqlColumnMap RequestColumnSqlColumnMap) (*models.SelectResult, error) {
	selectParams, err := util.ParseSelectQuery(Filter, Limiter, Order)
	if err != nil {
		return nil, err
	}

	var result *models.SelectResult
	if result, err = execSelectSql(TafDb, from, selectParams, SqlColumnMap, nil); err != nil {
		return nil, err
	}

	return result, nil
}

func replaceRequestColumn(requestColumnSqlColumnMap RequestColumnSqlColumnMap) ([]string, []string) {
	resultColumns := make([]string, 0, len(requestColumnSqlColumnMap))
	selectColumns := make([]string, 0, len(requestColumnSqlColumnMap))

	for k, sqlColumn := range requestColumnSqlColumnMap {
		resultColumns = append(resultColumns, k)
		selectColumns = append(selectColumns, sqlColumn.ColumnName)
	}
	return resultColumns, selectColumns
}

func replaceFilterColumns(requestColumnSqlColumnMap RequestColumnSqlColumnMap, filter *models.SelectRequestFilter) models.SelectRequestFilter {

	selectFilter := models.SelectRequestFilter{}

	replaceMapInterface := func(source models.MapInterface, destination *models.MapInterface) {
		if source == nil {
			return
		}

		if *destination == nil {
			*destination = make(models.MapInterface)
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

	replaceMapString := func(source models.MapString, destination *models.MapString) {
		if source == nil {
			return
		}

		if *destination == nil {
			*destination = make(models.MapString)
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

	replaceMapInterface(filter.Eq, &selectFilter.Eq)
	replaceMapInterface(filter.Ne, &selectFilter.Ne)
	replaceMapString(filter.Like, &selectFilter.Like)

	return selectFilter
}

func replaceOrderColumns(requestColumnSqlColumnMap RequestColumnSqlColumnMap, order *models.SelectRequestOrder) models.SelectRequestOrder {
	selectOrder := models.SelectRequestOrder{}
	for i := range *order {
		if sqlColumn, ok := requestColumnSqlColumnMap[(*order)[i].Column]; ok {
			if !strings.EqualFold((*order)[i].Order, "desc") {
				selectOrder = append(selectOrder, &models.SelectRequestOrderElem{Column: sqlColumn.ColumnName})
			} else {
				selectOrder = append(selectOrder, &models.SelectRequestOrderElem{Column: sqlColumn.ColumnName, Order: "DESC"})
			}
		}
	}
	return selectOrder
}

func selectSQLBuilderAppendFilter(filter *models.SelectRequestFilter, builder *sqrl.SelectBuilder) {
	builder.Where("1=1")
	builder.Where(sqrl.Eq(filter.Eq))
	builder.Where(sqrl.NotEq(filter.Ne))

	if filter.Like != nil {
		for k, v := range filter.Like {
			builder.Where(k + " like " + "\"" + strings.ReplaceAll(v, ".*", "%") + "\"")
		}
	}
}

func selectSQLBuilderAppendOrder(order models.SelectRequestOrder, builder *sqrl.SelectBuilder) {
	for i := range order {
		if strings.EqualFold(order[i].Order, "DESC") {
			builder.OrderBy(order[i].Column + " DESC")
		} else {
			builder.OrderBy(order[i].Column)
		}
	}
}

func selectSQLBuilderAppendLimiter(limiter *models.SelectRequestLimiter, builder *sqrl.SelectBuilder) {
	if *limiter.Offset != 0 {
		builder.Offset(uint64(*limiter.Offset))
	}
	if limiter.Rows != 0 {
		builder.Limit(uint64(limiter.Rows))
	}
}

func execFilterCount(db sqrl.BaseRunner, from string, filter *models.SelectRequestFilter, fixedFilter *models.SelectRequestFilter) (string, int32, error) {
	const countName = "FilterCount"
	selectBuilder := sqrl.Select("count(*)").From(from)
	if filter != nil {
		selectSQLBuilderAppendFilter(filter, selectBuilder)
	}
	if fixedFilter != nil {
		selectSQLBuilderAppendFilter(fixedFilter, selectBuilder)
	}
	row := selectBuilder.RunWith(db).QueryRow()
	var countValue int32
	if err := row.Scan(&countValue); err != nil {
		return countName, -1, err
	}
	return countName, countValue, nil
}

func execAllCount(db sqrl.BaseRunner, from string, fixedFilter *models.SelectRequestFilter) (string, int32, error) {

	const countName = "AllCount"

	selectBuilder := sqrl.Select("count(*)").From(from)

	if fixedFilter != nil {
		selectSQLBuilderAppendFilter(fixedFilter, selectBuilder)
	}

	row := selectBuilder.RunWith(db).QueryRow()

	var countValue int32
	if err := row.Scan(&countValue); err != nil {
		return countName, -1, err
	}
	return countName, countValue, nil
}

func execSelectSql(db sqrl.BaseRunner, from string, requestParams *util.SelectParams, requestColumnSqlColumnMap RequestColumnSqlColumnMap, fixedFilter *models.SelectRequestFilter) (*models.SelectResult, error) {
	var err error
	var rows *sql.Rows
	var data models.ArrayMapInterface

	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	var selectFilter models.SelectRequestFilter
	if requestParams.Filter != nil {
		selectFilter = replaceFilterColumns(requestColumnSqlColumnMap, requestParams.Filter)
	}

	for {
		resultColumns, selectColumns := replaceRequestColumn(requestColumnSqlColumnMap)
		if len(resultColumns) == 0 {
			data = make(models.ArrayMapInterface, 0)
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

		data = make(models.ArrayMapInterface, 0, 30)
		for rows.Next() {
			m := make(models.MapInterface, len(selectColumns))
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
		return nil, err
	}

	count := make(models.MapInt)

	var allCount, filterCount int32
	if _, allCount, err = execAllCount(db, from, fixedFilter); err != nil {
		return nil, err
	}
	count["AllCount"] = allCount

	if _, filterCount, err = execFilterCount(db, from, &selectFilter, fixedFilter); err != nil {
		return nil, err
	}
	count["FilterCount"] = filterCount

	return &models.SelectResult{Data: data, Count: count}, nil
}
