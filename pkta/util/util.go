package util

import "reflect"

/*
 * This func is for debug purpose.  It shows the type of the given var.
 * This helps getting more useful info how the var can be utilized.
 */
func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
