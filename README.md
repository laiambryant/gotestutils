# goTestUtils

[![Go Reference](https://pkg.go.dev/badge/github.com/laiambryant/gotestutils.svg)](https://pkg.go.dev/github.com/laiambryant/gotestutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/laiambryant/gotestutils)](https://goreportcard.com/report/github.com/laiambryant/gotestutils)
[![GitHub license](https://img.shields.io/github/license/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/issues)
[![GitHub stars](https://img.shields.io/github/stars/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/stargazers)
[![Coverage Status](https://coveralls.io/repos/github/laiambryant/goTestUtils/badge.svg?branch=main)](https://coveralls.io/github/laiambryant/goTestUtils?branch=main)

My Favourite testing utility methods for go. Includes utilities for:

- Characterization testing
- More to come...

## Installation

```bash
go get github.com/laiambryant/gotestutils
```

## Quick Start

### Testing Utilities

```go
import "github.com/laiambryant/gotestutils/ctesting"

// Example usage
testSuite := []ctesting.CharacterizationTest[int]{
    ctesting.NewCharacterizationTest(3, nil, func() (int, error) { 
        return sum(1, 2), nil 
    }),
}
results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, testSuite)
```
