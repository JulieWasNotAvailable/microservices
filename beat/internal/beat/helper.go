package beat

// import (
// 	"log"
// 	"reflect"
// )

// func fillTheStruct() {
// 	type Data struct {
// 		P1 *float64
// 		P2 *float64
// 		P3 *float64
// 	}
// 	data := &Data{}
// 	arr := []float64{1.234567543, 2.26752364, 3.325346456}

// 	val := reflect.ValueOf(data).Elem() // Dereference the pointer to the struct

// 	// Iterate over struct fields and assign from array
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		log.Println(field.Kind())
// 		if field.Kind() == reflect.Float64 && field.CanSet() {
// 			field.SetFloat(arr[i]) // Assign array value to struct field
// 		}
// 	}

// 	log.Println(data)
// }
