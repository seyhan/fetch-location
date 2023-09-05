package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)
//test
var cityUrl = ""
var countyUrl = ""
var districtUrl = ""

type Response []struct {
	PARENTID int    `json:"PARENTID"`
	ID       int    `json:"ID"`
	ADI      string `json:"ADI"`
}

func main() {

	var cities = getCities(0)
	cities.saveToFile("cities")

	var totalCounties Response
	var counties Response
	for _, city := range cities {
		counties = getCounties(city.ID)
		for _, county := range counties {
			totalCounties = append(totalCounties, county)
		}
	}
	totalCounties.saveToFile("counties")

	var totalDistricts Response
	var districts Response
	for _, city := range cities {
		var counties = getCounties(city.ID)
		for _, county := range counties {
			districts = getDistricts(county.ID)
			for _, district := range districts {
				totalDistricts = append(totalDistricts, district)
			}
		}
	}
	totalDistricts.saveToFile("districts")

}

func getCities(parentId int) Response {
	resp, err := http.PostForm(cityUrl, url.Values{})

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var cities Response
	if err := json.Unmarshal(body, &cities); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	for i := 0; i < len(cities); i++ {
		cities[i].PARENTID = parentId
	}

	return cities
}

func getCounties(cityId int) Response {
	resp, err := http.PostForm(countyUrl, url.Values{"cityId": {strconv.Itoa(cityId)}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var counties Response
	if err := json.Unmarshal(body, &counties); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	for i := 0; i < len(counties); i++ {
		counties[i].PARENTID = cityId
	}

	return counties
}

func getDistricts(countyId int) Response {
	resp, err := http.PostForm(districtUrl, url.Values{"countyId": {strconv.Itoa(countyId)}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var districts Response
	if err := json.Unmarshal(body, &districts); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	for i := 0; i < len(districts); i++ {
		districts[i].PARENTID = countyId
	}

	return districts
}

func (d Response) saveToFile(fileName string) {

	err := ioutil.WriteFile(fileName, []byte(d.toString()), 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func (d Response) toString() string {

	var s []string
	var c string
	for _, v := range d {
		c = strconv.Itoa(v.PARENTID) + "," + strconv.Itoa(v.ID) + "," + v.ADI
		s = append(s, c)
	}
	return strings.Join(s, "\n")
}
