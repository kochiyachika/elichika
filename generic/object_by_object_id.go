package generic

import (
	"elichika/utils"

	// "fmt"

	"encoding/json"
)

// Common pattern of [object_1_id, object_1, object_2_id, object_2...]
// if null give an empty array
// need to pass the pointer type of a type because I don't know golang

type ObjectByObjectIDReadInterface interface { // only unmarshal
	SetID(int64)
}

type ObjectByObjectIDRead[T ObjectByObjectIDReadInterface] struct {
	Objects []T
}

func (oboid *ObjectByObjectIDRead[T]) UnmarshalJSON(data []byte) error {
	oboid.Objects = []T{}
	if string(data) == "null" {
		return nil
	}
	arr := []any{}
	err := json.Unmarshal(data, &arr)
	utils.CheckErr(err)
	for i, _ := range arr {
		if i%2 == 1 {
			continue
		}
		bytes, err := json.Marshal(arr[i+1])
		utils.CheckErr(err)
		ptp := new(T)                         // pointer to pointer of the original type
		err = json.Unmarshal(bytes, ptp)      // unmarshal into the pointer of pointer the original type, create a pointer and then an object
		(*ptp).SetID(int64(arr[i].(float64))) // set the id
		utils.CheckErr(err)
		oboid.Objects = append(oboid.Objects, *ptp) // append the pointer of the original type
	}
	return nil
}

type ObjectByObjectIDWriteInterface interface { // only marshal
	ID() int64
}

type ObjectByObjectIDWrite[T ObjectByObjectIDWriteInterface] struct {
	Length  int
	Objects []T
}

func (oboid ObjectByObjectIDWrite[T]) MarshalJSON() ([]byte, error) {
	arr := []any{}
	for _, object := range oboid.Objects {
		arr = append(arr, object.ID())
		arr = append(arr, object)
	}
	bytes, err := json.Marshal(&arr)
	utils.CheckErr(err)
	return bytes, nil
}

// append a pointer
func (oboid *ObjectByObjectIDWrite[T]) PushBack(ptr T) {
	oboid.Length += 1
	oboid.Objects = append(oboid.Objects, ptr)
}

// append a new pointer at the end, the pointer point to a zero value
// then return the pointer and the index
func (oboid *ObjectByObjectIDWrite[T]) AppendNew() T {
	ptp := new(T) // pointer to pointer, that hold nothing
	// a bit dirty but it works
	err := json.Unmarshal([]byte("{}"), ptp) // unmarshal into the pointer of pointer the original type, create a pointer and then an object
	utils.CheckErr(err)

	oboid.Objects = append(oboid.Objects, *ptp) // append the pointer and return it
	oboid.Length += 1
	return *ptp
}
