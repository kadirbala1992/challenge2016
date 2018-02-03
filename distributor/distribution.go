package distribution

import (
	"bufio"
	"encoding/csv"
	//"encoding/json"
	"fmt"
	"io"
	"strings"
	//"io/ioutil"
	"log"
	"os"
)

type Cities struct {
	AreaCode    string `json:"area_code"`
	StateCode   string `json:"state_code"`
	CountryCode string `json:"country_code"`
	Area        string `json:"area"`
	State       string `json:"state"`
	Country     string `json:"country`
}

func PrepareCitiesJson() []Cities {
	csvFile, _ := os.Open("cities.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var cities []Cities
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		cities = append(cities, Cities{
			AreaCode:    line[0],
			StateCode:   line[1],
			CountryCode: line[2],
			Area:        line[3],
			State:       line[4],
			Country:     line[5],
		})
	}
	return cities
}

func getCountryNames(cities []Cities) []string {
	var tmp []string
	for _, item := range cities {
		tmp = append(tmp, item.Country)
	}

	return removeDuplicates(tmp)
}

func excludeAreaName(countries []string, cities []Cities, excludeState []string) []string {
	var tmp []string
	//var needful []map[string]interface{}
	var needful []string

	for _, country := range countries {
		for _, item := range cities {
			if item.Country == country {
				//fmt.Printf("%v\n", fmt.Sprintf("%v - %v - %v", item.Area, item.State, item.Country))
				needful = append(needful, fmt.Sprintf("%v _ %v  %v", item.Area, item.State, item.Country))
			}
		}
	}

	if len(excludeState) > 0 {
		for _, need := range needful {
			for _, item := range excludeState {
				state := strings.Split(need, " _ ")
				if state[1] != item {
					tmp = append(tmp, need)
				}
			}
		}
	} else {
		return removeDuplicates(needful)
	}

	return removeDuplicates(tmp)
}

func excludeStateName(countries []string, cities []Cities) []string {
	var tmp []string
	for _, country := range countries {
		for _, item := range cities {
			if item.Country == country {
				tmp = append(tmp, fmt.Sprintf("%v _ %v", item.State, item.Country))
			}
		}
	}

	return removeDuplicates(tmp)
}

func extendStateName(cities []Cities, countries []string) []string {
	var tmp []string

	for _, country := range countries {
		for _, item := range cities {
			if item.Country != country {
				tmp = append(tmp, fmt.Sprintf("%v _ %v", item.State, item.Country))
			}
		}
	}

	return removeDuplicates(tmp)
}

