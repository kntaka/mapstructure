package mapstructure

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// CreateStructBySkpFields create anonymous struct dynamically with skip fields
func CreateStructBySkpFields(targetStrut interface{}, skipFields []string) (interface{}, error) {
	structType := reflect.TypeOf(targetStrut)
	// get target struct's fields to take needed fields from target map
	var cols []string
	for i := 0; i < structType.NumField(); i++ {
		cols = append(cols, structType.Field(i).Name)
	}

	// create a new struct
	var fields []reflect.StructField
	for _, cc := range cols {
		if contains(skipFields, cc) {
			continue
		}
		f, ok := structType.FieldByName(cc)
		if !ok {
			return reflect.Value{}, fmt.Errorf("not found filed from %s val: %s", structType, cc)
		}
		fields = append(fields, f)
	}
	st := reflect.New(reflect.StructOf(fields)).Elem()
	fmt.Printf("\n%+v\n", st)
	return st.Interface(), nil
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// MapToStruct converts map to struct
func MapToStruct(targetMap map[string]string, targetStrut interface{}) (interface{}, error) {
	structType := reflect.TypeOf(targetStrut)

	// get target struct's fields to take needed fields from target map
	var cols []string
	for i := 0; i < structType.NumField(); i++ {
		cols = append(cols, structType.Field(i).Name)
	}

	// struct without type
	anonymousSt, _ := CreateStructBySkpFields(targetStrut, []string{"CreatedAt", "UpdatedAt"})
	st := reflect.New(reflect.TypeOf(anonymousSt)).Elem()

	// set value in map to the new struct
	for _, cc := range cols {
		val, ok := targetMap[cc]
		if !ok {
			continue
		}
		f, ok := structType.FieldByName(cc)
		if !ok {
			return nil, fmt.Errorf("not found filed from %s val: %s", structType, cc)
		}
		var tmp reflect.Value
		switch f.Type.Name() {
		case "string":
			tmp = reflect.ValueOf(val)
		case "int64":
			if val == "" {
				tmp = reflect.ValueOf(int64(0))
			} else {
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("int parse error. filed: %s val: %s.: %w", cc, val, err)
				}
				tmp = reflect.ValueOf(int64(intVal))
			}
		case "Time":
			layout := "2006-01-02T15:04:05Z07:00"
			t, err := time.Parse(layout, val)
			if err != nil {
				return nil, fmt.Errorf("time parse error. filed: %s val:%s : %w", cc, val, err)
			}
			tmp = reflect.ValueOf(t)
		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("bool parse error. filed: %s val:%s : %w", cc, val, err)
			}
			tmp = reflect.ValueOf(b)
		default:
			return nil, fmt.Errorf("unknown field type: %s", f.Type.Name())
		}
		st.FieldByName(cc).Set(tmp)
	}
	return st.Interface(), nil
}
