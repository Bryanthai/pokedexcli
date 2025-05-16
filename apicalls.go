package main

import (
	"encoding/json"
	"net/http"
	"time"
	"math/rand"

	"github.com/Bryanthai/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(10 * time.Second)

type area struct{
	Next string
	Previous string
	Results []struct{
		Name string
	}
}

func mapApiCall(config *Config, forward bool) ([]string, error) {

	var mapURL string
	if forward {
		mapURL = config.Next
	} else {
		mapURL = config.Previous
	}

	byteVal, ok := cache.Get(mapURL)
	if ok {
		var cachedVal area
		if err := json.Unmarshal(byteVal, &cachedVal); err != nil {
			return []string{}, err
		}
		var cacheStrArr []string

		for i:=0; i < len(cachedVal.Results); i++ {
			cacheStrArr = append(cacheStrArr, cachedVal.Results[i].Name)
		}
		config.Next = cachedVal.Next
		config.Previous = cachedVal.Previous
		return cacheStrArr, nil
	}

	res, err := http.Get(mapURL)
	if err != nil {
		return []string{}, err
	}

	var areaStruct area

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&areaStruct); err != nil {
		return []string{}, err
	}

	config.Next = areaStruct.Next
	config.Previous = areaStruct.Previous

	byteData, err := json.Marshal(areaStruct)
	if err != nil {
		return []string{}, err
	}
	cache.Add(mapURL, byteData)

	var strArr []string

	for i:=0; i < len(areaStruct.Results); i++ {
		strArr = append(strArr, areaStruct.Results[i].Name)
	}

	return strArr, nil
}

type areaInfo struct {
	Pokemon_Encounters	[]struct{
		Pokemon pokemon
	}
}

type pokemon struct {
	Name	string
}

func areaApiCall(area string) ([]string, error) {
	areaUrl := "https://pokeapi.co/api/v2/location-area/" + area

	byteVal, ok := cache.Get(areaUrl)
	if ok {
		var cachedVal areaInfo
		if err := json.Unmarshal(byteVal, &cachedVal); err != nil {
			return []string{}, err
		}

		var cacheStrArr []string
		for i:=0; i < len(cachedVal.Pokemon_Encounters); i++ {
			cacheStrArr = append(cacheStrArr, cachedVal.Pokemon_Encounters[i].Pokemon.Name)
		}

		return cacheStrArr, nil
	}

	res, err := http.Get(areaUrl)
	if err != nil {
		return []string{}, err
	}

	var areaValue areaInfo

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&areaValue); err != nil {
		return []string{}, err
	}

	byteData, err := json.Marshal(areaValue)
	if err != nil {
		return []string{}, err
	}
	cache.Add(areaUrl, byteData)

	var pokemonList []string
	for i:=0; i < len(areaValue.Pokemon_Encounters); i++ {
		pokemonList = append(pokemonList, areaValue.Pokemon_Encounters[i].Pokemon.Name)
	}

	return pokemonList, nil
}

type PokemonInfo struct {
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
	BaseExperience int `json:"base_experience"`
	Height    int `json:"height"`
	Name          string `json:"name"`
}

func catchPokemon(experience int) bool {
	rand.Seed(time.Now().UnixNano())
	randNum :=  rand.Intn(experience + 100 + int(experience/4))
	if randNum > experience{
		return true
	}
	return false
}

func pokemonApiCall(pokemonName string) (PokemonInfo, bool, error) {
	pokemonUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	byteVal, ok := cache.Get(pokemonUrl)
	if ok {
		var cachedVal PokemonInfo
		if err := json.Unmarshal(byteVal, &cachedVal); err != nil {
			return PokemonInfo{}, false, err
		}

		isCatch := catchPokemon(cachedVal.BaseExperience)

		return cachedVal, isCatch, nil
	}

	res, err :=  http.Get(pokemonUrl)
	if err != nil {
		return PokemonInfo{}, false, err
	}

	var pokemonValue PokemonInfo

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&pokemonValue); err != nil {
		return PokemonInfo{}, false, err
	}

	byteData, err := json.Marshal(pokemonValue)
	if err != nil {
		return PokemonInfo{}, false, err
	}
	cache.Add(pokemonUrl, byteData)

	isCatch := catchPokemon(pokemonValue.BaseExperience)
	return pokemonValue, isCatch, nil
}