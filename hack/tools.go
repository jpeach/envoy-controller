// +build tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
