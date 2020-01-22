package parse_test

import (
	"fmt"
	"strings"

	"github.com/looplab/fsm"
)

type Digit struct {
	Number string
	FSM    *fsm.FSM
}

func (d *Digit) enterState(e *fsm.Event) {
	eventName := e.Event // surprise!
	if strings.HasPrefix(e.Dst, "save_") {
		d.Number = d.Number + eventName
	}
	//fmt.Printf("State %s [%s]\n", e.Dst, eventName)
}

// Parse strings like "[-+]?[01]+$" to "-?[01]+"
// For example
// 101$ -> 101
// +10$ -> 10
// -10$ -> -10
// -+1$ -> error
// ++0$ -> error
func NewDigitParser() *Digit {
	d := &Digit{}

	d.FSM = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "$", Src: []string{"start"}, Dst: "error"},
			{Name: "+", Src: []string{"start"}, Dst: "skip_plus"},
			{Name: "-", Src: []string{"start"}, Dst: "save_minus"},
			{Name: "0", Src: []string{"start"}, Dst: "save_digit"},
			{Name: "1", Src: []string{"start"}, Dst: "save_digit"},
			{Name: "$", Src: []string{"skip_plus"}, Dst: "error"},
			{Name: "+", Src: []string{"skip_plus"}, Dst: "error"},
			{Name: "-", Src: []string{"skip_plus"}, Dst: "error"},
			{Name: "0", Src: []string{"skip_plus"}, Dst: "save_digit"},
			{Name: "1", Src: []string{"skip_plus"}, Dst: "save_digit"},
			{Name: "$", Src: []string{"save_minus"}, Dst: "error"},
			{Name: "+", Src: []string{"save_minus"}, Dst: "error"},
			{Name: "-", Src: []string{"save_minus"}, Dst: "error"},
			{Name: "0", Src: []string{"save_minus"}, Dst: "save_digit"},
			{Name: "1", Src: []string{"save_minus"}, Dst: "save_digit"},
			{Name: "$", Src: []string{"save_digit"}, Dst: "fin"},
			{Name: "+", Src: []string{"save_digit"}, Dst: "error"},
			{Name: "-", Src: []string{"save_digit"}, Dst: "error"},
			{Name: "0", Src: []string{"save_digit"}, Dst: "save_digit"},
			{Name: "1", Src: []string{"save_digit"}, Dst: "save_digit"},
		},
		fsm.Callbacks{
			"after_event": func(e *fsm.Event) { d.enterState(e) }, // enter_state works oddly
		},
	)

	return d
}

func Example_classicalDataDrivenFSM() {
	// Good:
	// - every pair (state, event) leads to certain state and certain action
	// - we obtain clear and compact transition table/function
	// - it guarantees O(N) and finiteness
	// What wrong here:
	// - what if we obtained an error while reading input? it is out of scope of our FSM
	// - what if we obtained an error in enterState processor? how to implement retries etc?
	for _, input := range []string{"11$", "+10$", "-+", "1+1"} {
		fmt.Println("--- Input:", input)
		p := NewDigitParser()
		for _, c := range input {
			//fmt.Println("Char:", string(c))
			p.FSM.Event(string(c))
		}
		fmt.Println("Final state:", p.FSM.Current())
		if p.FSM.Current() == "fin" {
			fmt.Println("Result", p.Number)
		}
	}
	// Output:
	// --- Input: 11$
	// Final state: fin
	// Result 11
	// --- Input: +10$
	// Final state: fin
	// Result 10
	// --- Input: -+
	// Final state: error
	// --- Input: 1+1
	// Final state: error
}
