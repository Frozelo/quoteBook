package appErrors

import "errors"

var ErrQuoteNotFound = errors.New("quote not found")
var ErrNoQuotes = errors.New("no quotes available")
