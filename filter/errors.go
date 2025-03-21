package filter

import "errors"

// Common errors for the filter package
var (
	ErrFilterNotFound        = errors.New("filter not found")
	ErrInvalidIntensity      = errors.New("intensity must be between 0 and 1")
	ErrInvalidOpacity        = errors.New("opacity must be between 0 and 1")
	ErrInvalidColor          = errors.New("invalid color")
	ErrNoFiltersSpecified    = errors.New("no filters specified")
	ErrFilterOptionsMismatch = errors.New("number of filter options must match number of filters")
)
