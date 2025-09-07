package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Find(array []string, value string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}

	return false
}

// converts array time []*string into []string
func ConvertStringPointerArrayToStringArray(pointerArray []*string) []string {
	convertedArray := make([]string, len(pointerArray))

	for i, v := range pointerArray {
		if v != nil {
			convertedArray[i] = *v
		}
	}

	return convertedArray
}

// function that finds parameters from the context
func GetIDParam(ctx *gin.Context, keys ...string) (int, error) {
	for _, key := range keys {
		if val := ctx.Param(key); val != "" {
			return strconv.Atoi(val)
		}
	}

	return 0, fmt.Errorf("no valid ID param found")
}
