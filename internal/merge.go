package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func DeserializeKeyValue(line string) (KeyValue, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 4 {
		return KeyValue{}, fmt.Errorf("invalid format: %s", line)
	}

	key := parts[0]
	valueStr := parts[1]
	valueType := parts[3]

	var value interface{}
	var err error

	switch valueType {
	case "int":
		value, err = strconv.Atoi(valueStr)
	case "string":
		value = valueStr
	// Add more types as needed
	default:
		err = fmt.Errorf("unsupported value type: %s", valueType)
	}

	if err != nil {
		return KeyValue{}, err
	}

	return KeyValue{Key: key, Value: value}, nil
}
func (m *Master) mergeFiles() (map[string][]int, error) {
	fmt.Println("starting merging ..")
	aggregateData := make(map[string][]int) // make type agnostic TODO
	for _, task := range m.Tasks {
		outputFile := task.Output
		file, err := os.Open(outputFile)
		if err != nil {
			return nil, fmt.Errorf("Error opening file for task %v-> %v", task.ID, task.Output)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			kv, err := DeserializeKeyValue(scanner.Text())
			if err != nil {
				return nil, fmt.Errorf("Error scanning file for task %v-> %v", task.ID, task.Output)
			}
			val, ok := kv.Value.(int) // make agnostic TODO
			if !ok {
				return nil, fmt.Errorf("Error value expected int got %T", kv.Value)

			}
			aggregateData[kv.Key] = append(aggregateData[kv.Key], val)
		}
	}
	// for k := range aggregateData {
	// 	values := aggregateData[k]
	// 	fmt.Printf("%v : ", k)
	// 	for _, val := range values {
	// 		fmt.Printf("%v ", val)
	// 	}
	// 	println("\n\n\n")
	// }
	return aggregateData, nil
}
