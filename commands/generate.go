package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/loader"
	"github.com/ms-henglu/armstrong/resource"
	"github.com/ms-henglu/armstrong/types"
)

type GenerateCommand struct {
	Ui                cli.Ui
	path              string
	useRawJsonPayload bool
	overwrite         bool
}

func (c *GenerateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("generate")

	fs.StringVar(&c.path, "path", "", "filepath of rest api to create arm resource example")
	fs.BoolVar(&c.useRawJsonPayload, "raw", false, "whether use raw json payload in `body`")
	fs.BoolVar(&c.overwrite, "overwrite", false, "whether overwrite existing terraform configurations")
	fs.Usage = func() { c.Ui.Error(c.Help()) }

	return fs
}

func (c GenerateCommand) Help() string {
	helpText := `
Usage: armstrong generate -path <filepath to example>
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c GenerateCommand) Synopsis() string {
	return "Generate testing files including terraform configuration for dependencies and testing resource."
}

func (c GenerateCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if len(c.path) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}
	if c.overwrite {
		_ = os.Remove("testing.tf")
		_ = os.Remove("dependency.tf")
	}
	err := ioutil.WriteFile("provider.tf", hclwrite.Format([]byte(hcl.ProviderHcl)), 0644)
	if err != nil {
		log.Fatalf("[Error] error writing provider.tf: %+v\n", err)
	}

	log.Println("[INFO] ----------- generate dependency and test resource ---------")
	// load dependencies
	log.Println("[INFO] loading dependencies")
	existDeps, deps := loadDependencies()

	// load example and generate hcl
	log.Println("[INFO] generating testing files")
	exampleFilepath := c.path
	exampleResource, err := resource.NewResourceFromExample(exampleFilepath)
	if err != nil {
		log.Fatalf("[Error] error reading example file: %+v\n", err)
	}

	dependencyHcl := exampleResource.DependencyHcl(existDeps, deps)
	err = appendFile("dependency.tf", dependencyHcl)
	if err != nil {
		log.Fatalf("[Error] error writing dependency.tf: %+v\n", err)
	}
	log.Println("[INFO] dependency.tf generated")

	testResourceHcl := exampleResource.Hcl(dependencyHcl, c.useRawJsonPayload)
	err = appendFile("testing.tf", testResourceHcl)
	if err != nil {
		log.Fatalf("[Error] error writing testing.tf: %+v\n", err)
	}
	log.Println("[INFO] testing.tf generated")
	return 0
}

func appendFile(filename string, hclContent string) error {
	content := hclContent
	if _, err := os.Stat(filename); err == nil {
		existingHcl, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("[WARN] reading %s: %+v", filename, err)
		} else {
			content = hcl.Combine(string(existingHcl), content)
		}
	}
	return ioutil.WriteFile(filename, hclwrite.Format([]byte(content)), 0644)
}

func loadDependencies() ([]types.Dependency, []types.Dependency) {
	mappingJsonLoader := loader.MappingJsonDependencyLoader{}
	hardcodeLoader := loader.HardcodeDependencyLoader{}

	deps := make([]types.Dependency, 0)
	depsMap := make(map[string]types.Dependency)
	if temp, err := mappingJsonLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	if temp, err := hardcodeLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	for _, dep := range depsMap {
		deps = append(deps, dep)
	}
	existDeps := hcl.LoadExistingDependencies()
	for i := range existDeps {
		ref := existDeps[i].ResourceType + "." + existDeps[i].ReferredProperty
		if dep, ok := depsMap[ref]; ok {
			existDeps[i].Pattern = dep.Pattern
		}
	}
	return existDeps, deps
}