func extendAreaName(cities []Cities, countries []string, states []string) []string {
	var tmp []string
	if len(countries) > 0 && len(states) > 0 {
		for _, country := range countries {
			for _, item := range cities {
				if item.Country != country {
					for _, state := range states {
						comp := fmt.Sprintf("%v _ %v", item.State, item.Country)
						if comp == "Chennai _ Tamil Nadu _ India" {
							fmt.Printf("State - %v\n", state)
							fmt.Printf("Country - %v\n", item.Country)
						}
						if comp != state {
							tmp = append(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
						}
					}
				}
			}
		}
	} else if len(countries) > 0 {
		for _, item := range cities {
			for _, country := range countries {
				if item.Country != country {
					tmp = append(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
				}
			}
		}
	} else if len(states) > 0 {
		for _, item := range cities {
			for _, state := range states {
				comp := fmt.Sprintf("%v _ %v", item.State, item.Country)
				if comp != state {
					tmp = append(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
				}
			}
		}
	} else {
		for _, item := range cities {
			tmp = append(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
		}
	}

	return removeDuplicates(tmp)
}

func subAllowedStateNames(userAccess map[string]interface{}, cities []Cities) []string {
	//fmt.Println("YOU ARE HERE\n")
	countryAccess, _ := userAccess["countries"].([]string)
	includeState, _ := userAccess["included_states"].([]string)
	excludeState, _ := userAccess["excluded_states"].([]string)

	/*	fmt.Printf("COUNTRIES ARE: %v\n", countryAccess)
		fmt.Printf("INCLUDES ARE: %v\n", includeState)
		fmt.Printf("EXCLUDES ARE: %v\n", excludeState)
		fmt.Println("*****************************************8\n")

		fmt.Printf("COUNTRIES ARE: %v\n", userAccess["countries"])
		fmt.Printf("INCLUDES ARE: %v\n", userAccess["included_states"])
		fmt.Printf("EXCLUDES ARE: %v\n", userAccess["excluded_states"])*/

	var tmp []string
	var custom []string

	for _, cnt := range countryAccess {
		for _, item := range cities {
			if item.Country == cnt {
				tmp = append(tmp, fmt.Sprintf("%v _ %v", item.State, item.Country))
			}
		}
	}

	for _, inc := range includeState {
		tmp = append(tmp, inc)
	}

	if len(excludeState) > 0 {
		for _, t := range tmp {
			for _, exc := range excludeState {
				if t != exc {
					custom = append(custom, t)
				}
			}
		}
	} else {
		for _, t := range tmp {
			custom = append(custom, t)
		}
	}

	return removeDuplicates(custom)
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] != true {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func subUserAreaNames(userAccess map[string]interface{}, cities []Cities, include []string, exclude []string) []string {
	//fmt.Println("YOU ARE HERE\n")
	countryAccess, _ := userAccess["countries"].([]string)
	includeArea, _ := userAccess["included_areas"].([]string)
	includeState, _ := userAccess["included_states"].([]string)
	excludeState, _ := userAccess["excluded_states"].([]string)
	excludeArea, _ := userAccess["excluded_areas"].([]string)

	var tmp []string

	if len(include) > 0 {
		countryAccess = []string{}
		includeState = include
		includeArea = []string{}
	} else if len(exclude) > 0 {
		excludeState = append(excludeState, exclude...)
	}

	for _, cnt := range countryAccess {
		for _, item := range cities {
			if item.Country == cnt {
				tmp = append(tmp, item.Area)
			}
		}
	}

	for _, inc := range includeArea {
		tmp = append(tmp, inc)
	}

	for _, ins := range includeState {
		for _, item := range cities {
			custom := strings.Split(ins, " _ ")
			if item.State == custom[0] {
				tmp = append(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
			}
		}
	}

	for _, exs := range excludeState {
		for _, item := range cities {
			custom := strings.Split(exs, " _ ")
			if item.State == custom[0] {
				tmp = remove(tmp, fmt.Sprintf("%v _ %v _ %v", item.Area, item.State, item.Country))
			}
		}
	}

	for _, exa := range excludeArea {
		for _, item := range tmp {
			if item == exa {
				tmp = remove(tmp, item)
			}
		}
	}

	return removeDuplicates(tmp)
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func PrepareRoorUser(input []string, cities []Cities) map[string]interface{} {
	currentUser := make(map[string]interface{})
	var countries []string
	var excludeStateAccess []string
	var excludeAreaAccess []string
	var rootStateAccess []string
	var rootAreaAccess []string

	for _, line := range input {
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "exclude") {
			custom := strings.Split(line, ":")
			custom1 := custom[1]
			//custom1, _ := strings.TrimSpace(tmp)
			if strings.Contains(custom1, "_") {
				detect := strings.Split(custom1, "_")
				if len(detect) == 2 {
					allowedStates := excludeStateName(countries, cities)
					stateName := ExistInArray(allowedStates, custom1)
					if stateName != "" {
						excludeStateAccess = append(excludeStateAccess, stateName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				} else if len(detect) == 3 {
					allowedAreas := excludeAreaName(countries, cities, excludeStateAccess)
					areaName := ExistInArray(allowedAreas, custom1)
					if areaName != "" {
						excludeAreaAccess = append(excludeAreaAccess, areaName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				}
			}
		} else if strings.Contains(lowerLine, "include") {
			allCountries := getCountryNames(cities)
			custom := strings.Split(line, ":")
			custom1 := custom[1]
			if strings.Contains(custom1, "_") {
				detect := strings.Split(custom1, "_")
				if len(detect) == 2 {
					allowedStates := extendStateName(cities, countries)
					stateName := ExistInArray(allowedStates, custom1)
					if stateName != "" {
						rootStateAccess = append(rootStateAccess, stateName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				} else if len(detect) == 3 {
					allowedAreas := extendAreaName(cities, countries, excludeStateAccess)
					areaName := ExistInArray(allowedAreas, custom1)
					if areaName != "" {
						rootAreaAccess = append(rootAreaAccess, areaName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				}
			} else {
				countryName := ExistInArray(allCountries, custom1)
				if countryName != "" {
					countries = append(countries, countryName)
				} else {
					currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
					return currentUser
				}
			}
		}
	}

	currentUser["countries"] = removeDuplicates(countries)
	currentUser["excluded_states"] = removeDuplicates(excludeStateAccess)
	currentUser["excluded_cities"] = removeDuplicates(excludeAreaAccess)
	currentUser["included_states"] = removeDuplicates(rootStateAccess)
	currentUser["included_cities"] = removeDuplicates(rootAreaAccess)

	return currentUser
}

func PrepareSubUser(input []string, cities []Cities, root map[string]interface{}) map[string]interface{} {
	currentUser := make(map[string]interface{})
	var excludeStateAccess []string
	var excludeAreaAccess []string
	var includedStateAccess []string
	var excludedStateAccess []string

	for _, line := range input {
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "exclude") {
			custom := strings.Split(line, ":")
			custom1 := custom[1]
			//custom1, _ := strings.TrimSpace(tmp)
			if strings.Contains(custom1, "_") {
				detect := strings.Split(custom1, "_")
				if len(detect) == 2 {
					allowedStates := subAllowedStateNames(root, cities)
					stateName := ExistInArray(allowedStates, custom1)
					if stateName != "" {
						excludeStateAccess = append(excludeStateAccess, stateName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				} else if len(detect) == 3 {
					allowedAreas := subUserAreaNames(root, cities, includedStateAccess, excludedStateAccess)
					areaName := ExistInArray(allowedAreas, custom1)
					if areaName != "" {
						excludeAreaAccess = append(excludeAreaAccess, areaName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted\n", line)
						return currentUser
					}
				}
			}
		} else if strings.Contains(lowerLine, "include") {
			custom := strings.Split(line, ":")
			custom1 := custom[1]
			if strings.Contains(custom1, "_") {
				detect := strings.Split(custom1, "_")
				if len(detect) == 2 {
					allowedStates := subAllowedStateNames(root, cities)
					stateName := ExistInArray(allowedStates, custom1)
					if stateName != "" {
						includedStateAccess = append(includedStateAccess, stateName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				} else if len(detect) == 3 {
					allowedAreas := subUserAreaNames(root, cities, includedStateAccess, excludedStateAccess)
					areaName := ExistInArray(allowedAreas, custom1)
					if areaName != "" {
						excludedStateAccess = append(excludedStateAccess, areaName)
					} else {
						currentUser["err"] = fmt.Sprintf("[%v] Not permitted", line)
						return currentUser
					}
				}
			}
		}
	}

	currentUser["excluded_states"] = removeDuplicates(excludeStateAccess)
	currentUser["excluded_cities"] = removeDuplicates(excludeAreaAccess)
	currentUser["included_states"] = removeDuplicates(includedStateAccess)
	currentUser["included_cities"] = removeDuplicates(excludedStateAccess)
	return currentUser
}

func ExistInArray(listOfItems []string, name string) string {
	//var item string
	for _, item := range listOfItems {
		if strings.EqualFold(strings.Replace(item, " ", "", -1), strings.Replace(name, " ", "", -1)) {
			return item
		}
	}
	return ""
}
