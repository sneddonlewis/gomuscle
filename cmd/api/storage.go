package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateWorkout() (*Workout, error)
	GetWorkouts() ([]*Workout, error)

	// FOR REMOVAL
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(user, dbName, password string) (*PostgresStore, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbName, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return nil
}

func (s *PostgresStore) CreateWorkout() (*Workout, error) {
	workout := NewWorkout()
	query := `
		INSERT INTO Workout
		(date) VALUES ($1) RETURNING id
	`
	var workoutID int
	err := s.db.QueryRow(query, workout.Date).Scan(&workoutID)
	if err != nil {
		return nil, err
	}

	workout.ID = workoutID
	return workout, nil
}

func (s *PostgresStore) GetWorkouts() ([]*Workout, error) {
	workoutQuery := "SELECT * FROM Workout"

	rows, err := s.db.Query(workoutQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []*Workout{}

	for rows.Next() {
		workout, err := scanIntoWorkout(rows)
		if err != nil {
			return nil, err
		}

		workingSets, err := s.GetWorkingSets(workout.ID)

		if err != nil {
			return nil, err
		}
		workout.Sets = workingSets

		workouts = append(workouts, workout)
	}

	return workouts, nil
}

func (s *PostgresStore) GetWorkingSets(workoutID int) ([]*WorkingSet, error) {
	setsQuery := "SELECT * FROM working_sets WHERE workout_id = $1"
	setsRows, err := s.db.Query(setsQuery, workoutID)
	if err != nil {
		return nil, err
	}
	defer setsRows.Close()

	workingSets := []*WorkingSet{}

	for setsRows.Next() {
		workingSet, err := scanIntoWorkingSet(setsRows)
		if err != nil {
			return nil, err
		}
		workingSets = append(workingSets, workingSet)
	}

	return workingSets, nil
}

func scanIntoWorkingSet(rows *sql.Rows) (*WorkingSet, error) {
	workingSet := new(WorkingSet)
	err := rows.Scan(
		&workingSet.ID,
		&workingSet.Exercise,
		&workingSet.ResistanceKg,
		&workingSet.Repetitions,
		&workingSet.NegativeRepetitions,
		&workingSet.StaticHoldSeconds,
		&workingSet.WorkoutID,
	)
	return workingSet, err
}

func scanIntoWorkout(rows *sql.Rows) (*Workout, error) {
	workout := new(Workout)
	err := rows.Scan(
		&workout.ID,
		&workout.Date,
	)
	return workout, err
}

// OLD ACCOUNT STUFF FOR REMOVAL

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists Account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		insert into Account 
		(first_name, last_name, number, balance, created_at)
		values
		($1, $2, $3, $4, $5)
	`
	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt)

	return err
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := "select * from account where id = $1"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := "select * from Account"

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := "delete from Account where id = $1"
	_, err := s.db.Query(query, id)
	return err
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	if err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt); err != nil {
		return nil, err
	}

	return account, nil
}
