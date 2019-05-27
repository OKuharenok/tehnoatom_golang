package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	if reflect.ValueOf(out).Type().Kind() != reflect.Ptr {
		return fmt.Errorf("invalid")
	}
	res := reflect.ValueOf(&out)
	res = res.Elem()
	outType := reflect.ValueOf(res.Interface()).Type()
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}
	tmp := reflect.ValueOf(out)
	if tmp.Kind() == reflect.Ptr {
		tmp = tmp.Elem()
	}
	tmp.Convert(outType)
	if tmp.Type().Kind() == reflect.Slice {
		tmp.Set(reflect.MakeSlice(tmp.Type(), reflect.ValueOf(data).Len(), reflect.ValueOf(data).Cap()))
		for j:= 0; j < tmp.Len(); j++ {
			copy := reflect.New(tmp.Index(j).Type())
			err := i2s(reflect.ValueOf(data).Index(j).Interface(), copy.Interface())
			if err != nil {
				return err
			}
			tmp.Index(j).Set(copy.Elem())
		}
		return nil
	}
	if reflect.ValueOf(data).Type().Kind() != reflect.Map {
		return fmt.Errorf("invalid")
	}
	for i := 0; i < outType.NumField(); i++ {
		switch tmp.Field(i).Type().Kind() {
		case reflect.Int, reflect.String, reflect.Bool:
			if reflect.ValueOf(data).Type().Kind() != reflect.Map {
				return fmt.Errorf("invalid")
			}
			value := reflect.ValueOf(data).MapIndex(reflect.ValueOf(outType.Field(i).Name)).Interface()
			if reflect.ValueOf(value).Type().Kind() == reflect.Float64 {
				value = int(value.(float64))
			}
			if tmp.Field(i).Type().Kind() != reflect.ValueOf(value).Type().Kind() {
				return fmt.Errorf("invalid")
			}
			tmp.Field(i).Set(reflect.ValueOf(value))
		case reflect.Slice:
			if reflect.ValueOf(data).Type().Kind() != reflect.Map {
				return fmt.Errorf("invalid")
			}
			value := reflect.ValueOf(data).MapIndex(reflect.ValueOf(outType.Field(i).Name)).Interface()
			if reflect.ValueOf(value).Type().Kind() != reflect.Slice {
				return fmt.Errorf("invalid")
			}
			tmp.Field(i).Set(reflect.MakeSlice(tmp.Field(i).Type(), reflect.ValueOf(value).Len(), reflect.ValueOf(value).Cap()))
			for j := 0; j < reflect.ValueOf(value).Len(); j++ {
				reflect.ValueOf(value).Convert(reflect.ValueOf(value).Type())
				copy := reflect.New(tmp.Field(i).Index(j).Type())
				err := i2s(reflect.ValueOf(value).Index(j).Interface(), copy.Interface())
				if err != nil {
					return err
				}
				tmp.Field(i).Index(j).Set(copy.Elem())
			}
		case reflect.Struct:
			value := reflect.ValueOf(data).MapIndex(reflect.ValueOf(outType.Field(i).Name)).Interface()
			copy := reflect.New(tmp.Field(i).Type())
			err := i2s(value, copy.Interface())
			if err != nil {
				return err
			}
			tmp.Field(i).Set(copy.Elem())
		default:
			return fmt.Errorf("invalid")
		}
	}
	return nil
}
