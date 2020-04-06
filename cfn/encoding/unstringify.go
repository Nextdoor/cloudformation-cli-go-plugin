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
	log.Printf("\nin convertStruct, m %+v\n", m)
	if !ok {
		log.Printf("\nCannot convert %T to struct\n", i)
		return zeroValue, fmt.Errorf("Cannot convert %T to struct", i)
	}

	out := reflect.New(t)
	log.Printf("\nout initial in convertStruct %+v\n", out)

	err := Unstringify(m, out.Interface())
	//err := mapstructure.Decode(m, out.Interface())
	log.Printf("after Unstringify m: %+v, out: %+v", m, out)
	if err != nil {
		log.Printf("\nreturning zeroValue from convertStruct\n")
		return zeroValue, err
	}

	if !pointer {
		log.Printf("\nout.Elem() %+v\n", out.Elem())
		out = out.Elem()
	}

	log.Printf("\nreturning out: %+v, nil\n", out)
	return out, nil
}

func convertSlice(i interface{}, t reflect.Type, pointer bool) (reflect.Value, error) {
	log.Printf("\nConverting slice of type %T\n", t)
	log.Printf("\n Type  is : %s\n", t.Kind())

	s, ok := i.([]interface{})
	if !ok {
		log.Printf("\nFailed to convert %T to slice\n", i)
		return zeroValue, fmt.Errorf("Cannot convert %T to slice", i)
	}

	out := reflect.New(t)
	log.Printf("Out is: %+v", out)
	log.Printf("\n\n")
	out.Elem().Set(reflect.MakeSlice(t, len(s), len(s)))

	for j, v := range s {
		log.Printf("\nIterating over slice; j: %+v, v: %+v\n", j, v)
		log.Printf("\nCalling convertType with t.Elem()%+v\n", t.Elem())

		val, err := convertType(t.Elem(), v)
		log.Printf("\nval: %+v, err: %+v\n", val, err)
		if err != nil {
			log.Printf("\nelem is not of type type's Kind is not Array, Chan, Map, Ptr, or Slice: %s\n", err.Error())
			return zeroValue, err
		}

		log.Printf("\nSetting out.Elem().Index(j).Set(val) with val %+v\n", val)
		out.Elem().Index(j).Set(val)
	}

	if !pointer {
		out = out.Elem()
	}

	log.Printf("\nFinal out is:  %+v\n", out)
	return out, nil
}

func convertMap(i interface{}, t reflect.Type, pointer bool) (reflect.Value, error) {
	m, ok := i.(map[string]interface{})
	if !ok {
		return zeroValue, fmt.Errorf("Cannot convert %T to map with string keys", i)
	}

	out := reflect.New(t)
	out.Elem().Set(reflect.MakeMap(t))

	for k, v := range m {
		val, err := convertType(t.Elem(), v)
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
	pointer := false
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		log.Printf("\nConverting struct %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertStruct(i, t, pointer)

	case reflect.Slice:
		log.Printf("\nConverting slice %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertSlice(i, t, pointer)

	case reflect.Map:
		log.Printf("\nConverting map %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertMap(i, t, pointer)

	case reflect.String:
		log.Printf("\nConverting string %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertString(i, pointer)

	case reflect.Bool:
		log.Printf("\nConverting bool %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertBool(i, pointer)

	case reflect.Int:
		log.Printf("\nConverting int %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertInt(i, pointer)

	case reflect.Float64:
		log.Printf("\nConverting float64 %v\n", t)
		log.Printf("\nInterface %v\n", i)
		return convertFloat64(i, pointer)

	default:
		log.Printf("\nUnsupported type %v\n", t)
		return zeroValue, fmt.Errorf("Unsupported type %v", t)
	}
}

// Unstringify takes a stringified representation of a value
// and populates it into the supplied interface
func Unstringify(data map[string]interface{}, v interface{}) error {
	log.Printf("Unstringify data: %+v, v: %+v", data, v)
	for k := range data {
		log.Printf("\nkey: %s\n", k)
	}
	t := reflect.TypeOf(v).Elem()
	log.Printf("\nt: %+v\n", t)

	val := reflect.ValueOf(v).Elem()
	log.Printf("\nUnstringify val: %+v\n", val)

	log.Printf("\nt.NumField(): %d \n", t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		log.Printf("\nUnstringify f: %+v\n", f)

		jsonName := f.Name
		log.Printf("\njsonName %s \n", jsonName)
		jsonTag := strings.Split(f.Tag.Get("json"), ",")
		log.Printf("\njsonTag %s \n", jsonTag)
		if len(jsonTag) > 0 && jsonTag[0] != "" {
			log.Println("setting to jsonTag[0]")
			jsonName = jsonTag[0]
		}

		log.Printf("\nUnstringify: data before indexing: %+v\n", data)
		log.Printf("\nUnstringify: jsonName before indexing: %+v\n", jsonName)
		log.Printf("\nUnstringify: data[jsonName]: %+v\n", data[jsonName])
		if value, ok := data[jsonName]; ok {
			log.Printf("\nUnstringify:  value: %+v\n", value)
			newValue, err := convertType(f.Type, value)
			log.Printf("\nUnstringify: newValue: %+v\n", newValue)
			if err != nil {
				log.Println("\nUnstringify: err in value check\n")
				return err
			}

			val.FieldByName(f.Name).Set(newValue)
		}
	}

	log.Printf("\nUnstringify val at end: %+v\n", val)
	return nil
}
