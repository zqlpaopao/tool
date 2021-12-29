package src

import (
	"encoding/json"
	"strings"
)

// MarshalWithInterface Uncertain type conversion. When uncertain,
//int will be converted to float, resulting in loss of precision
func MarshalWithInterface(req string) (test interface{}, err error) {
	decoder := json.NewDecoder(strings.NewReader(req))
	decoder.UseNumber()
	err = decoder.Decode(&test)
	return
}
