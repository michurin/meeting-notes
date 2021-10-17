package main

import (
	"encoding/json"
	"fmt"
	"os"
	"trimers/trimers"
)

func train(p *trimers.Predictor) { // TODO interface
	fh, err := os.Open("../data.json")
	if err != nil {
		panic(err)
	}
	data := map[string][]string{}
	err = json.NewDecoder(fh).Decode(&data)
	if err != nil {
		panic(err)
	}
	for class, ex := range data {
		p.Train(class, class)
		p.Train(class, ex...)
	}
}

func predict(p *trimers.Predictor, ex string) { // TODO interface
	classes := p.Predict(ex)
	fmt.Println(ex)
	for _, c := range classes {
		fmt.Println("\t" + c.(string))
	}
}

func main() {
	p := new(trimers.Predictor)
	train(p)
	for _, ex := range os.Args[1:] {
		predict(p, ex)
	}
}
