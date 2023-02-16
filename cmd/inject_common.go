//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/google/wire"
)

var analyzeProviderSet = wire.NewSet(
	provideOpa,
	provideOutputer,
	provideContext,
	analyzers.NewAnalyzer,
	provideGPTAnalyzer,
	skippers.NewSkipper,
	enricher.NewEnricherManager,
	collectors_manager.NewCollectorsManager,
	initializeAnalyzeExecutor,
	initializeAnalyzeGPTExecutor,
)
