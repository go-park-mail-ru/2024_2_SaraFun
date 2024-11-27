package errors

import stderr "errors"

var ErrWrongCredentials = stderr.New("wrong credentials")
var ErrInvalidSession = stderr.New("invalid session")
var ErrRegistrationUser = stderr.New("user register failed")
var ErrBadUsername = stderr.New("bad username")

var ErrUserNotFound = stderr.New("user not found")
var ErrCannotLikeSelf = stderr.New("cannot reaction/dislike self")
var ErrNoResult = stderr.New("no result found")
var ErrSmallAge = stderr.New("too young")
var ErrBigAge = stderr.New("too old")
