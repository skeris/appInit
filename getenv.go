package appInit

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	ENV     = "env"
	DEFAULT = "default"
)

func getEnv(mask interface{}) interface{} {
	r := reflect.ValueOf(mask)

	var argTypeRV reflect.Value
	argTypeRV = reflect.New(r.Type())

	for i := 0; i < r.NumField(); i++ {
		v := r.Type().Field(i).Tag.Get(ENV)
		d := r.Type().Field(i).Tag.Get(DEFAULT)

		env, ok := os.LookupEnv(v)
		var val string
		if ok {
			val = env
		} else {
			val = d
		}
		fmt.Println("ORA", v, d)
		fmt.Println("ORA1", reflect.ValueOf(&mask))
		fmt.Println("ORA12", reflect.ValueOf(mask))

		switch t := r.Type().Field(i).Type.Name(); t {
		case "string":
			argTypeRV.Field(i).SetString(val)
		case "bool":
			if strings.ToLower(val) == "true" {
				reflect.ValueOf(&mask).Elem().Field(i).SetBool(true)
			} else {
				reflect.ValueOf(&mask).Elem().Field(i).SetBool(false)
			}
		case "uint64":
			num, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				panic(err)
			}
			reflect.ValueOf(&mask).Elem().Field(i).SetUint(num)
		default:
			panic(errors.New("Something strange happend: " + t))
		}
	}

	return mask
}
