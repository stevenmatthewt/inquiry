package inquiry

import (
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
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("unexpected panic when unmarshalling query params: %+v", r)
			}
		}
	}()

	u := unmarshaller{
		queryMap: queryMap,
	}
	return u.unmarshal(v)

}

// unmashaller is used internally to help share information
// during the unmarshalling process
type unmarshaller struct {
	errors   []string
	queryMap map[string][]string
}

func (u *unmarshaller) unmarshal(v interface{}) error {
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

		d := decoder{
			queryFieldName: opts[0],
			opts:           opts[1:],
		}

		inputSlice := u.queryMap[d.queryFieldName]

		fieldValue := s.Field(i)
		if !fieldValue.IsValid() {
			continue
		}

		if fieldValue.CanSet() {
			switch fieldValue.Kind() {
			// This case statemtent could be made a part of decode()
			// but that would enable this to supprt nested arrays.
			// I'm not certain that's something that I want to
			// support at this time.
			case reflect.Slice, reflect.Array:
				outputSlice := reflect.MakeSlice(fieldValue.Type(), 0, len(inputSlice))
				for i := 0; i < len(inputSlice); i++ {
					tempValue := reflect.New(fieldValue.Type().Elem())
					data := inputSlice[i]
					err := decoder{
						queryFieldName: d.queryFieldName,
					}.decode([]string{data}, tempValue.Elem())
					if err != nil {
						u.addErr(err)
						continue
					}
					outputSlice = reflect.Append(outputSlice, tempValue.Elem())
				}
				fieldValue.Set(outputSlice)
			default:
				err := d.decode(inputSlice, fieldValue)
				if err != nil {
					u.addErr(err)
				}
			}
		}
	}

	return u.getErr()
}

func (u *unmarshaller) addErr(err error) {
	u.errors = append(u.errors, err.Error())
}

func (u *unmarshaller) getErr() error {
	if len(u.errors) > 0 {
		return fmt.Errorf("invalid query string format: [%s]", strings.Join(u.errors, ", "))
	}

	return nil
}

// decoder is responsible for decoding a **single** field
type decoder struct {
	opts           []string
	queryFieldName string
}

func (d decoder) decode(queryParams []string, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(queryParams) < 1 {
			return fmt.Errorf("must be at least one instance of %s", d.queryFieldName)
		}
		val, err := strconv.Atoi(queryParams[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid integer", d.queryFieldName)
		}
		if value.OverflowInt(int64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetInt(int64(val))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		val, err := strconv.Atoi(queryParams[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid unsigned integer", d.queryFieldName)
		}
		if val < 0 {
			return fmt.Errorf("value of %s underflows type %s", d.queryFieldName, value.Kind().String())
		}
		if value.OverflowUint(uint64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetUint(uint64(val))
	case reflect.Float32, reflect.Float64:
		if len(queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		size := 32
		if value.Kind() == reflect.Float64 {
			size = 64
		}
		val, err := strconv.ParseFloat(queryParams[0], size)
		if err != nil {
			return fmt.Errorf("%s is not a valid float", d.queryFieldName)
		}
		if value.OverflowFloat(float64(val)) {
			return fmt.Errorf("value of %s overflows type %s", d.queryFieldName, value.Kind().String())
		}
		value.SetFloat(float64(val))
	case reflect.String:
		if len(queryParams) > 1 {
			return fmt.Errorf("more than one instance of %s provider", d.queryFieldName)
		}
		if len(queryParams) < 1 {
			return fmt.Errorf("must be at least on instance of %s", d.queryFieldName)
		}
		value.SetString(queryParams[0])
	default:
		return fmt.Errorf("%s is not a supported type", value.Kind().String())
	}

	return nil
}
