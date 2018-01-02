package performer

import (
	"fmt"
	"github.com/sankarvj/syncadapter/core"
	"reflect"
)

func inImplementsCooker(in interface{}) bool {
	cookerin := reflect.TypeOf((*core.Cooker)(nil)).Elem()
	if reflect.TypeOf(in).Implements(cookerin) {
		return true
	} else {
		return false
	}
}

func inImplementsPasser(in interface{}) bool {
	passerin := reflect.TypeOf((*core.Passer)(nil)).Elem()
	if reflect.TypeOf(in).Implements(passerin) {
		return true
	} else {
		return false
	}
}

func PasserSlice(slice interface{}) (ret1 []core.Passer) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret1 = make([]core.Passer, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret1[i] = s.Index(i).Interface().(core.Passer)
	}
	return ret1
}

func CookerSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func needUpdate(serverupdated int64, localupdated int64) bool {
	if localupdated == serverupdated {
		return false
	} else {
		return true
	}
}

type SyncError struct {
	What string
	Stop bool
}

func (e SyncError) Error() string {
	return fmt.Sprintf("%v: %v", e.What, e.Stop)
}

func oops(errstr string, shouldstop bool) error {
	return SyncError{
		errstr,
		shouldstop,
	}
}
