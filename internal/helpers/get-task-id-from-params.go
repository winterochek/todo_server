package h

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTaskIDFromParams(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	id := int(taskID)
	return id, nil
}
