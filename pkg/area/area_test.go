/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-12-22 16:39:28
 */

package area

import (
	"fmt"
	"testing"
)

// TestCountry_GetProvinces 获取所有的省
func TestCountry_GetProvinces(t *testing.T) {
	provinces := GetProvinces()
	fmt.Println(provinces)
}

func TestCountry_GetCities(t *testing.T) {
	var province Province
	for _, provinceOne := range GetProvinces() {
		if provinceOne.Name == "山东省" {
			province = provinceOne
			break
		}
	}

	var city City
	for _, cityOne := range GetCities(province.Id) {
		if cityOne.Name == "枣庄市" {
			city = cityOne
			break
		}
	}
	fmt.Println(city)
}

func TestCountry_GetCounties(t *testing.T) {
	counties := GetCounties(370400)
	fmt.Println(counties)
}

func TestCountry_Get(t *testing.T) {
	province, city, county := Get(370402)
	fmt.Println(province, city, county)
}
