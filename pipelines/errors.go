/*
Copyright 2022 John Homan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipelines

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	status string
	code   int
}

func (a *ApiError) Error() string {
	return fmt.Sprintf("status=%s,code=%d", a.status, a.code)
}

func (a *ApiError) Code() int {
	return a.code
}

func NewConflict() *ApiError {
	return &ApiError{code: http.StatusConflict, status: "Conflict"}
}

func IsConflict(err error) bool {
	_, ok := err.(*ApiError)
	if !ok {
		return false
	}
	return err.(*ApiError).Code() == http.StatusConflict
}

func NewNotFound() *ApiError {
	return &ApiError{code: http.StatusNotFound, status: "NotFound"}
}

func IsNotFound(err error) bool {
	_, ok := err.(*ApiError)
	if !ok {
		return false
	}
	return err.(*ApiError).Code() == http.StatusNotFound
}
