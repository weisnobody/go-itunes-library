package itunes

import "fmt"

// ErrInvalidFormat is raised when an input file
// does not meet the expected itune library xml format
var ErrInvalidFormat = fmt.Errorf(
    "The given library file is not in the expected iTunes XML format")
