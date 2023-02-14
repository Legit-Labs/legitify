package outputer

import (
	"context"
	"io"

	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
)

type Outputer interface {
	Digest(inputChannel <-chan enricher.EnrichedData) group_waiter.Waitable
	Output(writer io.Writer) error
}

func NewOutputer(ctx context.Context, format formatter.FormatName, schemeType scheme.SchemeType, failedOnly bool) Outputer {
	return &outputer{
		format:     format,
		schemeType: schemeType,
		failedOnly: failedOnly,
	}
}

// -----------------------------------------------------------------------------

type outputer struct {
	format     formatter.FormatName
	schemeType scheme.SchemeType
	failedOnly bool
	output     []byte
	err        error
}

func enrichedDataToPolicyInfo(enrichedData enricher.EnrichedData) scheme.PolicyInfo {
	return scheme.PolicyInfo{
		Title:                    enrichedData.Title,
		Description:              enrichedData.Description,
		PolicyName:               enrichedData.PolicyName,
		FullyQualifiedPolicyName: enrichedData.FullyQualifiedPolicyName,
		Severity:                 enrichedData.Severity,
		Threat:                   enrichedData.Threat,
		RemediationSteps:         enrichedData.RemediationSteps,
		Namespace:                enrichedData.Namespace,
	}
}

func enrichedDataToViolation(enrichedData enricher.EnrichedData) scheme.Violation {
	return scheme.Violation{
		CanonicalLink:       enrichedData.CanonicalLink,
		ViolationEntityType: enrichedData.Entity.ViolationEntityType(),
		Aux:                 map_utils.ToKeySortedMap(enrichedData.Enrichers),
		Status:              enrichedData.Status,
	}
}

func (o *outputer) receiveViolations(inputChannel <-chan enricher.EnrichedData) *scheme.Flattened {
	violations := scheme.NewFlattenedScheme()

	for encrichedData := range inputChannel {
		policyName := encrichedData.FullyQualifiedPolicyName

		if _, ok := violations.AsOrderedMap().Get(policyName); !ok {
			violations.AsOrderedMap().Set(policyName, scheme.NewOutputData(enrichedDataToPolicyInfo(encrichedData)))
		}
		preAppend := violations.GetPolicyData(policyName)

		violation := enrichedDataToViolation(encrichedData)
		violations.AsOrderedMap().Set(policyName, scheme.AppendViolations(preAppend, violation))
	}

	return violations
}

func (o *outputer) Digest(inputChannel <-chan enricher.EnrichedData) group_waiter.Waitable {
	gw := group_waiter.New()

	gw.Do(func() {
		o.err = nil // zero err to allow reuse of the object
		violations := o.receiveViolations(inputChannel)
		sorted := violations.SortedBySeverity()

		if o.failedOnly {
			sorted = sorted.OnlyFailedViolations()
		}

		converted, err := converter.Convert(o.schemeType, sorted)
		if err != nil {
			o.err = err
			return
		}

		o.output, o.err = formatter.Format(o.format, formatter.DefaultOutputIndent, converted, o.failedOnly)
	})

	return gw
}

func (o *outputer) Output(writer io.Writer) error {
	if o.err != nil {
		return o.err
	}

	_, err := writer.Write(o.output)
	if err != nil {
		return err
	}

	return nil
}
