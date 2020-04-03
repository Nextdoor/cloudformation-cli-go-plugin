package encoding

import (
	"encoding/json"
	"log"
)

// Unmarshal converts stringified-JSON into the passed-in type
func Unmarshal(data []byte, v interface{}) error {
	var dataMap map[string]interface{}
	log.Printf("\nUnmarshal: byte data: %+v\n", data)
	log.Printf("\nUnmarshal: dataMap: %+v\n", dataMap)
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	Unstringify(dataMap, v)

	return nil
}
