package parse_test

import (
	"fmt"
	"strconv"

	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

type Storage struct {
	Data []interface{}
	i    int
}

func (s *Storage) Get() (interface{}, error) {
	var d interface{}
	if s.i < len(s.Data) {
		d = s.Data[s.i]
		s.i++
	} else {
		return nil, fmt.Errorf("no more data")
	}
	switch v := d.(type) {
	case error:
		return nil, v
	}
	return d, nil
}

func (s *Storage) String() string {
	return fmt.Sprintf("%v", s.Data)
}

// Fetch data with retries
// Cast data with fatal errors
func Example_processorDrivenFSM() {
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
		var a int
		var b int
		var result int
		var err error
		var state string
		for {
			state = fsm.Current()
			fmt.Println("--- State", state)
			switch state {
			case "fetch_A", "fetch_B":
				rawData, err = stor.Get()
				if err != nil {
					fsm.Event("error")
				} else {
					fsm.Event("data")
				}
			case "cast_A":
				a, err = strconv.Atoi(rawData.(string))
				if err != nil {
					fsm.Event("error")
				} else {
					fsm.Event("data")
				}
			case "cast_B":
				b, err = strconv.Atoi(rawData.(string))
				if err != nil {
					fsm.Event("error")
				} else {
					fsm.Event("data")
				}
			case "process":
				result = a + b
			}
			if state == "process" || state == "error" {
				break
			}
		}
		fmt.Printf("Final state: %s, Result: %d // %d+%d\n", state, result, a, b)
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
