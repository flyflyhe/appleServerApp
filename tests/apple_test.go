package tests

import (
	"fmt"
	"testing"

	"github.com/flyflyhe/appleServerApp/services/apple"
)

func TestGetTransactionHistory(t *testing.T) {
	result, err := apple.GetTransactionHistory("420000546563865", apple.GetAppleJwtToken(), "")
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(result)
}
