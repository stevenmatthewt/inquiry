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

			d := decoder{
				queryFieldName: queryFieldName,
				queryParams:    queryMap[queryFieldName],
			}

			switch fieldValue.Kind() {
			// This case statemtent could be made a part of decode()
			// but that would enable this to supprt nested arrays.
			// I'm not certain that's something that I want to
			// support at this time.
			case reflect.Slice, reflect.Array:
				inputSlice := queryMap[queryFieldName]
				outputSlice := reflect.MakeSlice(fieldValue.Type(), 0, len(inputSlice))
				for i := 0; i < len(inputSlice); i++ {
					tempValue := reflect.New(fieldValue.Type().Elem())
					data := inputSlice[i]
					err = decoder{
						queryFieldName: d.queryFieldName,
						queryParams:    []string{data},
					}.decode(tempValue.Elem())
					if err != nil {
						cumulativeErrs = append(cumulativeErrs, err)
						continue
					}
					outputSlice = reflect.Append(outputSlice, tempValue.Elem())
				}
				fieldValue.Set(outputSlice)
			default:
				err = d.decode(fieldValue)
				if err != nil {
					cumulativeErrs = append(cumulativeErrs, err)
				}
			}
		}
	}

	if len(cumulativeErrs) > 0 {
		return fmt.Errorf("%+v", cumulativeErrs)
	}

	return nil
}

type decoder struct {
	queryParams    []string
	queryFieldName string
}

func (d decoder) decode(value reflect.Value) error {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(d.queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(d.queryParams) < 1 {
			return fmt.Errorf("must be at least one instance of %s", d.queryFieldName)
		}
		val, err := strconv.Atoi(d.queryParams[0])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to map field %s", d.queryFieldName))
		}
		if value.OverflowInt(int64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetInt(int64(val))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(d.queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(d.queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		val, err := strconv.Atoi(d.queryParams[0])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to map field %s", d.queryFieldName))
		}
		if val < 0 {
			return fmt.Errorf("value of %s underflows type %s", d.queryFieldName, value.Kind().String())
		}
		if value.OverflowUint(uint64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetUint(uint64(val))
	case reflect.Float32, reflect.Float64:
		if len(d.queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(d.queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		size := 32
		if value.Kind() == reflect.Float64 {
			size = 64
		}
		val, err := strconv.ParseFloat(d.queryParams[0], size)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to map field %s", d.queryFieldName))
		}
		if value.OverflowFloat(float64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetFloat(float64(val))
	case reflect.String:
		if len(d.queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(d.queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		value.SetString(d.queryParams[0])
	default:
		return fmt.Errorf("%s is not a supported type", value.Kind().String())
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
