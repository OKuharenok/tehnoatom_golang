package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var FileName string = "dataset.xml"
var Token string = "123"

type Person struct {
	Id        int    `xml:"id" json:"Id"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	Name      string `json:"Name"`
	Age       int    `xml:"age" json:"Age"`
	Gender    string `xml:"gender" json:"Gender"`
	About     string `xml:"about" json:"About"`
}
type root struct {
	List []Person `xml:"row"`
}

func sorting(pr []Person, orderField string, orderBy int) {
	if orderBy == 0 {
		return
	}
	if orderField == "id" {
		sort.SliceStable(pr, func(i, j int) bool {
			if orderBy == OrderByDesc {
				return pr[i].Id > pr[j].Id
			} else {
				return pr[i].Id < pr[j].Id
			}
		})
	}
	if orderField == "age" {
		sort.SliceStable(pr, func(i, j int) bool {
			if orderBy == OrderByDesc {
				return pr[i].Age > pr[j].Age
			} else {
				return pr[i].Age < pr[j].Age
			}
		})
	}
	if orderField == "name" {
		sort.SliceStable(pr, func(i, j int) bool {
			if orderBy == OrderByDesc {
				return pr[i].Name > pr[j].Name
			} else {
				return pr[i].Name < pr[j].Name
			}
		})
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	prs := new(root)
	xmlFile, err := ioutil.ReadFile(FileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	xml.Unmarshal(xmlFile, prs)
	AccessToken := r.Header.Get("AccessToken")
	if AccessToken != Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 1 {
		u := ""
		data, _ := json.Marshal(&u)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(data))
		return
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	if orderField == "" {
		orderField = "name"
	}
	if orderField != "age" && orderField != "id" && orderField != "name" {
		u := SearchErrorResponse{"ErrorBadOrderField"}
		data, _ := json.Marshal(&u)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(data))
		return
	}
	orderBy, _ := strconv.Atoi(r.URL.Query().Get("order_by"))
	if orderBy != OrderByAsIs && orderBy != OrderByAsc && orderBy != OrderByDesc {
		u := SearchErrorResponse{"Bad OrderBy"}
		data, _ := json.Marshal(&u)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(data))
		return
	}
	result := []Person{}
	for _, value := range prs.List {
		if strings.Contains(value.FirstName+value.LastName, query) || strings.Contains(value.About, query) {
			value.Name = value.FirstName + " " + value.LastName
			result = append(result, value)
		}
	}
	sorting(result, orderField, orderBy)
	if offset > len(result) {
		data := ""
		jStr, _ := json.Marshal(&data)
		fmt.Fprint(w, string(jStr))
		return
	}
	result = result[offset:]
	if len(result) >= limit {
		result = result[:limit]
	}
	js, _ := json.Marshal(&result)
	fmt.Fprint(w, string(js))
}
