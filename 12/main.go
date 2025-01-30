package main

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	// struct tags tell decoders/encoders how to interpret struct fields
	// this tells json to treat the field Page as page in json
	Page  int      `json:"page"`
	Words []string `json:"words,omitempty"`
}

func main() {
	r := Response{
		Page:  1,
		Words: []string{"up", "in", "out"},
	}

	j, _ := json.Marshal(r)
	fmt.Println(string(j))

	var v Response
	json.Unmarshal(j, &v)
	fmt.Printf("%#v\n", v)
}
