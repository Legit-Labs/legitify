package opa_engine

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
)

type Enginer interface {
	Query(ctx context.Context, namespace string, input interface{}) ([]QueryResult, error)
	SetTracing(enabled bool)
	Namespaces() []string
	Modules() map[string]*ast.Module
	Annotations() *ast.AnnotationSet
}

func NewEnginer(modules map[string]*ast.Module, compiler *ast.Compiler) Enginer {
	return &enginer{
		modules:  modules,
		compiler: compiler,
	}
}

type QueryResult struct {
	PolicyName               string
	FullyQualifiedPolicyName string
	Annotations              *ast.Annotations
	ExtraData                interface{}
	IsViolation              bool
}

type enginer struct {
	modules       map[string]*ast.Module
	compiler      *ast.Compiler
	enableTracing bool
}

func (e *enginer) SetTracing(enabled bool) {
	e.enableTracing = enabled
}

func (engine *enginer) Modules() map[string]*ast.Module {
	return engine.modules
}

// Namespaces returns all the namespace in the engine.
func (engine *enginer) Namespaces() []string {
	namespaces := []string{}
	for _, module := range engine.modules {
		namespace := strings.Replace(module.Package.Path.String(), "data.", "", 1)
		for _, ns := range namespaces {
			if ns == namespace {
				continue
			}
		}

		namespaces = append(namespaces, namespace)
	}

	return namespaces
}

func (engine *enginer) Annotations() *ast.AnnotationSet {
	return engine.compiler.GetAnnotationSet()
}

func (engine *enginer) Query(ctx context.Context, namespace string, input interface{}) ([]QueryResult, error) {
	result, err := engine.queryPolicy(ctx, namespace, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %s %s", namespace, err)
	}

	return result, nil
}

func (engine *enginer) queryPolicy(ctx context.Context, namespace string, input interface{}) ([]QueryResult, error) {
	regoInstance := engine.buildRegoInstance(namespace, input)

	resultSet, err := regoInstance.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("query eval: %w", err)
	} else {
		return engine.parseResultsSet(resultSet), nil
	}
}

func (engine *enginer) buildRegoInstance(namespace string, input interface{}) *rego.Rego {
	return rego.New(
		rego.Query(fmt.Sprintf("data.%s", namespace)),
		rego.Input(input),
		rego.Compiler(engine.compiler),
		rego.Trace(engine.enableTracing),
		rego.StrictBuiltinErrors(true),
		rego.PrintHook(topdown.NewPrintHook(os.Stderr)),
	)
}

func (engine *enginer) parseResultsSet(rs rego.ResultSet) []QueryResult {
	result := make([]QueryResult, 0)

	for _, r := range rs {
		for _, exp := range r.Expressions {
			baseModule := exp.Text
			matchedPolicies := parseResults(exp.Value, baseModule)

			for _, m := range matchedPolicies {
				match := m.fullPolicyName
				split := strings.Split(match, ".")
				current := QueryResult{
					FullyQualifiedPolicyName: match,
					PolicyName:               split[len(split)-1],
					Annotations:              engine.findAnnotation(match),
					ExtraData:                m.extraData,
					IsViolation:              m.violation,
				}
				result = append(result, current)
			}
		}
	}

	return result
}

func (engine *enginer) findAnnotation(policyFullPath string) *ast.Annotations {
	for _, anno := range engine.Annotations().Flatten() {
		if anno.Path.String() == policyFullPath {
			return anno.Annotations
		}
	}

	return nil
}

type matchedPolicy struct {
	fullPolicyName string
	extraData      interface{}
	violation      bool
}

func parseResults(i interface{}, path string) []matchedPolicy {
	var result []matchedPolicy
	mapped := i.(map[string]interface{})

	for k, v := range mapped {
		fullPath := fmt.Sprintf("%s.%s", path, k)

		violation := true
		// policies that return value but didn't have a match returns an empty map and should be ignored
		extra, ok := v.(map[string]interface{})
		if ok && len(extra) == 0 {
			violation = false
		}
		violated, ok := v.(bool)
		if ok && !violated {
			violation = false
		}

		result = append(result, matchedPolicy{
			fullPolicyName: fullPath,
			extraData:      v,
			violation:      violation,
		})
	}

	return result
}
