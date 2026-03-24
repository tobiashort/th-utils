package cosine_similarity_test

import (
	"fmt"
	"testing"

	. "github.com/tobiashort/th-utils/lib/cosine-similarity"
)

func Test(t *testing.T) {
	s1 := "the best data science course"
	s2 := "data science is popular"
	sim := CosineSimilarity(s1, s2)
	fmt.Printf("[%s]<>[%s] = %f\n", s1, s2, sim)
	if sim < 0.43 || sim > 0.45 {
		t.Errorf("Expected around 0.44±0.01, got %f", sim)
	}
}
