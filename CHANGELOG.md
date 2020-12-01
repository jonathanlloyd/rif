# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Fixed
 - Validation error now correctly specifies missing field as `url` rather than
   `URL`

## [0.4.5] - 2019-04-22
### Fixed
 - Bug where missing URL / method fields caused a crash
 - Issue where only the first validation error is printed

## [0.4.4] - 2019-02-17
### Fixed
 - Bug where template variables were not interpolated into alternative output
   formats (HTTP / cURL)

## [0.4.3] - 2018-11-19
### Changed
 - Various refactorings

## [0.4.2] - 2018-04-11
### Changed
 - Improved error message when rif\_version field is missing

## [0.4.1] - 2018-03-21
### Fixed
 - Fixed bug where non-string default values were rendered incorrectly in error string

## [0.4.0] - 2018-03-18
### Added
 - cURL output format (rif <filename> --output=curl)
### Changed
 - Improved error message when a variable is given an unknown type in a RIF file
 - Improved error message when required variables are not given

## [0.3.2] - 2018-02-17
### Changed
 - Add a unique(ish) prefix to the build number to distinguish builds made
   from the same commit.

## [0.3.1] - 2018-02-04
### Fixed
 - Bug which prevented the use of the HTTP output format on requests with
   bodies.

## [0.3.0] - 2018-02-01
### Added
 - HTTP output format (rif <filename> --output=http)

## [0.2.0] - 2018-01-28
### Added
 - url, header & body templating
 - rif to HTTP request transformation
 - .rif file parsing

## [0.1.0] - 2017-12-29
### Added
 - Basic command implementation


[Unreleased]: https://github.com/jonathanlloyd/rif/compare/0.4.5...HEAD
[0.4.5]: https://github.com/jonathanlloyd/rif/compare/0.4.4...0.4.5
[0.4.4]: https://github.com/jonathanlloyd/rif/compare/0.4.3...0.4.4
[0.4.3]: https://github.com/jonathanlloyd/rif/compare/0.4.2...0.4.3
[0.4.2]: https://github.com/jonathanlloyd/rif/compare/0.4.1...0.4.2
[0.4.1]: https://github.com/jonathanlloyd/rif/compare/0.4.0...0.4.1
[0.4.0]: https://github.com/jonathanlloyd/rif/compare/0.3.2...0.4.0
[0.3.2]: https://github.com/jonathanlloyd/rif/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/jonathanlloyd/rif/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/jonathanlloyd/rif/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/jonathanlloyd/rif/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/jonathanlloyd/rif/releases/tag/0.1.0
