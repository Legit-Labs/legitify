package formatter_test

import "encoding/json"

func DeserializeJson(bytes []byte) (output map[string]interface{}, err error) {
	err = json.Unmarshal(bytes, &output)
	return
}
