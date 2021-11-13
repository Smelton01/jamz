package cmd

import "errors"

var ErrInvalidQuery = errors.New("search query is not valid")
var ErrNoActiveDevice = errors.New("no active device detected")
