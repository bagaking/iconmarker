// Package filter provides image filtering capabilities for Icon Marker
package filter

import (
	"image/draw"
)

// FilterOption defines options for filtering operations
type FilterOption interface {
	// ValidateOption validates the filtering options
	ValidateOption() error
}

// Filter defines the interface for all image filters
type Filter interface {
	// Apply applies the filter to the given image
	Apply(img draw.Image, options FilterOption) error
}

// TintOption defines options for tint filter
type TintOption struct {
	// Color is the tint color
	Color [3]uint8
	// Intensity is between 0 and 1, where 0 means no effect and 1 means full tint
	Intensity float64
}

// ValidateOption validates the tint options
func (o TintOption) ValidateOption() error {
	if o.Intensity < 0 || o.Intensity > 1 {
		return ErrInvalidIntensity
	}
	return nil
}

// GrayscaleOption defines options for grayscale filter
type GrayscaleOption struct {
	// PreserveAlpha determines whether to preserve the alpha channel
	PreserveAlpha bool
}

// ValidateOption validates the grayscale options
func (o GrayscaleOption) ValidateOption() error {
	return nil // No validation needed
}

// OpacityOption defines options for opacity filter
type OpacityOption struct {
	// Opacity is between 0 and 1, where 0 means fully transparent and 1 means fully opaque
	Opacity float64
}

// ValidateOption validates the opacity options
func (o OpacityOption) ValidateOption() error {
	if o.Opacity < 0 || o.Opacity > 1 {
		return ErrInvalidOpacity
	}
	return nil
}

// InvertOption defines options for invert filter
type InvertOption struct {
	// InvertAlpha determines whether to invert the alpha channel as well
	InvertAlpha bool
}

// ValidateOption validates the invert options
func (o InvertOption) ValidateOption() error {
	return nil // No validation needed
}

// FilterManager manages and applies filters to images
type FilterManager struct {
	filters map[string]Filter
}

// NewFilterManager creates a new filter manager
func NewFilterManager() *FilterManager {
	manager := &FilterManager{
		filters: make(map[string]Filter),
	}

	// Register default filters
	manager.Register("grayscale", NewGrayscaleFilter())
	manager.Register("tint", NewTintFilter())
	manager.Register("opacity", NewOpacityFilter())
	manager.Register("invert", NewInvertFilter())

	return manager
}

// Register registers a filter with a name
func (m *FilterManager) Register(name string, filter Filter) {
	m.filters[name] = filter
}

// Get returns a filter by name
func (m *FilterManager) Get(name string) (Filter, bool) {
	filter, ok := m.filters[name]
	return filter, ok
}

// Apply applies a named filter to an image
func (m *FilterManager) Apply(img draw.Image, name string, options FilterOption) error {
	filter, ok := m.Get(name)
	if !ok {
		return ErrFilterNotFound
	}

	return filter.Apply(img, options)
}
