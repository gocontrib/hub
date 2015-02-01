package sock

import "encoding/json"

func toJSON(payload interface{}) []byte {
	var b, err = json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return b
}
