package trimers

// Низкоуровневая модель. Работает только с рунами
// и массивами рун (чтобы избежать множественных
// засад типа len(string) != len([]rune(string)))
//
// Мы не храним счётчик триграмм, считаем, что
// - лучше сблизить скоры для топа результатов
// - у нас очень маленький набор синонимов
//   и перекрытия, — скорее недогляд, чем намерение
// Иными словами: триграмма может принадлежать
// нескольким классам, но тогда все классы получают
// равный буст
type predictor struct {
	m []map[rune]map[interface{}]struct{}
}

func (p *predictor) train(class interface{}, ex ...[]rune) {
	for _, e := range ex {
		for i, s := range trimers(e) {
			for len(p.m) <= i {
				p.m = append(p.m, map[rune]map[interface{}]struct{}{})
			}
			r := p.m[i]
			c := r[s]
			if c == nil {
				r[s] = map[interface{}]struct{}{class: {}}
			} else {
				c[class] = struct{}{}
			}
		}
	}
}

func (p *predictor) predict(ex ...[]rune) map[interface{}]int {
	score := map[interface{}]int{}
	for _, e := range ex {
		for j, s := range trimers(e) {
			for d, w := range []int{1, 2, 5, 2, 1} {
				i := j + d - 2
				if i < 0 || i >= len(p.m) {
					continue
				}
				r := p.m[i]
				c := r[s]
				if c != nil {
					for class := range c {
						score[class] = score[class] + w
					}
				}
			}
		}
	}
	return score
}

// Возвращает не сами триграммы, а суммы
// Суммы вычисляются без потери информции, для букв одного алфавита
// 37 — простое число, большее количества символов
func trimers(s []rune) []rune {
	l := len(s) - 2
	if l <= 0 {
		return nil
	}
	r := make([]rune, l)
	for i := 0; i < l; i++ {
		r[i] = (s[i]*37+s[i+1])*37 + s[i+2]
	}
	return r
}
