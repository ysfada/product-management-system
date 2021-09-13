package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppErrWithoutDetail(t *testing.T) {
	err := AppErr{
		Message: "error",
	}

	assert.NotNil(t, err)
	var appErr AppErr
	assert.IsType(t, appErr, err)
	assert.Nil(t, err.Detail)
	// json
	assert.Equal(t, "{\"Message\":\"error\",\"Detail\":null}", err.Error())
}

func TestAppErrWithDetail(t *testing.T) {
	err := AppErr{
		Message: "error",
		Detail:  "detail",
	}

	assert.NotNil(t, err)
	var appErr AppErr
	assert.IsType(t, appErr, err)
	assert.NotNil(t, err.Detail)
	// json
	assert.Equal(t, "{\"Message\":\"error\",\"Detail\":\"detail\"}", err.Error())
}
