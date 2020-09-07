package qdb

import "errors"

var (
	// ErrNoSuchQuote is returned in the event that a quote was
	// requested or configuration specified that returns no data.
	ErrNoSuchQuote = errors.New("no quote matches the given parameters")

	// ErrNoSuchBackend is returned when a QuoteDB backend is
	// requested that does not exist.
	ErrNoSuchBackend = errors.New("backend specified does not exist")

	// ErrInternal is returned for all uncategorized internal
	// database errors.
	ErrInternal = errors.New("an internal database error has occured")
)
