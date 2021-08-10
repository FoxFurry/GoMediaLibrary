package db

import (
	"fmt"
	"github.com/foxfurry/simple-rest/app/env_setup"
	"testing"
)

func init(){
	env_setup.TestConfig()
}

func TestBookDBRepository_GetBook(t *testing.T) {
	fmt.Printf("Running test 1")
}