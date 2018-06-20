package db

import "fmt"

var data = make(map[string]interface{}, 0)

func Insert(id string, v interface{}) error {
	data[id] = v
	return nil
}

func Query(id string) (interface{}, error) {
	if d, ok := data[id]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("%s not foud", id)
}

func Delete(id string) {
	if _, ok := data[id]; ok {
		delete(data, id)
	}
}



