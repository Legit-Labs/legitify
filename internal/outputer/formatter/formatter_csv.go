package formatter

import (
	"encoding/csv"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type CsvFormatter struct {
	colorizer humanColorizer
}

func newCSVFormatter() OutputFormatter {
	return &CsvFormatter{
		colorizer: humanColorizer{},
	}
}

func (f *CsvFormatter) csvFailedPolicies(output *scheme.Flattened, csvwriter *csv.Writer) bool {
	failedPolicies := output.OnlyFailedViolations()
	headers := []string{"#", "Policy Name", "Namespace", "Severity", "Threat", "Violations", "Remediation Steps"}
	err := csvwriter.Write(headers)
	
	
	for i, policyName := range failedPolicies.AsOrderedMap().Keys() {
		policyData := output.GetPolicyData(policyName)
		policyInfo := policyData.PolicyInfo
		rowNum := i+1
		policyName := policyInfo.PolicyName
		Namespace := policyInfo.Namespace
		Severity := policyInfo.Severity
		Threat := strings.Join([]string(policyInfo.Threat), "\n")

		RemediationSteps := strings.Join([]string(policyInfo.RemediationSteps), "\n")
		var entityType, Link, violationString string
		var violationsSummery []string
		for _, violation := range policyData.Violations {
			entityType = (&violation).ViolationEntityType
			Link = violation.CanonicalLink
			violationString = entityType + " " + Link
			violationsSummery = append(violationsSummery, violationString)
		}
		violationsPolicy := strings.Join([]string(violationsSummery), "\n")

		row := []string{strconv.Itoa(rowNum), policyName, Namespace, Severity, Threat, violationsPolicy, RemediationSteps}
		err := csvwriter.Write(row)
		if err != nil {
			panic(err)
		}

	}
	csvwriter.Write([]string{"\n"})
	if err != nil {
		panic(err)
	}


	return true
}


func (f *CsvFormatter) formatSummery(output *scheme.Flattened, csvwriter *csv.Writer) bool {
	headers := []string{"#", "Namespace", "Policy", "Severity", "Passed", "Failed", "Skipped"}
	err := csvwriter.Write(headers)

	for i, policyName := range output.AsOrderedMap().Keys() {
		rowNum := i+1
		data := output.GetPolicyData(policyName)
		policyInfo := data.PolicyInfo
		title := policyInfo.Title
		severity := policyInfo.Severity
		namespace := policyInfo.Namespace

		var passed, failed, skipped int
		for _, violation := range data.Violations {
			switch violation.Status {
			case analyzers.PolicyPassed:
				passed++
			case analyzers.PolicyFailed:
				failed++
			case analyzers.PolicySkipped:
				skipped++
			}
		}


		row := []string{strconv.Itoa(rowNum), namespace, title, severity, strconv.Itoa(passed), strconv.Itoa(failed), strconv.Itoa(skipped)}
		err := csvwriter.Write(row)
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	return true
}

func (f *CsvFormatter) Format(output scheme.Scheme, failedOnly bool) ([]byte, error) {
	var csvBuffer bytes.Buffer
	csvWriter := csv.NewWriter(&csvBuffer)

	typedOutput, ok := output.(*scheme.Flattened)
	if !ok {
		return nil, UnsupportedScheme{output}
	}

	f.csvFailedPolicies(typedOutput, csvWriter)
	if !failedOnly {
		f.formatSummery(typedOutput, csvWriter)
	}
	csvWriter.Flush()

	// Check for errors during flushing
	if err := csvWriter.Error(); err != nil {
		fmt.Printf("Error: %v\n", err)
		panic(err)
	}

	csvData := csvBuffer.Bytes()

	return csvData, nil
}

func (f *CsvFormatter) IsSchemeSupported(schemeType string) bool {
	return schemeType == scheme.TypeFlattened
}
