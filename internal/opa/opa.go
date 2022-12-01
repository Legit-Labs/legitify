package opa

import (
	"embed"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"os"
	"path"
	"strings"

	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/policies"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/loader"
)

func Load(policyPaths []string, scm scm_type.ScmType) (opa_engine.Enginer, error) {
	loadedPolicies, err := loader.NewFileLoader().
		WithProcessAnnotation(true).
		Filtered(policyPaths, isRegoFile)
	if err != nil {
		return nil, opa_engine.NewErrPolicyLoad(err)
	}

	if len(policyPaths) != 0 && len(loadedPolicies.Modules) == 0 {
		return nil, opa_engine.NewErrNoPolicies(policyPaths)
	}

	modules := loadedPolicies.ParsedModules()
	compiler := ast.NewCompiler().WithEnablePrintStatements(true)

	bundledModules, err := loadModules(scm)
	if err != nil {
		return nil, err
	}

	for _, m := range bundledModules {
		modules[m.Package.Location.File] = m
	}

	compiler.Compile(modules)

	if compiler.Failed() {
		return nil, fmt.Errorf("compiler: %w", compiler.Errors)
	}

	engine := opa_engine.NewEnginer(modules, compiler)

	return engine, nil
}

func loadModules(scmType scm_type.ScmType) ([]*ast.Module, error) {
	switch scmType {
	case scm_type.GitHub:
		return loadModulesFromFs(policies.GitHubBundle, path.Dir(""))
	case scm_type.GitLab:
		return loadModulesFromFs(policies.GitLabBundle, path.Dir(""))
	default:
		return nil, fmt.Errorf("unknown scm type %s", scmType)
	}
}

func loadModulesFromFs(fs embed.FS, p string) ([]*ast.Module, error) {
	bundledModules := make([]*ast.Module, 0)
	files, err := fs.ReadDir(p)

	if err != nil {
		return nil, err
	}

	for _, de := range files {
		if de.IsDir() {
			c, err := loadModulesFromFs(fs, path.Join(p, de.Name()))

			if err != nil {
				return nil, err
			}

			bundledModules = append(bundledModules, c...)
		} else {
			data, err := fs.ReadFile(path.Join(p, de.Name()))

			if err != nil {
				return nil, err
			}

			module, err := ast.ParseModuleWithOpts(de.Name(), string(data), ast.ParserOptions{
				ProcessAnnotation: true,
			})

			if err != nil {
				return nil, err
			}

			bundledModules = append(bundledModules, module)
		}
	}

	return bundledModules, nil
}

func isRegoFile(_ string, info os.FileInfo, depth int) bool {
	return !info.IsDir() && !strings.HasSuffix(info.Name(), bundle.RegoExt)
}
