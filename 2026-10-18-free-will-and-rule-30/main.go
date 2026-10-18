package main

import "fmt"

const rule = 30 // 0 // 4 // 50 // 18 // 30
const gens = 28

const emp = "••" // "··" // ".." // empty
const occ = "█▊" // "[]" // occupied

const colorOn = "\033[1;92m"
const colorOff = "\033[0m"

func top[T any](f bool, a, b T) T { // helper ternary operator: f ? a : b
	if f {
		return a
	}
	return b
}

func toz[T any](f bool, a T) T { // helper
	return top(f, a, *new(T))
}

func dumpRules(n int) {
	fmt.Printf("%[2]sRule: %[1]d%[3]s (%08[1]b)\n", n, colorOn, colorOff)
	for i := range 8 {
		fmt.Printf("%s%3d%s", colorOn, 7-i, colorOff)
		for j := range 3 {
			fmt.Print(top((i<<j)&4 > 0, emp, occ))
		}
		fmt.Print("  ")
	}
	fmt.Println()
	for i := range 8 {
		fmt.Print(toz(i != 0, "    "))
		fmt.Printf("%3d  ", 1<<(7-i))
		fmt.Print(colorOn + top(n&128 == 0, emp, occ) + colorOff)
		n = n << 1
	}
	fmt.Println()
}

func printLine(a []int) {
	for i, v := range a {
		c := i*2+1 == len(a) // it is center element
		fmt.Print(toz(c, colorOn) + top(v == 0, emp, occ) + toz(c, colorOff))
	}
	fmt.Println()
}

func shiftMask(v int, leftBit int) int {
	return ((v << 1) | leftBit) & 7
}

func generation(r int, a []int) []int { // generate next generation
	b := make([]int, len(a))
	v := shiftMask(0, a[0]) // start from v=0, consider left most a
	for i := range len(a) - 1 {
		v = shiftMask(v, a[i+1]) // look ahead (i+1) to consider right nearest neighbor
		b[i] = (r >> v) & 1      // just take bit v from rule r
	}
	b[len(a)-1] = (r >> shiftMask(v, 0)) & 1 // calculate right most b, assuming outbound a is 0
	return b
}

func main() {
	const size = gens*2 + 1
	fmt.Print("\033[2J\033[H") // clear, home
	dumpRules(rule)
	fmt.Println()
	l := make([]int, size) // allowed values 0 and 1 only
	l[size/2] = 1
	printLine(l)
	for range gens {
		l = generation(rule, l)
		printLine(l)
	}
}
