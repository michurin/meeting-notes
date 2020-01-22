package parse_test

import (
	"fmt"
	"strconv"

	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

type Accumulator struct {
	A int
	B int
}

func event(err error) string {
	if err != nil {
		return "error"
	}
	return "data"
}

// Fetch data with retries
// Cast data with fatal errors
func Example_processorDrivenFSMWithProcMap() {
	for _, stor := range []*Storage{
		{Data: []interface{}{"1", "2"}},
		{Data: []interface{}{"3", errors.New("Error"), "4"}},
		{Data: []interface{}{"x"}},
	} {
		fmt.Println("Input:", stor)
		fsm := fsm.NewFSM(
			"fetch_A",
			fsm.Events{
				{Name: "data", Src: []string{"fetch_A"}, Dst: "cast_A"},
				{Name: "error", Src: []string{"fetch_A"}, Dst: "fetch_A"},
				{Name: "data", Src: []string{"cast_A"}, Dst: "fetch_B"},
				{Name: "error", Src: []string{"cast_A"}, Dst: "error"},
				{Name: "data", Src: []string{"fetch_B"}, Dst: "cast_B"},
				{Name: "error", Src: []string{"fetch_B"}, Dst: "fetch_B"},
				{Name: "data", Src: []string{"cast_B"}, Dst: "process"},
				{Name: "error", Src: []string{"cast_B"}, Dst: "error"},
			},
			fsm.Callbacks{},
		)
		var rawData interface{}
		var acc Accumulator
		var result int
		var err error
		var state string

		// Just POC how we could use maps and keep everything typed
		fetch := func(stor *Storage) (interface{}, error) {
			return stor.Get()
		}
		opMapRead := map[string]func(storage *Storage) (interface{}, error){
			"fetch_A": fetch,
			"fetch_B": fetch,
		}
		castA := func(rawData interface{}, acc *Accumulator) (err error) {
			acc.A, err = strconv.Atoi(rawData.(string))
			return
		}
		castB := func(rawData interface{}, acc *Accumulator) (err error) {
			acc.B, err = strconv.Atoi(rawData.(string))
			return
		}
		opMapCast := map[string]func(rawData interface{}, acc *Accumulator) error{
			"cast_A": castA,
			"cast_B": castB,
		}

		for {
			state = fsm.Current()
			fmt.Println("--- State", state)
			readOperation, ok := opMapRead[state]
			if ok {
				rawData, err = readOperation(stor)
				fsm.Event(event(err))
			}
			castOperation, ok := opMapCast[state]
			if ok {
				err = castOperation(rawData, &acc)
				fsm.Event(event(err))
			}
			switch state { // All the rest could be moved to maps too for sure
			case "process":
				result = acc.A + acc.B
			}
			if state == "process" || state == "error" {
				break
			}
		}
		fmt.Printf("Final state: %s, Result: %d // %d+%d\n", state, result, acc.A, acc.B)
	}
	// Output:
	// Input: [1 2]
	// --- State fetch_A
	// --- State cast_A
	// --- State fetch_B
	// --- State cast_B
	// --- State process
	// Final state: process, Result: 3 // 1+2
	// Input: [3 Error 4]
	// --- State fetch_A
	// --- State cast_A
	// --- State fetch_B
	// --- State fetch_B
	// --- State cast_B
	// --- State process
	// Final state: process, Result: 7 // 3+4
	// Input: [x]
	// --- State fetch_A
	// --- State cast_A
	// --- State error
	// Final state: error, Result: 0 // 0+0
}
