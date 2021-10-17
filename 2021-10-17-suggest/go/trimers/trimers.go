package trimers

import (
	"sort"
	"strings"
	"unicode"
)

// Высокоуровневый Predicotor: инкапсулирует модель и берёт
// на себя всю работу со строками, нормализацию символов и пр.
type Predictor struct {
	p predictor
}

// Допускаются многословные обучающие примеры. Ок ли это?
func (p *Predictor) Train(class interface{}, examples ...string) {
	for _, ex := range examples {
		p.p.train(class, canonicalSplit(ex)...)
	}
}

// Возвращаем только те классы, которые смэтчились. Клиент в праве
// добавить другие или специальные: "другое", "любое"
func (p *Predictor) Predict(example string) []interface{} {
	return scoreToList(p.p.predict(canonicalSplit(example)...))
}

func canonicalSplit(s string) [][]rune {
	r := [][]rune(nil)
	for _, w := range strings.FieldsFunc(s, isNotLetter) {
		r = append(r, []rune(strings.ReplaceAll(strings.ToLower(w), "ё", "е")))
	}
	return r
}

func isNotLetter(c rune) bool {
	return !unicode.IsLetter(c)
}

// TODO эта функция должна обеспечить устрочивую сторировку
// - знать о природе классов?
// - знать дефолтный порядок?
func scoreToList(score map[interface{}]int) []interface{} {
	list := make([]interface{}, len(score))
	i := 0
	for k := range score {
		list[i] = k
		i++
	}
	sort.Slice(list, func(a, b int) bool { return score[list[a]] > score[list[b]] })
	return list
}
