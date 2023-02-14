package bucket

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hellofresh/stats-go/log"
)

const (
	sectionsDelimiter   = ":"
	logSuspiciousMetric = "Second level ID auto-discover found suspicious metric"

	// SectionTestTrue is a name for "stats.TestAlwaysTrue" test callback function
	SectionTestTrue = "true"
	// SectionTestIsNumeric is a name for "stats.TestIsNumeric" test callback function
	SectionTestIsNumeric = "numeric"
	// SectionTestIsNotEmpty is a name for "stats.TestIsNotEmpty" test callback function
	SectionTestIsNotEmpty = "not_empty"
)

var (
	// ErrInvalidFormat error indicates that sections string has invalid format
	ErrInvalidFormat = errors.New("invalid sections format")

	// ErrUnknownSectionTest error indicates that section has unknown test callback name
	ErrUnknownSectionTest = errors.New("unknown section test")

	// ErrLooksLikeID error indicates that second level ID auto-discover found suspicious metric
	ErrLooksLikeID = errors.New("metric looks like ID")
)

// PathSection type represents single path section string
type PathSection string

// SectionTestCallback type represents section test callback function
type SectionTestCallback func(PathSection) bool

// SectionTestDefinition type represents section test callback definition
type SectionTestDefinition struct {
	Name     string
	Callback SectionTestCallback
}

// SectionsTestsMap type represents section test callbacks definitions map
type SectionsTestsMap map[PathSection]SectionTestDefinition

// String returns pretty formatted string representation of SectionsTestsMap
func (m SectionsTestsMap) String() string {
	var sections []string
	for k, v := range m {
		sections = append(sections, fmt.Sprintf("%s: %s", k, v.Name))
	}

	sort.Strings(sections)
	return fmt.Sprintf("[%s]", strings.Join(sections, ", "))
}

// TestAlwaysTrue section test callback function that gives true result to any section
func TestAlwaysTrue(PathSection) bool {
	return true
}

// TestIsNumeric section test callback function that gives true result if section is numeric
func TestIsNumeric(s PathSection) bool {
	_, err := strconv.Atoi(string(s))
	return err == nil
}

// TestIsNotEmpty section test callback function that gives true result if section is not empty placeholder ("-")
func TestIsNotEmpty(s PathSection) bool {
	return string(s) != MetricEmptyPlaceholder
}

var (
	sectionsTestSync     sync.Mutex
	sectionsTestRegistry = map[string]SectionTestCallback{
		SectionTestTrue:       TestAlwaysTrue,
		SectionTestIsNumeric:  TestIsNumeric,
		SectionTestIsNotEmpty: TestIsNotEmpty,
	}
)

// SecondLevelIDConfig configuration struct for second level ID callback
type SecondLevelIDConfig struct {
	HasIDAtSecondLevel    SectionsTestsMap
	AutoDiscoverThreshold uint
	AutoDiscoverWhiteList []string

	autoDiscoverStorage  *metricStorage
	autoDiscoverWhiteMap map[string]bool
}

// NewHasIDAtSecondLevelCallback returns HttpMetricNameAlterCallback implementation that checks for IDs
// on the second level of HTTP Request path
func NewHasIDAtSecondLevelCallback(config *SecondLevelIDConfig) HTTPMetricNameAlterCallback {
	if config.AutoDiscoverThreshold > 0 {
		config.autoDiscoverStorage = newMetricStorage(config.AutoDiscoverThreshold)

		// convert array to map for easier search
		config.autoDiscoverWhiteMap = make(map[string]bool, len(config.AutoDiscoverWhiteList))
		for _, val := range config.AutoDiscoverWhiteList {
			config.autoDiscoverWhiteMap[val] = true
		}
	}

	return func(operation MetricOperation, r *http.Request) MetricOperation {
		firstFragment := "/"
		for _, fragment := range strings.Split(r.URL.Path, "/") {
			if fragment != "" {
				firstFragment = fragment
				break
			}
		}

		if testFunction, ok := config.HasIDAtSecondLevel[PathSection(firstFragment)]; ok {
			if testFunction.Callback(PathSection(operation[2])) {
				operation[2] = MetricIDPlaceholder
			}
		} else if config.AutoDiscoverThreshold > 0 {
			if _, ok := config.autoDiscoverWhiteMap[firstFragment]; !ok {
				if config.autoDiscoverStorage.LooksLikeID(firstFragment, operation[2]) {
					log.Log("Second level ID auto-discover found suspicious metric", map[string]interface{}{
						"method":    r.Method,
						"path":      r.URL.Path,
						"operation": operation,
					}, ErrLooksLikeID)

					operation[2] = MetricIDPlaceholder
				}
			}
		}

		return operation
	}
}

// RegisterSectionTest registers new section test callback function with its name
func RegisterSectionTest(name string, callback SectionTestCallback) {
	sectionsTestSync.Lock()
	defer sectionsTestSync.Unlock()

	sectionsTestRegistry[name] = callback
}

// GetSectionTestCallback returns section test callback function by name
func GetSectionTestCallback(name string) SectionTestCallback {
	sectionsTestSync.Lock()
	defer sectionsTestSync.Unlock()

	return sectionsTestRegistry[name]
}

// ParseSectionsTestsMap parses string into SectionsTestsMap.
// In most cases string comes as config to the application e.g. from env.
// Valid string formats are:
// 1. <section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>:<section-2>:<test-callback-name-2>
// 2. <section-0>:<test-callback-name-0>\n<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>
// 3. <section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>
func ParseSectionsTestsMap(s string) (SectionsTestsMap, error) {
	result := make(SectionsTestsMap)
	var parts []string

	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			for _, part := range strings.Split(strings.TrimSpace(line), sectionsDelimiter) {
				if strings.TrimSpace(part) != "" {
					parts = append(parts, part)
				}
			}
		}
	}
	if len(parts)%2 != 0 {
		return nil, ErrInvalidFormat
	}

	for i := 0; i < len(parts); i += 2 {
		pathSection := PathSection(parts[i])
		sectionTestName := parts[i+1]

		sectionTestCallback := GetSectionTestCallback(sectionTestName)
		if sectionTestCallback == nil {
			return nil, ErrUnknownSectionTest
		}
		result[pathSection] = SectionTestDefinition{sectionTestName, sectionTestCallback}
	}

	return result, nil
}
