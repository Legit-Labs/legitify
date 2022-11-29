package opa_test

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"log"
	"testing"

	"github.com/Legit-Labs/legitify/internal/opa"
)

func TestEngineSanity(t *testing.T) {
	ctx := context.Background()
	// Doesn't matter which scm type we use here
	engine, err := opa.Load([]string{"./testdata"}, scm_type.GitHub)

	if err != nil {
		t.Errorf("Unable to engine with policies")
	}

	engine.SetTracing(true)

	input := map[string]interface{}{
		"bla": "o2k",
	}

	result, err := engine.Query(ctx, "test", input)

	if err != nil {
		t.Errorf("Failed to query engine: %s", err)
	} else if len(result) != 4 {
		log.Println(result)
	}
}
