package models

import "errors"

var (
ErrNotFound = errors.New("Requested item is not found!")
ErrEmailTaken = errors.New("Email Already Taken")
)
