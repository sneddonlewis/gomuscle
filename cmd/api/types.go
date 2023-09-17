package main

import (
	"math/rand"
	"time"
)

type Workout struct {
	ID   int           `json:"id"`
	Date time.Time     `json:"date"`
	Sets []*WorkingSet `json:"sets"`
}

type WorkingSet struct {
	ID                  int    `json:"id"`
	Exercise            string `json:"exercise"`
	ResistanceKg        int    `json:"resistance_kg"`
	Repetitions         int    `json:"repitions"`
	NegativeRepetitions int    `json:"negative_repitions"`
	StaticHoldSeconds   int    `json:"static_hold_seconds"`
	WorkoutID           int    `json:"workout_id"`
}

func NewWorkout() *Workout {
	return &Workout{
		Date: time.Now(),
		Sets: []*WorkingSet{},
	}
}

type TransferRequest struct {
	ToAccount int `json:"to_account"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int       `json:"number"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Intn(100000),
		CreatedAt: time.Now().UTC(),
	}
}
