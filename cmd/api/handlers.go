package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *ApiServer) health(w http.ResponseWriter, r *http.Request) error {
	healthData := map[string]string{
		"status":      "available",
		"environment": s.cfg.env,
		"version":     version,
	}

	return writeJSON(w, http.StatusOK, healthData)
}

func (s *ApiServer) workout(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.getWorkouts(w, r)
	}

	if r.Method == http.MethodPost {
		return s.createWorkout(w, r)
	}

	return fmt.Errorf("Http Request Method Type Not Supported: %s", r.Method)
}

func (s *ApiServer) createWorkout(w http.ResponseWriter, r *http.Request) error {
	workout, err := s.store.CreateWorkout()

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, workout)
}

func (s *ApiServer) getWorkouts(w http.ResponseWriter, r *http.Request) error {
	workouts, err := s.store.GetWorkouts()

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, workouts)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Http Request Method Type Not Supported: %s", r.Method)
}

// FOR REMOVAL
func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getId(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountById(id)
		if err != nil {
			return err
		}

		return writeJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}

	account := NewAccount(createAccReq.FirstName, createAccReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenStr, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println(tokenStr)

	return writeJSON(w, http.StatusCreated, account)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return writeJSON(w, http.StatusOK, transferReq)
}
