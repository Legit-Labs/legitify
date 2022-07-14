package scorecard

import (
	"context"
	"github.com/ossf/scorecard/v4/checker"
	docs "github.com/ossf/scorecard/v4/docs/checks"
	sclog "github.com/ossf/scorecard/v4/log"
	"github.com/ossf/scorecard/v4/pkg"
	"github.com/ossf/scorecard/v4/policy"
	"log"
	"os"
)

func init() {
	// needed to enable webhook checks
	_ = os.Setenv("SCORECARD_V6", "")
}

type Result struct {
	Score  float64             `json:"score"`
	Result pkg.ScorecardResult `json:"result"`
}

func Calculate(ctx context.Context, repoUrl string, isPrivate bool) (*Result, error) {
	logger := sclog.NewLogger(sclog.DebugLevel)
	repo, repoClient, fuzzClient, ciiClient, vulnClient, err := checker.GetClients(ctx, repoUrl, "", logger)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = repoClient.Close()
		if err != nil {
			log.Printf("Failed to close repository client %s", err)
		}
		err = fuzzClient.Close()
		if err != nil {
			log.Printf("Failed to close fuzz client %s", err)
		}
	}()

	checks := []string{
		"Binary-Artifacts",
		"Branch-Protection",
		"Code-Review",
		"Contributors",
		"Dangerous-Workflow",
		"Dependency-Update-Tool",
		"Maintained",
		"Pinned-Dependencies",
		"SAST",
		"Token-Permissions",
		"Vulnerabilities",
		"Webhooks",
	}

	if !isPrivate {
		checks = append(checks, []string{
			"Packaging",
			"Security-Policy",
			"CII-Best-Practices",
			"Fuzzing",
			"License",
			"Signed-Releases",
		}...)
	}

	enabledChecks, err := policy.GetEnabled(nil, checks, nil)
	if err != nil {
		return nil, err
	}

	d, err := docs.Read()
	if err != nil {
		return nil, err
	}

	repoResult, err := pkg.RunScorecards(
		ctx,
		repo,
		"HEAD",
		enabledChecks,
		repoClient,
		fuzzClient,
		ciiClient,
		vulnClient,
	)
	if err != nil {
		return nil, err
	}

	score, err := repoResult.GetAggregateScore(d)
	if err != nil {
		return nil, err
	}

	return &Result{
		Score:  score,
		Result: repoResult,
	}, nil
}
