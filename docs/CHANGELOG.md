# Docker Versioner Releases

## 0.7.0

### Changed

* Updated error output to show the tag we were able to find if it cannot be parsed as SemVer.
* Corrected the examples in the readme.

## 0.6.0

### Added

* Allow users to provide a prefix, so the output version will be: `<prefix>major.minor.patch-b<timestamp>`, this can have issues with with some tooling so if you get errors about invalid SemVer please remove the prefix for compatibility.

### Fixed

* Change format slightly to play nicely with Helm, new format: `major.minor.patch-b<timestamp>`

## 0.5.0

### Added

* Implement test framework.
* Add expected variables to the readme.

### Fixed

* If no tag is found return a default instead of failing, this brings it back to the way the old versioner worked.
* Include `git` in the Docker image.