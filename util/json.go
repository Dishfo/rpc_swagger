package util

import "encoding/json"

func JsonContent(data interface{}) string {
	str, _ := json.Marshal(data)
	return string(str)
}
