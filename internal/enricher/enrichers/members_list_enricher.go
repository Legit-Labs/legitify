package enrichers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

const MembersList = "violatedUsers"

func NewMembersListEnricher() membersListEnricher {
	return membersListEnricher{}
}

type membersListEnricher struct {
}

func (e membersListEnricher) Enrich(_ context.Context, data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createMembersListEnrichment(data.ExtraData)
	if err != nil {
		return nil, false
	}
	return result, true
}

func (se membersListEnricher) Parse(data interface{}) (Enrichment, error) {
	if val, ok := data.([]githubcollected.OrganizationMember); !ok {
		return nil, fmt.Errorf("expecting []githubcollected.OrganizationMember, found %T", data)
	} else {
		return MembersListEnrichment(val), nil
	}
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

	// sort by user id to maintain a determenstic order
	sort.Slice(result, func(i, j int) bool {
		return *result[i].User.ID < *result[j].User.ID
	})

	return MembersListEnrichment(result), nil
}

type MembersListEnrichment []githubcollected.OrganizationMember

func (se MembersListEnrichment) HumanReadable(prepend string, linebreak string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, member := range []githubcollected.OrganizationMember(se) {
		sb.WriteStringf("%d. %s (%d)%s", i+1, *member.User.HTMLURL, *member.User.ID, linebreak)
	}

	return sb.String()
}
