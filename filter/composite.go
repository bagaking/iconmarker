package filter

import (
	"image/draw"
)

// CompositeFilter applies multiple filters in sequence
type CompositeFilter struct {
	filters []Filter
}

// CompositeOption defines options for the composite filter
type CompositeOption struct {
	Filters     []Filter       // Filters to apply
	Options     []FilterOption // Options corresponding to each filter
	StopOnError bool           // Whether to stop on first error
}

// ValidateOption validates the composite filter options
func (o CompositeOption) ValidateOption() error {
	if len(o.Filters) == 0 {
		return ErrNoFiltersSpecified
	}

	if len(o.Options) > 0 && len(o.Options) != len(o.Filters) {
		return ErrFilterOptionsMismatch
	}

	return nil
}

// NewCompositeFilter creates a new composite filter with the given filters
func NewCompositeFilter(filters ...Filter) *CompositeFilter {
	return &CompositeFilter{
		filters: filters,
	}
}

// Apply applies all filters in sequence
func (f *CompositeFilter) Apply(img draw.Image, options FilterOption) error {
	opt, ok := options.(CompositeOption)
	if !ok {
		// If no options provided, just apply filters with nil options
		for _, filter := range f.filters {
			if err := filter.Apply(img, nil); err != nil {
				return err
			}
		}
		return nil
	}

	// Validate options
	if err := opt.ValidateOption(); err != nil {
		return err
	}

	// Apply each filter with its corresponding option if available
	for i, filter := range opt.Filters {
		var option FilterOption
		if i < len(opt.Options) {
			option = opt.Options[i]
		}

		if err := filter.Apply(img, option); err != nil {
			if opt.StopOnError {
				return err
			}
			// Log error but continue if StopOnError is false
			// In a real implementation, you might want to use a logger
		}
	}

	return nil
}
