package main

import (
	"course-project/app"
	"course-project/entities"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type any interface{}

func coalesce(items ...any) any {
	for _, item := range items {
		if item != nil {
			return item
		}
	}
	return nil
}

func IngestPlaygrounds(d app.DAO, playgroundsPath string) error {
	rawData, err := os.ReadFile(playgroundsPath)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return fmt.Errorf("Error while unmarshaling geojson data: %v", err)
	}
	features := data["features"].([]interface{})
	for _, rawFeature := range features {
		feature := rawFeature.(map[string]interface{})
		props := feature["properties"].(map[string]interface{})
		geom := feature["geometry"].(map[string]interface{})
		coords := geom["coordinates"].([]interface{})[0].([]interface{})
		areaVal, ok := coalesce(props["new_plost"], props["plost_old"], "0").(string)
		var area int
		if ok {
			area, err = strconv.Atoi(areaVal)
			if err != nil {
				area = 0
			}
		} else {
			area = 0
		}
		playground := entities.Playground{
			SiteNumber: coalesce(props["nobekt_new"], props["nobekt_old"], "").(string),
			Latitude:   coords[0].(float64),
			Longitude:  coords[1].(float64),
			Area:       area,
			Location:   coalesce(props["new_mestopolozh"], props["mestopolozh_old"], "").(string),
			Ownership:  coalesce(props["new_sobstvenos"], props["sobstvenost_old"], "").(string),
		}
		err = d.CreatePlayground(&playground)
		if err != nil {
			return fmt.Errorf("Error while creating playground: %v", err)
		}
	}
	return nil
}
