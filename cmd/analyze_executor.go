package cmd

import (
	"os"

	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/outputer"
)

type analyzeExecutor struct {
	manager         collectors_manager.CollectorManager
	analyzer        analyzers.Analyzer
	enricherManager enricher.EnricherManager
	out             outputer.Outputer
}

func initializeAnalyzeExecutor(manager collectors_manager.CollectorManager,
	analyzer analyzers.Analyzer,
	enricherManager enricher.EnricherManager,
	outputer outputer.Outputer) *analyzeExecutor {
	return &analyzeExecutor{
		manager:         manager,
		analyzer:        analyzer,
		enricherManager: enricherManager,
		out:             outputer,
	}
}

func (r *analyzeExecutor) Run() error {
	defer errlog.FlushAll()

	// let progress bar run in the background
	pWaiter := progressbar.Run()

	// start all pipeline parts in the background
	collectionChan := r.manager.Collect()
	analyzedDataChan := r.analyzer.Analyze(collectionChan)
	enrichedDataChan := r.enricherManager.Enrich(analyzedDataChan)
	outputWaiter := r.out.Digest(enrichedDataChan)

	// wait for progress bars to finish before outputting
	pWaiter.Wait()

	// wait for output to be digested
	outputWaiter.Wait()

	return r.out.Output(os.Stdout)
}
