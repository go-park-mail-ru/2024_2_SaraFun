package errors

import stderr "errors"

var ErrWrongCredentials = stderr.New("wrong credentials")
