package godog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/gherkin-go/v11"
	"github.com/cucumber/messages-go/v10"
)

var pathLineRe = regexp.MustCompile(`:([\d]+)$`)

func extractFeaturePathLine(p string) (string, int) {
	line := -1
	retPath := p
	if m := pathLineRe.FindStringSubmatch(p); len(m) > 0 {
		if i, err := strconv.Atoi(m[1]); err == nil {
			line = i
			retPath = p[:strings.LastIndexByte(p, ':')]
		}
	}
	return retPath, line
}

func parseFeatureFile(path string, newIDFunc func() string) (*feature, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	var buf bytes.Buffer
	gherkinDocument, err := gherkin.ParseGherkinDocument(io.TeeReader(reader, &buf), newIDFunc)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", path, err)
	}

	gherkinDocument.Uri = path
	pickles := gherkin.Pickles(*gherkinDocument, path, newIDFunc)

	f := feature{GherkinDocument: gherkinDocument, pickles: pickles, content: buf.Bytes()}
	return &f, nil
}

func parseFeatureDir(dir string, newIDFunc func() string) ([]*feature, error) {
	var features []*feature
	return features, filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(p, ".feature") {
			return nil
		}

		feat, err := parseFeatureFile(p, newIDFunc)
		if err != nil {
			return err
		}

		features = append(features, feat)
		return nil
	})
}

func parsePath(path string) ([]*feature, error) {
	var features []*feature

	path, line := extractFeaturePathLine(path)

	fi, err := os.Stat(path)
	if err != nil {
		return features, err
	}

	newIDFunc := (&messages.Incrementing{}).NewId

	if fi.IsDir() {
		return parseFeatureDir(path, newIDFunc)
	}

	ft, err := parseFeatureFile(path, newIDFunc)
	if err != nil {
		return features, err
	}

	// filter scenario by line number
	var pickles []*messages.Pickle
	for _, pickle := range ft.pickles {
		sc := ft.findScenario(pickle.AstNodeIds[0])

		if line == -1 || uint32(line) == sc.Location.Line {
			pickles = append(pickles, pickle)
		}
	}
	ft.pickles = pickles

	return append(features, ft), nil
}

func parseFeatures(filter string, paths []string) ([]*feature, error) {
	var order int

	featureIdxs := make(map[string]int)
	uniqueFeatureURI := make(map[string]*feature)
	for _, path := range paths {
		feats, err := parsePath(path)

		switch {
		case os.IsNotExist(err):
			return nil, fmt.Errorf(`feature path "%s" is not available`, path)
		case os.IsPermission(err):
			return nil, fmt.Errorf(`feature path "%s" is not accessible`, path)
		case err != nil:
			return nil, err
		}

		for _, ft := range feats {
			if _, duplicate := uniqueFeatureURI[ft.Uri]; duplicate {
				continue
			}

			uniqueFeatureURI[ft.Uri] = ft
			featureIdxs[ft.Uri] = order

			order++
		}
	}

	var features = make([]*feature, len(uniqueFeatureURI))
	for uri, feature := range uniqueFeatureURI {
		idx := featureIdxs[uri]
		features[idx] = feature
	}

	features = filterFeatures(filter, features)

	return features, nil
}

func filterFeatures(tags string, features []*feature) (result []*feature) {
	for _, ft := range features {
		ft.pickles = applyTagFilter(tags, ft.pickles)

		if ft.Feature != nil && len(ft.pickles) > 0 {
			result = append(result, ft)
		}
	}

	return
}

func applyTagFilter(tags string, pickles []*messages.Pickle) (result []*messages.Pickle) {
	if len(tags) == 0 {
		return pickles
	}

	for _, pickle := range pickles {
		if matchesTags(tags, pickle.Tags) {
			result = append(result, pickle)
		}
	}

	return
}
