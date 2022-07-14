package analyzers

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/opa"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Legit-Labs/legitify/internal/collectors"
)

func TestAnalyzerSanity(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	ctx = context_utils.NewContextWithTokenScopes(ctx, permissions.TokenScopes{})
	data := make(chan collectors.CollectedData, 3)

	engine, _ := opa.Load([]string{})
	analyzer := NewAnalyzer(ctx, engine, skippers.NewSkipper(ctx))
	require.NotNilf(t, analyzer, "failed to create analyzer")

	type nullEntity struct {
		githubcollected.CollectedEntity
	}

	someData := collectors.CollectedData{
		Namespace: "test",
		Entity:    nullEntity{},
	}

	data <- someData
	data <- someData
	data <- someData
	close(data)

	// Run
	analyzer.Analyze(data)
}
