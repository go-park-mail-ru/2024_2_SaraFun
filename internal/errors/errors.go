package errors

import stderr "errors"

var ErrWrongCredentials = stderr.New("wrong credentials")
var ErrInvalidSession = stderr.New("invalid session")
var ErrRegistrationUser = stderr.New("user register failed")
var ErrBadUsername = stderr.New("bad username")

var ErrUserNotFound = stderr.New("user not found")
var ErrCannotLikeSelf = stderr.New("cannot like/dislike self")
