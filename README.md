# Docker Semantic Versioning

## Purpose

Allow for easier generation of semantic versions based on [semver.org](https://semver.org).

## Building the image

```bash
> docker build -t docker-versioner .
```

### Available Variables

| ENV Var                | Example | Description                                                                       |
|------------------------|---------|-----------------------------------------------------------------------------------|
| version_base           | `3.0.0` | The base version to use, ignored unless `ignorePrevious` is true.                 |
| version_ignorePrevious | `true`  | If true the previous tags on the repo will be ignored when determining a version. |
| version_pattern        | `*.^.*` | How to increment a version, `*` will leave alone, `^` will increment by 1         |
| version_prerelease     | 432345  | Append pre-release data, this defaults to being a timestamp. After a `-` character| 
| version_metadata       | alpha   | Append release metadata, no default. After a `+` character.                       |

## Usage

### Basic

```shell
> docker run -v "$(pwd):/tmp" -w /tmp --rm -e version_base=1.3.5 -e version_ignorePrevious=true quay.cnqr.delivery/da-workflow/docker-versioner
# output 1.4.0
```

### Customizing the increment

```shell
> docker run -v "$(pwd):/tmp" -w /tmp --rm -e version_base=2.4.0 -e version_ignorePrevious=true -e version_pattern="^.*.*" quay.cnqr.delivery/da-workflow/docker-versioner
# output 3.0.0
```

```shell
> docker run -v "$(pwd):/tmp" -w /tmp --rm -e version_base=2.1.0 -e version_pattern="*.*.^" quay.cnqr.delivery/da-workflow/docker-versioner
# output 2.1.1
```
