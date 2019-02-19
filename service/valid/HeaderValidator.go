package valid

import (
	"net/http"
	"errors"
)

type HttpHeaderValidator struct {
	Header *http.Header
}


func (that *HttpHeaderValidator) valid(schemeField string, target string) error {
	srcVal := that.Header.Get(schemeField)
	if target == "" {
		return nil
	}

	if srcVal == "" {
		return errors.New(schemeField + ":is nil or empty")
	}

	if srcVal == target {
		return nil
	}
	return errors.New("value not equal:[" + srcVal + "]=[" + target + "]")
}