package inquiry

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const tagName = "query"

func UnmarshalMap(queryMap map[string][]string, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// var ok bool
			// err, ok = r.(error)
			// if !ok {
			panic(r)
			// }
		}
	}()
	var cumulativeErrs []error
	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	s := reflect.ValueOf(v).Elem()
	t := s.Type()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("output type must be of type struct; %s was given", s.Kind().String())
	}

	// Get the type and kind of our user variable
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get(tagName)

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		opts := strings.Split(tag, ",")
		if len(opts) < 1 {
			return errors.New("Invalid struct tag format")
		}
		queryFieldName := opts[0]
		opts = opts[1:]

		fieldValue := s.Field(i)
		if !fieldValue.IsValid() {
			continue
		}

		if fieldValue.CanSet() {
			// var queryVal interface{}
			// switch fieldValue.Kind() {
			// case reflect.Slice, reflect.Array:
			// 	queryVal = queryMap[queryFieldName]
			// default:
			// 	if len(queryMap[queryFieldName]) <= 0 {
			// 		cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("must be at least on instance of %s", queryFieldName)))
			// 	}
			// 	if len(queryMap[queryFieldName]) > 1 {
			// 		cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("more than one instance of %s provider", queryFieldName)))
			// 		continue
			// 	}
			// 	queryVal = queryMap[queryFieldName][0]
			// }
			// err = storeValue(queryVal, fieldValue.Addr().Interface())
			// if err != nil {
			// 	cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("failed to unmarshal field %s", queryFieldName)))
			// }

			switch fieldValue.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if len(queryMap[queryFieldName]) > 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("more than one instance of %s provider", queryFieldName))
					continue
				}
				if len(queryMap[queryFieldName]) < 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("must be at least on instance of %s", queryFieldName))
					continue
				}
				val, err := strconv.Atoi(queryMap[queryFieldName][0])
				if err != nil {
					cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("failed to map field %s", queryFieldName)))
					continue
				}
				if fieldValue.OverflowInt(int64(val)) {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("value of %s overflows type %s", queryFieldName, fieldValue.Kind().String()))
				}
				fieldValue.SetInt(int64(val))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if len(queryMap[queryFieldName]) > 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("more than one instance of %s provider", queryFieldName))
					continue
				}
				if len(queryMap[queryFieldName]) < 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("must be at least on instance of %s", queryFieldName))
					continue
				}
				val, err := strconv.Atoi(queryMap[queryFieldName][0])
				if err != nil {
					cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("failed to map field %s", queryFieldName)))
					continue
				}
				if val < 0 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("value of %s underflows type %s", queryFieldName, fieldValue.Kind().String()))
				}
				if fieldValue.OverflowUint(uint64(val)) {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("value of %s overflows type %s", queryFieldName, fieldValue.Kind().String()))
				}
				fieldValue.SetUint(uint64(val))
			case reflect.Float32, reflect.Float64:
				if len(queryMap[queryFieldName]) > 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("more than one instance of %s provider", queryFieldName))
					continue
				}
				if len(queryMap[queryFieldName]) < 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("must be at least on instance of %s", queryFieldName))
					continue
				}
				size := 32
				if fieldValue.Kind() == reflect.Float64 {
					size = 64
				}
				val, err := strconv.ParseFloat(queryMap[queryFieldName][0], size)
				if err != nil {
					cumulativeErrs = append(cumulativeErrs, errors.Wrap(err, fmt.Sprintf("failed to map field %s", queryFieldName)))
					continue
				}
				if fieldValue.OverflowFloat(float64(val)) {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("value of %s overflows type %s", queryFieldName, fieldValue.Kind().String()))
				}
				fieldValue.SetFloat(float64(val))
			case reflect.String:
				if len(queryMap[queryFieldName]) > 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("more than one instance of %s provider", queryFieldName))
					continue
				}
				if len(queryMap[queryFieldName]) < 1 {
					cumulativeErrs = append(cumulativeErrs, fmt.Errorf("must be at least on instance of %s", queryFieldName))
					continue
				}
				fieldValue.SetString(queryMap[queryFieldName][0])
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				fmt.Println("WE GOT A BIG ONE!")
			}
		}
	}

	if len(cumulativeErrs) > 0 {
		return fmt.Errorf("%+v", cumulativeErrs)
	}

	return nil
}

func storeValue(data interface{}, variable interface{}) error {
	var inputStruct struct {
		Temp interface{}
	}
	var outputStruct struct {
		Temp interface{}
	}
	outputStruct.Temp = variable

	inputStruct.Temp = data
	bytes, err := json.Marshal(inputStruct)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &outputStruct)
	fmt.Printf("just stored: %x\n", variable)
	return err
}
