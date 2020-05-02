package encoding

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func convertStruct(i interface{}, t reflect.Type, pointer bool) (reflect.Value, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return zeroValue, fmt.Errorf("Cannot convert %T to struct", i)
	}

	out := reflect.New(t)

	err := Unstringify(m, out.Interface())
	if err != nil {
		return zeroValue, err
	}

	if !pointer {
		out = out.Elem()
	}

	return out, nil
}

func convertSlice(i interface{}, t reflect.Type, pointer bool) (reflect.Value, error) {
	s, ok := i.([]interface{})
	if !ok {
		return zeroValue, fmt.Errorf("Cannot convert %T to slice", i)
	}

	out := reflect.New(t)
	out.Elem().Set(reflect.MakeSlice(t, len(s), len(s)))

	for j, v := range s {
		val, err := convertType(t.Elem(), v)
		if err != nil {
			return zeroValue, err
		}

		out.Elem().Index(j).Set(val)
	}

	if !pointer {
		out = out.Elem()
	}

	return out, nil
}

func convertMap(i interface{}, t reflect.Type, pointer bool) (reflect.Value, error) {
	m, ok := i.(map[string]interface{})
	log.Printf("\tm %+v", m)
	if !ok {
		return zeroValue, fmt.Errorf("Cannot convert %T to map with string keys", i)
	}

	out := reflect.New(t)
	out.Elem().Set(reflect.MakeMap(t))

	for k, v := range m {
		log.Printf("\tconvertMap")
		log.Printf("\tt.Elem() %+v", t.Elem())
		val, err := convertType(t.Elem(), v)
		log.Printf("\tval %+v", val)
		if err != nil {
			return zeroValue, err
		}

		out.Elem().SetMapIndex(reflect.ValueOf(k), val)
	}

	if !pointer {
		out = out.Elem()
	}

	return out, nil
}

func convertString(i interface{}, pointer bool) (reflect.Value, error) {
	s, ok := i.(string)

	if !ok {
		return zeroValue, fmt.Errorf("Cannot convert %T to string", i)
	}

	if pointer {
		return reflect.ValueOf(&s), nil
	}

	return reflect.ValueOf(s), nil
}

func convertBool(i interface{}, pointer bool) (reflect.Value, error) {
	var b bool
	var err error

	switch v := i.(type) {
	case bool:
		b = v

	case string:
		b, err = strconv.ParseBool(v)
		if err != nil {
			return zeroValue, err
		}

	default:
		return zeroValue, fmt.Errorf("Cannot convert %T to bool", i)
	}

	if pointer {
		return reflect.ValueOf(&b), nil
	}

	return reflect.ValueOf(b), nil
}

func convertInt(i interface{}, pointer bool) (reflect.Value, error) {
	var n int

	switch v := i.(type) {
	case int:
		n = v

	case float64:
		n = int(v)

	case string:
		n64, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			return zeroValue, err
		}

		n = int(n64)

	default:
		return zeroValue, fmt.Errorf("Cannot convert %T to bool", i)
	}

	if pointer {
		return reflect.ValueOf(&n), nil
	}

	return reflect.ValueOf(n), nil
}

func convertFloat64(i interface{}, pointer bool) (reflect.Value, error) {
	var f float64
	var err error

	switch v := i.(type) {
	case float64:
		f = v

	case int:
		f = float64(v)

	case string:
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return zeroValue, err
		}

	default:
		return zeroValue, fmt.Errorf("Cannot convert %T to bool", i)
	}

	if pointer {
		return reflect.ValueOf(&f), nil
	}

	return reflect.ValueOf(f), nil
}

func convertType(t reflect.Type, i interface{}) (reflect.Value, error) {
	log.Printf("in convertType converting i: %+v to type %+v", i, t)
	pointer := false
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		log.Printf("convertType calling convertStruct")
		return convertStruct(i, t, pointer)

	case reflect.Slice:
		log.Printf("convertType calling convertSlice")
		return convertSlice(i, t, pointer)

	case reflect.Map:
		log.Printf("convertType calling convertMap on %+v", i)
		return convertMap(i, t, pointer)

	case reflect.String:
		log.Printf("convertType calling convertString")
		return convertString(i, pointer)

	case reflect.Bool:
		log.Printf("convertType calling convertBool")
		return convertBool(i, pointer)

	case reflect.Int:
		log.Printf("convertType calling convertInt")
		return convertInt(i, pointer)

	case reflect.Float64:
		log.Printf("convertType calling convertFloat")
		return convertFloat64(i, pointer)

	default:
		log.Printf("convertType UNSUPPORTED")
		return zeroValue, fmt.Errorf("Unsupported type %v", t)
	}
}

// Unstringify takes a stringified representation of a value
// and populates it into the supplied interface
func Unstringify(data map[string]interface{}, v interface{}) error {
	clean := make(map[string]interface{})
	log.Printf("UNSTRINGIFY data map before clean %+v", data)
	for k := range data {
		val := data[k]
		log.Printf("UNSTRINGIFY key %s: val %+v", k, val)
		strippedKey := strings.Replace(k, "/", "", 1)
		clean[strippedKey] = val
		log.Printf("UNSTRINGIFY strippedKey %s: %+v", strippedKey, clean[strippedKey])
	}

	t := reflect.TypeOf(v).Elem()
	log.Printf("UNSTRINGIFY t %+v", t)

	val := reflect.ValueOf(v).Elem()
	log.Printf("UNSTRINGIFY val2 %+v", val)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		log.Printf("UNSTRINGIFY f %+v", f)

		jsonName := f.Name
		jsonTag := strings.Split(f.Tag.Get("json"), ",")
		if len(jsonTag) > 0 && jsonTag[0] != "" {
			jsonName = jsonTag[0]
		}

		if value, ok := clean[jsonName]; ok {
			log.Printf("UNSTRINGIFY %s: %+v", jsonName, value)
			newValue, err := convertType(f.Type, value)
			if err != nil {
				return err
			}

			val.FieldByName(f.Name).Set(newValue)
		}
	}

	return nil
}
