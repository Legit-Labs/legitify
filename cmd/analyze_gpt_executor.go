package cmd

import (
	"context"
	"fmt"
	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	"github.com/Legit-Labs/legitify/internal/gpt"
	"github.com/Legit-Labs/legitify/internal/screen"
	"github.com/fatih/color"
	"strings"
)

type analyzeGPTExecutor struct {
	manager  collectors_manager.CollectorManager
	context  context.Context
	analyzer *gpt.Analyzer
}

func initializeAnalyzeGPTExecutor(analyzer *gpt.Analyzer,
	manager collectors_manager.CollectorManager,
	ctx context.Context) *analyzeGPTExecutor {
	return &analyzeGPTExecutor{
		manager:  manager,
		context:  ctx,
		analyzer: analyzer,
	}
}

func formatResults(results []gpt.Result) string {
	const title = "GPT Recommendations:"
	sb := strings.Builder{}
	for _, r := range results {
		sb.WriteString("\n")
		sb.WriteString(color.HiCyanString("%s:\n", r.EntityType))
		sb.WriteString(color.HiCyanString("%s\n", strings.Repeat("-", len(r.EntityType))))
		sb.WriteString(fmt.Sprintf("Id: %d\n", r.Entity.ID()))
		sb.WriteString(fmt.Sprintf("Name: %s\n", r.Entity.Name()))
		sb.WriteString(fmt.Sprintf("Url: %s\n", r.Entity.CanonicalLink()))
		sb.WriteString("\n")
		sb.WriteString(color.MagentaString("%s\n", title))
		sb.WriteString(color.MagentaString("%s", strings.Repeat("-", len(title))))
		sb.WriteString("\n")
		sb.WriteString(r.GPTResult)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r *analyzeGPTExecutor) Run() error {
	pWaiter := progressbar.Run()
	collected := r.manager.Collect()
	var results []gpt.Result
	for res := range r.analyzer.Analyze(collected) {
		results = append(results, res)
	}

	pWaiter.Wait()
	screen.Printf("%s\n", formatResults(results))
	return nil
}
