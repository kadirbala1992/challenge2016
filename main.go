package main

import (
	distributor "./distributor"
	//"bufio"
	//"encoding/json"
	"fmt"
	//"os"
	//"os/exec"
	//"strconv"
	//"strings"
)

var cities = distributor.PrepareCitiesJson()
var distributorMap = make(map[string]interface{})

//var allCountries []string
var directUserList []string
var indirectUserList []string

func main() {

	for {
		var distType string
		if len(directUserList) == 0 {
			fmt.Printf("By default you need to create a Direct Distributor initially\n")
			distType = "direct"
		} else {
			distType = distributor.GetDistType()
		}

		if distType == "direct" {
			permission := distributor.GetInput()
			//permission := []string{"ILAYA", "INCLUDE: INDIA", "INCLUDE: UNITEDSTATES", "EXCLUDE: KARNATAKA-INDIA", "EXCLUDE: CHENNAI-TAMILNADU-INDIA"}
			valid := distributor.ExistInArray(directUserList, permission[0])
			if valid == "" {
				prepareDirectUser(permission)
			} else {
				fmt.Printf("Direct user already exist with this name\n")
			}

		} else {
			valid := distributor.ExistInArray(directUserList, distType)
			if valid != "" {
				permission := distributor.GetInput()
				//permission := []string{"KADIR", "INCLUDE: KERALA-INDIA", "INCLUDE: PUNJAB-INDIA", "EXCLUDE: GUJARAT-INDIA"}

				prepareInDirectUser(permission, distributorMap[valid].(map[string]interface{}), valid)
			}
		}
	}

}

/*CRITICAL: SALVA,BN,RO,Salva,Bistrita-Nasaud,Romania*/

func prepareDirectUser(permission []string) {

	currentUser := distributor.PrepareRoorUser(permission, cities)
	if currentUser["err"] == nil {
		directUserList = append(directUserList, permission[0])
		currentUser["type"] = "direct"
		distributorMap[permission[0]] = currentUser
		fmt.Printf("%v", distributorMap[permission[0]])
	} else {
		fmt.Printf("%v", currentUser["err"])
	}

	//fmt.Printf("%v", distributorMap[permission[0]])
}

func prepareInDirectUser(permission []string, root map[string]interface{}, parent string) {
	currentUser := distributor.PrepareSubUser(permission, cities, root)
	if currentUser["err"] == nil {
		currentUser["type"] = "indirect"
		currentUser["parent"] = parent
		indirectUserList = append(indirectUserList, permission[0])
		distributorMap[permission[0]] = currentUser
		fmt.Printf("%v", distributorMap[permission[0]])
	} else {
		fmt.Printf("%v", currentUser["err"])
	}

	//fmt.Printf("%v", distributorMap[permission[0]])
}
