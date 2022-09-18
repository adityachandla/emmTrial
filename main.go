package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var tagToFieldNameMapping map[string]string

type houseInfo struct {
	Suburb         string  `csv:"Suburb"`
	Rooms          int     `csv:"Rooms"`
	HouseType      string  `csv:"Type"`
	Price          int     `csv:"Price"`
	Method         string  `csv:"Method"`
	SellerName     string  `csv:"SellerG"`
	DistFromCenter float64 `csv:"Distance"`
	Bedrooms       int     `csv:"Bedroom2"`
	Bathrooms      int     `csv:"Bathroom"`
	Car            int     `csv:"Car"`
	LandSize       float64 `csv:"Landsize"`
	BuildingArea   float64 `csv:"BuildingArea"`
	YearBuilt      int     `csv:"YearBuilt"`
	CouncilArea    string  `csv:"CouncilArea"`
	Region         string  `csv:"Regionname"`
	PropertyCount  string  `csv:"Propertycount"`
}

func main() {
	makeMapping()
	houses := readHouses()
	fmt.Println(len(houses))
}

func makeMapping() {
	tagToFieldNameMapping = make(map[string]string)
	fields := reflect.VisibleFields(reflect.TypeOf(houseInfo{}))
	for _, field := range fields {
		csvName := field.Tag.Get("csv")
		tagToFieldNameMapping[csvName] = field.Name
	}
}

func readHouses() []*houseInfo {
	file, _ := os.Open("melb_data.csv")
	reader := bufio.NewReader(file)

	line, _, _ := reader.ReadLine()
	headers := strings.Split(string(line), ",")

	houses := make([]*houseInfo, 0, 100)

	line, _, err := reader.ReadLine()
	for err == nil {
		entry := strings.Split(string(line), ",")
		houseInfo := parseInfo(entry, headers)
		houses = append(houses, houseInfo)
		line, _, err = reader.ReadLine()
	}
	return houses
}

func parseInfo(entry []string, headers []string) *houseInfo {
	if len(entry) != len(headers) {
		fmt.Println(entry)
		panic("Invalid row encountered")
	}

	info := &houseInfo{}
	for idx, header := range headers {
		name, exists := tagToFieldNameMapping[header]
		if exists && entry[idx] != "" {
			fieldVal := reflect.ValueOf(info).Elem().FieldByName(name)
			kind := fieldVal.Kind().String()
			if kind == "string" {
				fieldVal.SetString(entry[idx])
			} else if kind == "int" {
				val, err := strconv.ParseFloat(entry[idx], 64)
				check(err)
				fieldVal.SetInt(int64(val))
			} else if kind == "float64" {
				val, err := strconv.ParseFloat(entry[idx], 64)
				check(err)
				fieldVal.SetFloat(val)
			}
		}
	}
	return info
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
