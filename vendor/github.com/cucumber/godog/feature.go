package godog

import (
	"github.com/cucumber/messages-go/v10"
)

type feature struct {
	*messages.GherkinDocument
	pickles []*messages.Pickle
	content []byte
}

type sortFeaturesByName []*feature

func (s sortFeaturesByName) Len() int           { return len(s) }
func (s sortFeaturesByName) Less(i, j int) bool { return s[i].Feature.Name < s[j].Feature.Name }
func (s sortFeaturesByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (f feature) findScenario(astScenarioID string) *messages.GherkinDocument_Feature_Scenario {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil && sc.Id == astScenarioID {
			return sc
		}
	}

	return nil
}

func (f feature) findBackground(astScenarioID string) *messages.GherkinDocument_Feature_Background {
	var bg *messages.GherkinDocument_Feature_Background

	for _, child := range f.GherkinDocument.Feature.Children {
		if tmp := child.GetBackground(); tmp != nil {
			bg = tmp
		}

		if sc := child.GetScenario(); sc != nil && sc.Id == astScenarioID {
			return bg
		}
	}

	return nil
}

func (f feature) findExample(exampleAstID string) (*messages.GherkinDocument_Feature_Scenario_Examples, *messages.GherkinDocument_Feature_TableRow) {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil {
			for _, example := range sc.Examples {
				for _, row := range example.TableBody {
					if row.Id == exampleAstID {
						return example, row
					}
				}
			}
		}
	}

	return nil, nil
}

func (f feature) findStep(astStepID string) *messages.GherkinDocument_Feature_Step {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil {
			for _, step := range sc.GetSteps() {
				if step.Id == astStepID {
					return step
				}
			}
		}

		if bg := child.GetBackground(); bg != nil {
			for _, step := range bg.GetSteps() {
				if step.Id == astStepID {
					return step
				}
			}
		}
	}

	return nil
}
