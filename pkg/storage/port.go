package storage

import "io"

// Port signature for storage pkg.
type Port interface {
	// Store save given reader and use given string as the id of that reader.
	// Depends on the implementation this may be a file name or some uuid. This
	// Store will be responsible for closing the reader.
	Store(io.ReadCloser, string)
	// Remove do remove given file id. Depends on the implementation this may
	// be a file path and filename or some uuid.
	Remove(string)
}
