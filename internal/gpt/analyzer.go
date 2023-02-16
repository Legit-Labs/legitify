package gpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/collected"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"log"
	"strings"
)

type Analyzer struct {
	context   context.Context
	gptClient *gogpt.Client
}

type Result struct {
	Entity     collected.Entity
	EntityType string
	GPTResult  string
}

func NewAnalyzer(ctx context.Context, gptToken string) *Analyzer {
	return &Analyzer{
		context:   ctx,
		gptClient: gogpt.NewClient(gptToken),
	}
}

func cleanData(entity collected.Entity) (marshalled []byte, entityType string, err error) {
	switch v := entity.(type) {
	case ghcollected.Organization:
		entityType = "Github Organization"
		marshalled, err = json.Marshal(v)
	case ghcollected.Repository:
		entityType = "Github Repository"
		v.Collaborators = nil
		marshalled, err = json.Marshal(v)
	case gitlab_collected.Organization:
		entityType = "Gitlab Organization"
		v.Projects = nil
		marshalled, err = json.Marshal(v)
	case gitlab_collected.Repository:
		entityType = "Gitlab Repository"
		v.Members = nil
		marshalled, err = json.Marshal(v)
	default:
		err = fmt.Errorf("unknow about type %T!\n", v)
		return
	}

	return
}

func generatePrompt(toAnalyze []byte, entityType string) string {
	prompt := fmt.Sprintf("Explain the security posture of the below %s, provide the answer as a list of recommendations, all must be relevant to the provided configuration (atleast 10),"+
		" recommendations needs to be related to the provided data and start the text with the numbered list:\n", entityType)
	return fmt.Sprintf("%s\n%s\n", prompt, toAnalyze)
}

func streamResults(stream *gogpt.CompletionStream, barName string) (string, error) {
	progressbar.Report(progressbar.OptionalBarCreation{
		BarName:       barName,
		TotalEntities: 0,
	})

	aggregate := strings.Builder{}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			progressbar.Report(progressbar.BarUpdate{
				BarName:     barName,
				TriggerDone: true,
			})
			break
		}

		if err != nil {
			progressbar.Report(progressbar.BarClose{
				BarName: barName,
			})
			return "", err
		}

		if len(response.Choices) > 0 {
			progressbar.Report(progressbar.BarUpdate{
				BarName:     barName,
				Change:      1,
				TotalChange: 4,
			})
			aggregate.WriteString(response.Choices[0].Text)
		}
	}

	progressbar.Report(progressbar.BarClose{
		BarName: barName,
	})
	return aggregate.String(), nil
}

func (a *Analyzer) Analyze(dataChannel <-chan collectors.CollectedData) chan Result {
	result := make(chan Result)

	go func() {
		defer close(result)
		gw := group_waiter.New()
		for data := range dataChannel {
			data := data
			gw.Do(func() {
				raw, entityType, err := cleanData(data.Entity)
				if err != nil {
					log.Println(err)
					return
				}

				prompt := generatePrompt(raw, entityType)

				stream, err := a.gptClient.CreateCompletionStream(a.context, gogpt.CompletionRequest{
					Model:       gogpt.GPT3TextDavinci003,
					Prompt:      prompt,
					Temperature: 1.0,
					MaxTokens:   1000,
				})

				if err != nil {
					log.Println(err)
					return
				}

				gptResult, err := streamResults(stream,
					fmt.Sprintf("GPT generating tokens for %s", data.Entity.Name()))

				if err != nil {
					log.Println(err)
					return
				}

				result <- Result{
					Entity:     data.Entity,
					GPTResult:  gptResult,
					EntityType: entityType,
				}
			})
		}
		gw.Wait()
	}()

	return result
}
