package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status Status `json:"status"`
}

func init() {
	go AutoReloadJSON()
}

func main() {
	http.HandleFunc("/", AutoReloadWeb)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}

func AutoReloadJSON() {
	for {
		data := Data{}
		data.Status.Water = RandomNumberGenerator(20, 1)
		data.Status.Wind = RandomNumberGenerator(20, 1)

		encodedData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Encrypt error: ", err.Error())
		}

		_ = ioutil.WriteFile("data.json", encodedData, 0644)

		time.Sleep(time.Second * 15)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	var data Data
	waterLevel := data.Status.Water
	windLevel := data.Status.Wind

	// read file json
	dataJSON, errJSON := os.Open("data.json")
	if errJSON != nil {
		fmt.Println("Read file json error: ", errJSON.Error())
	}

	// defer file close
	defer dataJSON.Close()

	byteValue, errByte := ioutil.ReadAll(dataJSON)
	if errByte != nil {
		fmt.Println("Read file json error: ", errByte.Error())
	}

	errUnmarshal := json.Unmarshal(byteValue, &data)
	if errUnmarshal != nil {
		fmt.Println("Unmarshaling error: ", errUnmarshal.Error())
	}

	// logic check status
	waterStatus := StatusChecker("water", waterLevel)
	windStatus := StatusChecker("wind", windLevel)

	var newData = map[string]string{
		"waterLevel":  strconv.Itoa(waterLevel),
		"windLevel":   strconv.Itoa(windLevel),
		"waterStatus": waterStatus,
		"windStatus":  windStatus,
	}

	t, errParse := template.ParseFiles("index.html")
	if errParse != nil {
		fmt.Println("Parsing file html error: ", errParse.Error())
	}

	errExec := t.Execute(w, newData)
	if errExec != nil {
		fmt.Println("Execute error: ", errExec.Error())
	}
}

func RandomNumberGenerator(max int, min int) (level int) {
	level = rand.Intn(max-min) + min
	return
}

func StatusChecker(attribute string, level int) (status string) {
	if attribute == "water" {
		if level < 5 {
			status = "Aman"
		} else if level <= 8 {
			status = "Siaga"
		} else if level > 8 {
			status = "Bahaya"
		}
	} else if attribute == "wind" {
		if level < 6 {
			status = "Aman"
		} else if level <= 15 {
			status = "Siaga"
		} else if level > 15 {
			status = "Bahaya"
		}
	}
	return
}
