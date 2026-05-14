package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"myapp/model"
	"myapp/utils/httpResp"

	"github.com/gorilla/mux"
)

func getUserId(userIdParam string) (int64, error) {
	return strconv.ParseInt(userIdParam, 10, 64)
}

func AddStudent(w http.ResponseWriter, r *http.Request) {
	var s model.Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := s.Create(); err != nil {
		if errors.Is(err, model.ErrFirstNameRequired) || errors.Is(err, model.ErrEmailRequired) {
			httpResp.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusCreated, s)
}

func GetStud(w http.ResponseWriter, r *http.Request) {
	id, err := getUserId(mux.Vars(r)["sid"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}
	s := model.Student{StdId: id}
	if err := s.Read(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpResp.RespondWithError(w, http.StatusNotFound, "Student not found")
			return
		}
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, s)
}

func GetAllStuds(w http.ResponseWriter, r *http.Request) {
	students, err := model.GetAllStudents()
	if err != nil {
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, students)
}

func UpdateStud(w http.ResponseWriter, r *http.Request) {
	oldID, err := getUserId(mux.Vars(r)["sid"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	var s model.Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if oldID != 0 {
		s.StdId = oldID
	}
	if err := s.Update(oldID); err != nil {
		if errors.Is(err, model.ErrFirstNameRequired) || errors.Is(err, model.ErrEmailRequired) || errors.Is(err, model.ErrStudentIDRequired) {
			httpResp.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			httpResp.RespondWithError(w, http.StatusNotFound, "Student not found")
			return
		}
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, s)
}

func DeleteStud(w http.ResponseWriter, r *http.Request) {
	id, err := getUserId(mux.Vars(r)["sid"])
	if err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}
	s := model.Student{StdId: id}
	if err := s.Delete(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpResp.RespondWithError(w, http.StatusNotFound, "Student not found")
			return
		}
		httpResp.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
