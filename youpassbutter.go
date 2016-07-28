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
	defaultDataHost = "localhost"
	defaultDataPort = 3306
	defaultPort     = 8080
	defaultError    = "{\"error\":\"Something went wrong\""
)

type serverConfig struct {
	DataHost     string `json:"dataHost"`
	DataPort     int    `json:"dataPort"`
	DataName     string `json:"dataName"`
	DataUsername string `json:"dataUsername"`
	DataPassword string `json:"dataPassword"`
	Port         int    `json:"port"`
}

func assignServerConfigDefaultValues(config serverConfig) serverConfig {
	if config.DataHost == "" {
		config.DataHost = defaultDataHost
	}

	if config.DataPort == 0 {
		config.DataPort = defaultDataPort
	}

	if config.Port == 0 {
		config.Port = defaultPort
	}

	return config
}

func getServerConfig(filename string) (serverConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return serverConfig{}, err
	}

	var config serverConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return serverConfig{}, err
	}

	return assignServerConfigDefaultValues(config), nil
}

func getQueries(filename string) (map[string]string, error) {
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

func getFilenames() (string, string, error) {
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

func getDBConnectionAndQueries() (*sql.DB, serverConfig, map[string]string, error) {
	configFile, queriesFile, err := getFilenames()
	if err != nil {
		return nil, serverConfig{}, nil, err
	}

	config, err := getServerConfig(configFile)
	if err != nil {
		return nil, serverConfig{}, nil, err
	}

	connectionString := fmt.Sprintf("postgres://%v:%v@%v/%v", config.DataUsername, config.DataPassword, config.DataHost, config.DataName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, serverConfig{}, nil, err
	}

	queries, err := getQueries(queriesFile)
	if err != nil {
		return nil, serverConfig{}, nil, err
	}

	return db, config, queries, nil
}

var (
	db, config, queries, startupError = getDBConnectionAndQueries()
)

func getQueryAndParams(r *http.Request) (string, []string, error) {
	queryStringMap := r.URL.Query()
	queries := queryStringMap["q"]
	if len(queries) == 0 {
		return "", nil, errors.New("You must provide a query (q=)")
	}

	query := queries[0]
	params := queryStringMap["p"]
	return query, params, nil
}

type errorJson struct {
	Error string `json:"error"`
}

func logErrorMessage(r *http.Request, errorString string) {
	log.Println("ERROR", r.RemoteAddr, r.RequestURI, errorString)
}

func logInfoMessage(r *http.Request, response interface{}) {
	log.Println("INFO", r.RemoteAddr, r.RequestURI, response)
}

func writeErrorMessage(w http.ResponseWriter, r *http.Request, errorString string) {
	logErrorMessage(r, errorString)

	w.WriteHeader(http.StatusBadRequest)
	errorJson := errorJson{errorString}
	json, err := json.Marshal(errorJson)
	if err != nil {
		fmt.Fprintf(w, defaultError)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		fmt.Fprintf(w, defaultError)
	}
}

func isSelectQuery(query string) bool {
	return strings.ToLower(query[0:6]) == "select"
}

func getTypedInterface(s string) interface{} {
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

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	query, params, err := getQueryAndParams(r)
	if err != nil {
		writeErrorMessage(w, r, err.Error())
		return
	}

	storedQuery, queryIsPresent := queries[query]
	if !queryIsPresent {
		writeErrorMessage(w, r, "Passed query is not present in queries.json")
		return
	}

	parameterCount := strings.Count(storedQuery, "$")
	if parameterCount != len(params) {
		writeErrorMessage(w, r, "Count of passed parameters are unequal to expected parameters")
		return
	}

	paramsInterfaced := make([]interface{}, len(params))
	for i, v := range params {
		paramsInterfaced[i] = v
	}

	if isSelectQuery(storedQuery) {
		rows, err := db.Query(storedQuery, paramsInterfaced...)
		if err != nil {
			writeErrorMessage(w, r, err.Error())
			return
		}

		columns, err := rows.Columns()
		if err != nil {
			writeErrorMessage(w, r, err.Error())
			return
		}

		rawResult := make([][]byte, len(columns))
		result := make([]*interface{}, len(columns))

		destination := make([]interface{}, len(columns))
		for i := range rawResult {
			destination[i] = &rawResult[i]
		}

		response := []map[string]*interface{}{}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(destination...)
			if err != nil {
				writeErrorMessage(w, r, err.Error())
				return
			}

			for i, raw := range rawResult {
				if raw == nil {
					result[i] = nil
				} else {
					temp := getTypedInterface(string(raw))
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
			writeErrorMessage(w, r, err.Error())
			return
		}

		logInfoMessage(r, response)
		_, err = w.Write(js)
		if err != nil {
			writeErrorMessage(w, r, err.Error())
			return
		}

	} else {
		result, err := db.Exec(storedQuery, paramsInterfaced...)
		if err != nil {
			writeErrorMessage(w, r, err.Error())
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
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		log.Fatalln(startupError)
	}
}
