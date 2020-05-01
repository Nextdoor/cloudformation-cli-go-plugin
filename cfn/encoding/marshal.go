package encoding

import (
	"encoding/json"
	"log"
)

// Marshal converts a value into stringified-JSON
func Marshal(v interface{}) ([]byte, error) {
	stringified, err := Stringify(v)
	log.Printf("Marshalling %s", stringified)
	if err != nil {
		return nil, err
	}

	return json.Marshal(stringified)
}
