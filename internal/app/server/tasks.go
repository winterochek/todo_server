package api

import (
	"encoding/json"
	"net/http"

	"github.com/winterochek/todo-server/internal/app/model"
)

// /tasks POST controller
func (s *server) HandleTaskCreate() http.HandlerFunc {
	type request struct {
		Content   string
		Completed bool
	}
	type response struct {
		Task *model.Task `json:"task"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		userId, ok := r.Context().Value(context_userId_key).(int)
		if !ok {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		t := model.Task{
			UserID:    userId,
			Content:   req.Content,
			Completed: req.Completed,
		}

		err = t.Validate()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Task().Create(&t)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		res := &response{
			Task: &t,
		}
		s.respond(w, r, http.StatusOK, res)
	}
}
