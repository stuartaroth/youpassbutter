package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	DATA_HOST = "localhost"
	DATA_PORT = 3306
	PORT      = 8080
)

type Config struct {
	DataHost     string
	DataPort     int
	DataName     string
	DataUsername string
	DataPassword string
	Port         int
}

func AssignConfigDefaultValues(config Config) Config {
	if config.DataHost == "" {
		config.DataHost = DATA_HOST
	}

	if config.DataPort == 0 {
		config.DataPort = DATA_PORT
	}

	if config.Port == 0 {
		config.Port = PORT
	}

	return config
}

func GetConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return AssignConfigDefaultValues(config), nil
}

func GetQueries(filename string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var queries map[string]string
	err = json.Unmarshal(data, &queries)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

func GetFilenames() (string, string, error) {
	var configFile string
	var queriesFile string
	flag.StringVar(&configFile, "c", "", "config.json file")
	flag.StringVar(&queriesFile, "q", "", "queries.json file")
	flag.Parse()
	if configFile == "" || queriesFile == "" {
		return configFile, queriesFile, errors.New("You must specify a config.json and a queries.json file using -c and -q")
	}

	return configFile, queriesFile, nil
}

func GetDBConnectionAndQueries() (*sql.DB, Config, map[string]string, error) {
	configFile, queriesFile, err := GetFilenames()
	if err != nil {
		return nil, Config{}, nil, err
	}

	config, err := GetConfig(configFile)
	if err != nil {
		return nil, Config{}, nil, err
	}

	connectionString := fmt.Sprintf("postgres://%v:%v@%v/%v", config.DataUsername, config.DataPassword, config.DataHost, config.DataName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, Config{}, nil, err
	}

	queries, err := GetQueries(queriesFile)
	if err != nil {
		return nil, Config{}, nil, err
	}

	return db, config, queries, nil
}

var (
	db, config, queries, startupError = GetDBConnectionAndQueries()
)

func GetQueryAndParams(r *http.Request) (string, []string, error) {
	queryStringMap := r.URL.Query()
	queries := queryStringMap["q"]
	if len(queries) == 0 {
		return "", nil, errors.New("You must provide a query (q=)")
	}

	query := queries[0]
	params := queryStringMap["p"]
	return query, params, nil
}

type ErrorJson struct {
	Error string `json:"error"`
}

func LogErrorMessage(r *http.Request, errorString string) {
	log.Println("ERROR", r.RemoteAddr, r.RequestURI, errorString)
}

func LogInfoMessage(r *http.Request, response interface{}) {
	log.Println("INFO", r.RemoteAddr, r.RequestURI, response)
}

func WriteErrorMessage(w http.ResponseWriter, r *http.Request, errorString string) {
	LogErrorMessage(r, errorString)

	w.WriteHeader(http.StatusBadRequest)
	errorJson := ErrorJson{errorString}
	json, err := json.Marshal(errorJson)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"Something went wrong\"")
		return
	}

	w.Write(json)
}

func IsSelectQuery(query string) bool {
	return strings.ToLower(query[0:6]) == "select"
}

func GetTypedInterface(s string) interface{} {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	b, err := strconv.ParseBool(s)
	if err == nil {
		return b
	}

	return s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	query, params, err := GetQueryAndParams(r)
	if err != nil {
		WriteErrorMessage(w, r, err.Error())
		return
	}

	storedQuery, queryIsPresent := queries[query]
	if !queryIsPresent {
		WriteErrorMessage(w, r, "Passed query is not present in queries.json")
		return
	}

	parameterCount := strings.Count(storedQuery, "$")
	if parameterCount != len(params) {
		WriteErrorMessage(w, r, "Count of passed parameters are unequal to expected parameters")
		return
	}

	paramsInterfaced := make([]interface{}, len(params))
	for i, v := range params {
		paramsInterfaced[i] = v
	}

	if IsSelectQuery(storedQuery) {
		rows, err := db.Query(storedQuery, paramsInterfaced...)
		if err != nil {
			WriteErrorMessage(w, r, err.Error())
			return
		}

		columns, err := rows.Columns()
		if err != nil {
			WriteErrorMessage(w, r, err.Error())
			return
		}

		rawResult := make([][]byte, len(columns))
		result := make([]*interface{}, len(columns))

		destination := make([]interface{}, len(columns))
		for i, _ := range rawResult {
			destination[i] = &rawResult[i]
		}

		response := []map[string]*interface{}{}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(destination...)
			if err != nil {
				WriteErrorMessage(w, r, err.Error())
				return
			}

			for i, raw := range rawResult {
				if raw == nil {
					result[i] = nil
				} else {
					temp := GetTypedInterface(string(raw))
					result[i] = &temp
				}
			}
			row := make(map[string]*interface{})
			for i, v := range result {
				row[columns[i]] = v
			}
			response = append(response, row)
		}

		js, err := json.Marshal(response)
		if err != nil {
			WriteErrorMessage(w, r, err.Error())
			return
		}

		LogInfoMessage(r, response)
		w.Write(js)
	} else {
		result, err := db.Exec(storedQuery, paramsInterfaced...)
		if err != nil {
			WriteErrorMessage(w, r, err.Error())
			return
		}

		log.Println(r, result)
		fmt.Fprintf(w, "{\"success\":\"Query executed without errors\"}")
	}
}

func main() {
	if startupError != nil {
		log.Fatalln(startupError)
		return
	}

	listenPort := fmt.Sprintf(":%d", config.Port)
	http.HandleFunc("/", Handler)
	http.ListenAndServe(listenPort, nil)
}
