package enrichers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

const MembersList = "violatedUsers"

func NewMembersListEnricher(ctx context.Context) Enricher {
	return &membersListEnricher{}
}

type membersListEnricher struct {
}

func (e *membersListEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createMembersListEnrichment(data.ExtraData)
	if err != nil {
		return nil, false
	}
	return result, true
}

func createMembersListEnrichment(extraData interface{}) (Enrichment, error) {
	casted, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid membersList extra data")
	}

	var result []githubcollected.OrganizationMember
	for k := range casted {
		var member githubcollected.OrganizationMember
		err := json.Unmarshal([]byte(k), &member)
		if err != nil {
			return nil, err
		}

		result = append(result, member)
	}

	return &MembersListEnrichment{
		Members: result,
	}, nil
}

func (e *membersListEnricher) Name() string {
	return MembersList
}

type MembersListEnrichment struct {
	Members []githubcollected.OrganizationMember
}

func (se *MembersListEnrichment) Name() string {
	return MembersList
}

func (se *MembersListEnrichment) HumanReadable(prepend string, linebreak string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, member := range se.Members {
		sb.WriteString(fmt.Sprintf("%d. %s (%d)%s", i+1, *member.User.HTMLURL, *member.User.ID, linebreak))
	}

	return sb.String()
}
