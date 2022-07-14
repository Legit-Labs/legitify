package opa

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/policies"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/loader"
)

func Load(policyPaths []string) (opa_engine.Enginer, error) {
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

	bundledModules, err := loadModules()
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

func loadModules() ([]*ast.Module, error) {
	return loadModulesFromFs(policies.Bundle, path.Dir(""))
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
