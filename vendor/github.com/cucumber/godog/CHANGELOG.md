# CHANGE LOG

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org).

This document is formatted according to the principles of [Keep A CHANGELOG](http://keepachangelog.com).

----

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

## [v0.10.0]

### Added
- Added concurrency support to the pretty formatter ([275](https://github.com/cucumber/godog/pull/275) - [lonnblad])
- Added concurrency support to the events formatter ([274](https://github.com/cucumber/godog/pull/274) - [lonnblad])
- Added concurrency support to the cucumber formatter ([273](https://github.com/cucumber/godog/pull/273) - [lonnblad])
- Added an example for how to use assertion pkgs like testify with godog ([289](https://github.com/cucumber/godog/pull/289) - [lonnblad])
- Added the new TestSuiteInitializer and ScenarioInitializer ([294](https://github.com/cucumber/godog/pull/294) - [lonnblad])
- Added an in-mem storage for pickles ([304](https://github.com/cucumber/godog/pull/304) - [lonnblad])
- Added Pickle and PickleStep results to the in-mem storage ([305](https://github.com/cucumber/godog/pull/305) - [lonnblad])
- Added features to the in-mem storage ([306](https://github.com/cucumber/godog/pull/306) - [lonnblad])
- Broke out some code from massive files into new files ([307](https://github.com/cucumber/godog/pull/307) - [lonnblad])
- Added support for concurrent scenarios ([311](https://github.com/cucumber/godog/pull/311) - [lonnblad])

### Changed
- Broke out snippets gen and added sorting on method name ([271](https://github.com/cucumber/godog/pull/271) - [lonnblad])
- Updated so that we run all tests concurrent now ([278](https://github.com/cucumber/godog/pull/278) - [lonnblad])
- Moved fmt tests to a godog_test pkg and restructured the fmt output tests ([295](https://github.com/cucumber/godog/pull/295) - [lonnblad])
- Moved builder tests to a godog_test pkg ([296](https://github.com/cucumber/godog/pull/296) - [lonnblad])
- Made the builder tests run in parallel ([298](https://github.com/cucumber/godog/pull/298) - [lonnblad])
- Refactored suite_context.go  ([300](https://github.com/cucumber/godog/pull/300) - [lonnblad])
- Added better testing of the Context Initializers and TestSuite{}.Run() ([301](https://github.com/cucumber/godog/pull/301) - [lonnblad])
- Updated the README.md  ([302](https://github.com/cucumber/godog/pull/302) - [lonnblad])
- Unexported some exported properties in unexported structs ([303](https://github.com/cucumber/godog/pull/303) - [lonnblad])
- Refactored some states in the formatters and feature struct ([310](https://github.com/cucumber/godog/pull/310) - [lonnblad])

### Deprecated
- Deprecated SuiteContext and ConcurrentFormatter ([314](https://github.com/cucumber/godog/pull/314) - [lonnblad])

### Removed
- Removed pre go112 build code ([293](https://github.com/cucumber/godog/pull/293) - [lonnblad])
- Removed the deprecated feature hooks ([312](https://github.com/cucumber/godog/pull/312) - [lonnblad])

### Fixed
- Fixed failing builder tests due to the v0.9.0 change ([lonnblad])
- Update paths to screenshots for examples ([270](https://github.com/cucumber/godog/pull/270) - [leviable])
- Made progress formatter verification a bit more accurate ([lonnblad])
- Added comparison between single and multi threaded runs ([272](https://github.com/cucumber/godog/pull/272) - [lonnblad])
- Fixed issue with empty feature file causing nil pointer deref ([288](https://github.com/cucumber/godog/pull/288) - [lonnblad])
- Updated linting checks in circleci config and fixed linting issues ([290](https://github.com/cucumber/godog/pull/290) - [lonnblad])
- Readded some legacy doc for FeatureContext ([297](https://github.com/cucumber/godog/pull/297) - [lonnblad])
- Fixed an issue with calculating time for junit testsuite ([308](https://github.com/cucumber/godog/pull/308) - [lonnblad])
- Fixed so that we don't execute features with zero scenarios ([315](https://github.com/cucumber/godog/pull/315) - [lonnblad])
- Fixed the broken --random flag ([317](https://github.com/cucumber/godog/pull/317) - [lonnblad])

## [0.9.0]

### Added

### Changed

- Run godog features in CircleCI in strict mode ([jaysonesmith])
- Removed TestMain call in `suite_test.go` for CI. ([jaysonesmith])
- Migrated to [gherkin-go - v11.0.0](https://github.com/cucumber/gherkin-go/releases/tag/v11.0.0). ([240](https://github.com/cucumber/godog/pull/240) - [lonnblad])

### Deprecated

### Removed

### Fixed

- Fixed the time attributes in the JUnit formatter. ([232](https://github.com/cucumber/godog/pull/232) - [lonnblad])
- Re enable custom formatters. ([238](https://github.com/cucumber/godog/pull/238) - [ericmcbride])
- Added back suite_test.go ([jaysonesmith])
- Normalise module paths for use on Windows ([242](https://github.com/cucumber/godog/pull/242) - [gjtaylor])
- Fixed panic in indenting function `s` ([247](https://github.com/cucumber/godog/pull/247) - [titouanfreville])
- Fixed wrong version in API example ([263](https://github.com/cucumber/godog/pull/263) - [denis-trofimov])

## [0.8.1]

### Added

- Link in Readme to the Slack community. ([210](https://github.com/cucumber/godog/pull/210) - [smikulcik])
- Added run tests for Cucumber formatting. ([214](https://github.com/cucumber/godog/pull/214), [216](https://github.com/cucumber/godog/pull/216) - [lonnblad])

### Changed

- Renamed the `examples` directory to `_examples`, removing dependencies from the Go module ([218](https://github.com/cucumber/godog/pull/218) - [axw])

### Deprecated

### Removed

### Fixed

- Find/Replaced references to DATA-DOG/godog -> cucumber/godog for docs. ([209](https://github.com/cucumber/godog/pull/209) - [smikulcik])
- Fixed missing links in changelog to be correctly included! ([jaysonesmith])

## [0.8.0]

### Added

- Added initial CircleCI config. ([jaysonesmith])
- Added concurrency support for JUnit formatting ([lonnblad])

### Changed

- Changed code references to DATA-DOG/godog to cucumber/godog to help get things building correctly. ([jaysonesmith])

### Deprecated

### Removed

### Fixed

<!-- Releases -->
[Unreleased]: https://github.com/cucumber/cucumber/compare/godog/v0.8.1...master
[0.8.0]:      https://github.com/cucumber/cucumber/compare/godog/v0.8.0...godog/v0.8.1
[0.8.0]:      https://github.com/cucumber/cucumber/compare/godog/v0.7.13...godog/v0.8.0

<!-- Contributors -->
[axw]:              https://github.com/axw
[jaysonesmith]:     https://github.com/jaysonesmith
[lonnblad]:         https://github.com/lonnblad
[smikulcik]:        https://github.com/smikulcik
[ericmcbride]:      https://github.com/ericmcbride
[gjtaylor]:         https://github.com/gjtaylor
[titouanfreville]:  https://github.com/titouanfreville
[denis-trofimov]:   https://github.com/denis-trofimov
[leviable]:         https://github.com/leviable
