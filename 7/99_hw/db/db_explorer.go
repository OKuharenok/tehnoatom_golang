package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type JR map[string]interface{}

type Table struct {
	Name    string
	Id      string //имя ключа в таблице
	Columns []Field
}

type Field struct {
	Name       string
	Type       string
	IsNullable bool
	IsKey      bool
}

func (field *Field)defaultValue() interface{} {
	if field.Type == "varchar" || field.Type == "text" {
		return ""
	} else {
		return 0
	}
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	TabNames, err := getTableName(db)
	if err != nil {
		panic(err)
	}
	tables := make(map[string]Table, len(TabNames))
	for _, name := range TabNames {
		var id string
		rows, err := db.Query("select column_name, data_type, if(is_nullable='YES', true, false) as is_nullable, if(column_key = 'PRI', true, false) as is_key from information_schema.columns where table_schema=database() and table_name = ?", name)
		if err != nil {
			panic(err)
		}
		var res []Field
		for rows.Next() {
			var f Field
			rows.Scan(&f.Name, &f.Type, &f.IsNullable, &f.IsKey)
			if f.IsKey {
				id = f.Name
			}
			res = append(res, f)
		}
		rows.Close()
		var tmp Table
		tmp.Name = name
		tmp.Id = id
		tmp.Columns = res
		tables[name] = tmp
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/" {
				resp, err := json.Marshal(JR{"response": JR{"tables": TabNames}})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(resp)
				return
			}
			str := strings.Split(r.URL.Path, "/")
			if len(str) == 2 {
				if _, isExist := tables[str[1]]; !isExist {
					resp, err := json.Marshal(JR{"error": "unknown table"})
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusNotFound)
					w.Write(resp)
					return
				}
				limit := 5
				offset := 0
				if r.URL.Query().Get("limit") != "" {
					limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
					if err != nil {
						limit = 5
					}
				}
				if r.URL.Query().Get("offset") != "" {
					offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
					if err != nil {
						offset = 0
					}
				}
				records, err := getRecords(db, str[1], tables, limit, offset)
				if err != nil {
					panic(err)
				}
				resp, err := json.Marshal(JR{"response": JR{"records": records}})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(resp)
				return
			}
			if len(str) == 3 {
				if _, isExist := tables[str[1]]; !isExist {
					resp, err := json.Marshal(JR{"error": "unknown table"})
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusNotFound)
					w.Write(resp)
					return
				}
				id, err := strconv.Atoi(str[2])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				record, err := getRow(db, tables[str[1]], id)
				if err != nil {
					resp, err := json.Marshal(JR{"error" : "record not found"})
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusNotFound)
					w.Write(resp)
					return
				}
				resp, err := json.Marshal(JR{"response" : JR{"record" : record}})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(resp)
				return
			}
		case http.MethodPut :
			str := strings.Split(r.URL.Path, "/")
			if _, isExist := tables[str[1]]; !isExist {
				resp, err := json.Marshal(JR{"error": "unknown table"})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(resp)
				return
			}
			var params map[string]interface{}
			body, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(body, &params)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			insertId := insert(db, tables[str[1]], params)
			resp, err := json.Marshal(JR{"response" : JR{tables[str[1]].Id : insertId}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		case http.MethodPost :
			str := strings.Split(r.URL.Path, "/")
			if _, isExist := tables[str[1]]; !isExist {
				resp, err := json.Marshal(JR{"error": "unknown table"})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(resp)
				return
			}
			recordId := str[2]
			var params map[string]interface{}
			body, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(body, &params)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			error := validate(tables[str[1]], params)
			if error != nil {
				resp, err := json.Marshal(JR{"error" : error.Error()})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write(resp)
				return
			}
			updated := update(db, tables[str[1]], params, recordId)
			resp, err := json.Marshal(JR{"response" : JR{"updated" : updated}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		case http.MethodDelete :
			str := strings.Split(r.URL.Path, "/")
			if _, isExist := tables[str[1]]; !isExist {
				resp, err := json.Marshal(JR{"error": "unknown table"})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(resp)
				return
			}
			recordId := str[2]
			id, err := strconv.Atoi(recordId)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			deleted := deleteRecord(db, tables[str[1]], id)
			resp, err := json.Marshal(JR{"response" : JR{"deleted" : deleted}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		}
	})
	return mux, nil
}

func deleteRecord(db *sql.DB, table Table, id int) int64 {
	result, err := db.Exec(fmt.Sprintf("delete from %s where %s = ?", table.Name, table.Id), id)
	if err != nil {
		return 0
	} else {
		res, _ := result.RowsAffected()
		return res
	}
}

func update(db *sql.DB, table Table, params map[string]interface{}, recordId string) int64 {
	id, _ := strconv.Atoi(recordId)
	var values []string
	var mas []interface{}
	for _, field := range table.Columns {
		if val, isExist := params[field.Name]; isExist {
			values = append(values, fmt.Sprintf("%s = ?", field.Name))
			mas = append(mas, val)
		}
	}
	mas = append(mas, id)
	result, err := db.Exec(fmt.Sprintf("update %s set %s where %s = ?", table.Name, strings.Join(values, ", "), table.Id), mas...)
	if err != nil {
		panic(err)
	}
	res, _ := result.RowsAffected()
	return res
}

func validate(table Table, params map[string]interface{}) error {
	for _, field := range table.Columns {
		if value, isExist := params[field.Name]; isExist {
			if value == nil && !field.IsNullable {
				return fmt.Errorf("field %s have invalid type", field.Name)
			}
			if field.Name == table.Id {
				return fmt.Errorf("field %s have invalid type", field.Name)
			}
			switch value.(type) {
			case float64 :
				if field.Type != "int" {
					return fmt.Errorf("field %s have invalid type", field.Name)
				}
			case string :
				if field.Type != "varchar" && field.Type != "text" {
					return fmt.Errorf("field %s have invalid type", field.Name)
				}
			}
		}
	}
	return nil
}

func insert(db *sql.DB, table Table, params map[string]interface{}) int64 {
	mas := make([]interface{}, len(table.Columns))
	names := make([]string, len(table.Columns))
	placeholders := make([]string, len(table.Columns))
	for i, field := range table.Columns {
		names[i] = field.Name
		placeholders[i] = "?"
		if table.Id == field.Name {
			continue
		}
		if _, isExist := params[field.Name]; isExist {
			mas[i] = params[field.Name]
		} else {
			if field.IsNullable {
				mas[i] = nil
			} else {
				mas[i] = field.defaultValue()
			}
		}
	}
	result, err := db.Exec(fmt.Sprintf("insert into %s (%s) values (%s)", table.Name, strings.Join(names, ", "), strings.Join(placeholders, ", ")), mas...)
	if err != nil {
		panic(err)
	}
	id, _ := result.LastInsertId()
	return id
}

func getRow(db *sql.DB, table Table, id int) (JR, error) {
	row := db.QueryRow(fmt.Sprintf("select * from %s where %s = %d", table.Name, table.Id, id))
	record := make([]interface{}, len(table.Columns))
	for i, field := range table.Columns {
		switch field.Type {
		case "varchar" :
			record[i] = new(sql.NullString)
		case "text" :
			record[i] = new(sql.NullString)
		case "int" :
			record[i] = new(sql.NullInt64)
		}
	}
	err := row.Scan(record...)
	if err != nil {
		return nil, err
	}
	item := make(map[string]interface{}, len(record))
	for i, v := range record {
		switch v.(type) {
		case *sql.NullString:
			if value, ok := v.(*sql.NullString); ok {
				if value.Valid {
					item[table.Columns[i].Name] = value.String
				} else {
					item[table.Columns[i].Name] = nil
				}
			}
		case *sql.NullInt64:
			if value, ok := v.(*sql.NullInt64); ok {
				if value.Valid {
					item[table.Columns[i].Name] = value.Int64
				} else {
					item[table.Columns[i].Name] = nil
				}
			}

		}
	}
	return item, nil
}

func getRecords(db *sql.DB, TabName string, tables map[string]Table, limit int, offset int) ([]interface{}, error) {
	rows, err := db.Query(fmt.Sprintf("select * from %s limit %d offset %d", TabName, limit, offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []interface{}
	for rows.Next() {
		row := make([]interface{}, len(tables[TabName].Columns))
		for i, field := range tables[TabName].Columns {
			switch field.Type {
			case "varchar":
				row[i] = new(sql.NullString)
			case "text":
				row[i] = new(sql.NullString)
			case "int":
				row[i] = new(sql.NullInt64)

			}
		}
		rows.Scan(row...)
		item := make(map[string]interface{}, len(row))
		for i, v := range row {
			switch v.(type) {
			case *sql.NullString:
				if value, ok := v.(*sql.NullString); ok {
					if value.Valid {
						item[tables[TabName].Columns[i].Name] = value.String
					} else {
						item[tables[TabName].Columns[i].Name] = nil
					}
				}
			case *sql.NullInt64:
				if value, ok := v.(*sql.NullInt64); ok {
					if value.Valid {
						item[tables[TabName].Columns[i].Name] = value.Int64
					} else {
						item[tables[TabName].Columns[i].Name] = nil
					}
				}

			}
		}
		result = append(result, item)
	}
	return result, nil
}

func getTableName(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	var res []string
	for rows.Next() {
		value := ""
		rows.Scan(&value)
		res = append(res, value)
	}
	rows.Close()
	return res, nil
}
