package oops

import "errors"

var InvalidKey = errors.New("invalid key")
var KeyNotFound = errors.New("key not found")
var KeyAlreadyExists = errors.New("key already exists")
var ValidationError = errors.New("validation error")
