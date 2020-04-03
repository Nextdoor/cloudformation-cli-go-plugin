package encoding

import (
	"encoding/json"
	"log"
	"strconv"
)

// Unmarshal converts stringified-JSON into the passed-in type
func Unmarshal(data []byte, v interface{}) error {
	var dataMap map[string]interface{}
	log.Printf("\nUnmarshal: byte data: %+v\n", data)
	log.Printf("\nUnmarshal: string data: %+v\n", string(data))
	log.Printf("\nUnmarshal: dataMap: %+v\n", dataMap)

	dataStr, err := strconv.Unquote(string(data))
	if err != nil {
		log.Printf("Unmarshal: err Unquoting")
	}
	log.Printf("\nUnmarshal: dataStr %s\n", dataStr)
	err = json.Unmarshal([]byte(dataStr), &dataMap)
	if err != nil {
		return err
	}

	Unstringify(dataMap, v)

	return nil
}
