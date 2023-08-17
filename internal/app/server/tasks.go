package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/helpers"
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
		// get userId from context, provided by AuthMiddleware
		userId, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		// decode request body
		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// create task model
		t := model.Task{
			UserID:    userId,
			Content:   req.Content,
			Completed: req.Completed,
		}

		// validate task model
		err = t.Validate()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// create storage record. task model with updated by new values
		err = s.store.Task().Create(&t)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// create response model
		res := &response{
			Task: &t,
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

// tasks GET controller for single task
func (s *server) HandleSingleTaskRead() http.HandlerFunc {
	type response struct {
		Task *model.Task `json:"task"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		taskID, err := h.GetTaskIDFromParams(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		t, err := s.store.Task().ReadOne(taskID, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusNotFound, ErrNotFound)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		res := &response{
			Task: t,
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

// tasks GET controller for multiple tasks, that belong to one user
func (s *server) HandleMultipleTasksRead() http.HandlerFunc {
	type response struct {
		Tasks []*model.Task `json:"tasks"`
		Count int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// get userId from context, provided by AuthMiddleware
		userId, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		tasks, err := s.store.Task().ReadAll(userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		count := len(tasks)
		if count == 0 {
			tasks = []*model.Task{}
		}
		res := &response{
			Tasks: tasks,
			Count: count,
		}
		s.respond(w, r, http.StatusOK, res)
	}

}

// tasks PUT controller. Created for updating the whole task info, not only to complete
func (s *server) HandleTasksUpdate() http.HandlerFunc {
	type request struct {
		Content   string
		Completed bool
	}

	type response struct {
		Task *model.Task `json:"task"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		// get userId from context, provided by AuthMiddleware
		userId, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		// get taskId from request params
		taskID, err := h.GetTaskIDFromParams(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// decode request body
		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// create task model
		t := model.Task{
			UserID:    userId,
			Content:   req.Content,
			Completed: req.Completed,
			ID:        taskID,
		}
		// validation task model => ensure that content is provided. If completed is not provided - it will be set as false
		err = t.Validate()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		// update record and get updated model back
		err = s.store.Task().Update(&t)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		// create response model
		res := &response{
			Task: &t,
		}

		s.respond(w, r, http.StatusOK, res)
	}
}

// controller for /tasks DELETE
func (s *server) HandleTasksDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		taskId, err := h.GetTaskIDFromParams(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		err = s.store.Task().Delete(taskId, userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}
