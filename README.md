# GoZix Echo

[documentation-img]: https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff
[documentation-url]: https://pkg.go.dev/github.com/gozix/echo/v3
[license-img]: https://img.shields.io/github/license/gozix/echo.svg?style=for-the-badge
[license-url]: https://github.com/gozix/echo/blob/master/LICENSE
[release-img]: https://img.shields.io/github/tag/gozix/echo.svg?label=release&color=24B898&logo=github&style=for-the-badge
[release-url]: https://github.com/gozix/echo/releases/latest
[build-status-img]: https://img.shields.io/github/actions/workflow/status/gozix/echo/go.yml?logo=github&style=for-the-badge
[build-status-url]: https://github.com/gozix/echo/actions
[go-report-img]: https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge
[go-report-url]: https://goreportcard.com/report/github.com/gozix/echo
[code-coverage-img]: https://img.shields.io/codecov/c/github/gozix/echo.svg?style=for-the-badge&logo=codecov
[code-coverage-url]: https://codecov.io/gh/gozix/echo

[![License][license-img]][license-url]
[![Documentation][documentation-img]][documentation-url]

[![Release][release-img]][release-url]
[![Build Status][build-status-img]][build-status-url]
[![Go Report Card][go-report-img]][go-report-url]
[![Code Coverage][code-coverage-img]][code-coverage-url]

The bundle provide a Echo integration to GoZix application.

## Installation

```shell
go get github.com/gozix/echo/v3
```

## Dependencies

* [validator](https://github.com/gozix/validator)
* [viper](https://github.com/gozix/viper)
* [zap](https://github.com/gozix/zap)

## Configuration

```json
{
  "echo": {
    "debug": false,
    "level": "debug",
    "static": {
      "prefix": "/",
      "root": ""
    },
    "hide_banner": false,
    "hide_port": false 
  }
}
```

## Built-in Tags

| Symbol                | Value                         | Description               | 
| --------------------- | ----------------------------- | ------------------------- |
| TagController         | echo.controller               | Add a controller          |
| TagConfigurator       | echo.configurator             | Add a configurator        |
| TagMiddleware         | echo.middleware               | Add a middleware          |

## Documentation

You can find documentation on [pkg.go.dev][documentation-url] and read source code if needed.

## Questions

If you have any questions, feel free to create an issue.
