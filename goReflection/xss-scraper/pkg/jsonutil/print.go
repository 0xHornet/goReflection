package jsonutil

import (
	"encoding/json"
	"fmt"
)

func Print(v interface{}) {
	bs, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(bs))
}
