package k6exec

import (
	"context"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/grafana/k6deps"
)

type module struct {
	name    string
	version *semver.Version
}

type modules map[string]*module

func (mods *modules) unmarshalVersionOutput(text []byte) error {
	*mods = make(modules)

	if allmatch := reK6.FindAllSubmatch(text, -1); allmatch != nil {
		match := allmatch[0]

		version, err := semver.NewVersion(string(match[idxK6Version]))
		if err != nil {
			return err
		}

		(*mods)[k6module] = &module{name: k6module, version: version}
	}

	for _, match := range reExtension.FindAllSubmatch(text, -1) {
		version, err := semver.NewVersion(string(match[idxExtVersion]))
		if err != nil {
			return err
		}

		name := string(match[idxExtName])

		(*mods)[name] = &module{name: name, version: version}
	}

	return nil
}

func (mods modules) fulfill(deps k6deps.Dependencies) bool {
	for _, dep := range deps {
		mod, found := mods[dep.Name]
		if !found || (dep.Constraints != nil && !dep.Constraints.Check(mod.version)) {
			return false
		}
	}

	return true
}

func (mods modules) merge(deps k6deps.Dependencies) k6deps.Dependencies {
	merged := make(k6deps.Dependencies, len(deps))

	for name, dep := range deps {
		merged[name] = dep
	}

	for name := range mods {
		if _, found := merged[name]; !found {
			merged[name] = &k6deps.Dependency{Name: name}
		}
	}

	if _, found := merged[k6module]; !found {
		merged[k6module] = &k6deps.Dependency{Name: k6module}
	}

	return merged
}

func unmarshalVersionOutput(ctx context.Context, cmd string) (modules, error) {
	c := exec.CommandContext(ctx, cmd, "version")

	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var mods modules

	err = mods.unmarshalVersionOutput(out)
	if err != nil {
		return nil, err
	}

	return mods, nil
}

//nolint:gochecknoglobals
var (
	reK6          = regexp.MustCompile(`k6 (?P<k6Version>[^ ]+) .*`)
	reExtension   = regexp.MustCompile(`  (?P<extModule>[^ ]+) (?P<extVersion>[^,]+), (?P<extName>[^ ]+) \[([^\]]+)\]`)
	idxK6Version  = reK6.SubexpIndex("k6Version")
	idxExtVersion = reExtension.SubexpIndex("extVersion")
	idxExtName    = reExtension.SubexpIndex("extName")
)

const k6module = "k6"
