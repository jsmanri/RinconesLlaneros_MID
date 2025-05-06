package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
)

func Metodo_get_all(host, endpoint string) ([]byte, error) {
	url := beego.AppConfig.String(host) + endpoint
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

func ProcesarJsonArreglos(datos []byte) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	err := json.Unmarshal(datos, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ProcesarJson(datos []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err5 := json.Unmarshal(datos, &result)
	if err5 != nil {
		log.Fatal(err5)
		return nil, err5
	}
	return result, nil
}

func ConvertInterfaceToSliceMap(input interface{}) ([]map[string]interface{}, error) {
	// Afirmar que es un slice de interface{}
	list, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("no es un []interface{}")
	}

	// Convertir cada elemento a map[string]interface{}
	var result []map[string]interface{}
	for i, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("elemento %d no es un map[string]interface{}", i)
		}
		result = append(result, m)
	}

	return result, nil
}

func ConvertToSliceOfMaps(input interface{}) ([]map[string]interface{}, error) {
	// Afirmar que el input es un slice de interfaces
	rawSlice, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("input is not a slice")
	}

	// Crear el slice de map[string]interface{}
	result := make([]map[string]interface{}, len(rawSlice))

	for i, item := range rawSlice {
		// Intentar convertir cada elemento a map[string]interface{}
		elem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a map[string]interface{}", i)
		}
		result[i] = elem
	}

	return result, nil

}
func GroupByID(items []map[string]interface{}) map[interface{}][]map[string]interface{} {
    grouped := make(map[interface{}][]map[string]interface{})

    for i, item := range items {
		sitioturistico := items[i]["IdSitiosTuristicos"]
		idsitio := sitioturistico.(map[string]interface{})["Id"]
        id := idsitio
        grouped[id] = append(grouped[id], item)
    }

    return grouped
}

func PromedioNumeros(data []interface{}) float64 {
    var suma float64
    var cantidad int

    for _, v := range data {
        switch num := v.(type) {
        case int:
            suma += float64(num)
            cantidad++
        case float64:
            suma += num
            cantidad++
        case float32:
            suma += float64(num)
            cantidad++
        case int64:
            suma += float64(num)
            cantidad++
        case int32:
            suma += float64(num)
            cantidad++
        default:
            // Ignorar tipos no numéricos
        }
    }

    if cantidad == 0 {
        return 0 // evitar división por cero
    }

    return suma / float64(cantidad)
}
