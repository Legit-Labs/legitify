package cmd

import (
	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/outputer"
	"log"
	"os"
)

type analyzeExecutor struct {
	manager         collectors_manager.CollectorManager
	analyzer        analyzers.Analyzer
	enricherManager enricher.EnricherManager
	out             outputer.Outputer
	log             *log.Logger
}

func initializeAnalyzeExecutor(manager collectors_manager.CollectorManager,
	analyzer analyzers.Analyzer,
	enricherManager enricher.EnricherManager,
	outputer outputer.Outputer,
	log *log.Logger) *analyzeExecutor {
	return &analyzeExecutor{
		manager:         manager,
		analyzer:        analyzer,
		enricherManager: enricherManager,
		out:             outputer,
		log:             log,
	}
}

func (r *analyzeExecutor) Run() error {
	r.log.Printf("Gathering collection metadata...")
	collectionMetadata := r.manager.CollectMetadata()
	progressBar := progressbar.NewProgressBar(collectionMetadata)

	// TODO progressBar should run before collection starts and wait for channels to read from
	collectionChannels := r.manager.Collect()
	pWaiter := progressBar.Run(collectionChannels.Progress)
	analyzedDataChan := r.analyzer.Analyze(collectionChannels.Collected)
	enrichedDataChan := r.enricherManager.Enrich(analyzedDataChan)
	outputWaiter := r.out.Digest(enrichedDataChan)

	// Wait for progress bars to finish before outputting
	pWaiter.Wait()

	// Wait for output to be digested
	outputWaiter.Wait()

	return r.out.Output(os.Stdout)
}
