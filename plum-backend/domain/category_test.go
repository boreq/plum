package domain

import (
	"testing"
)

func TestCategoriesArePresentWithoutData(t *testing.T) {
	summary := NewSummary()

	for _, category := range Categories {
		if _, ok := summary.Categories[category]; !ok {
			t.Errorf("missing category %q", category)
		}
	}
}
