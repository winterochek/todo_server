package sqlstore_test

import (
	"os"
	"testing"
)

var dbURI string

func TestMain(m *testing.M){
	dbURI = "host=localhost user=admin password=admin dbname=test sslmode=disable"
	os.Exit(m.Run())
}