package models

type DBEvent int

const (
	DELETE DBEvent = iota
	INSERT
	UPDATE
)
