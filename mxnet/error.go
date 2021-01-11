package mxnet

/*
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <mxnet/c_predict_api.h>
*/
import "C"
import (
	"github.com/pkg/errors"
)

// get the last error happeneed.
// go binding for MXGetLastError
func GetLastError() error {
	if err := C.MXGetLastError(); err != nil {
		return errors.Errorf("error in mxnet :: %v", C.GoString(err))
	}
	return nil
}
