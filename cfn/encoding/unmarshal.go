package encoding

import (
	"encoding/json"
	"log"
)

// Unmarshal converts stringified-JSON into the passed-in type
func Unmarshal(data []byte, v interface{}) error {
	var dataMap map[string]interface{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	log.Printf("dataMap %+v", dataMap)
	Unstringify(dataMap, v)

	return nil
}
