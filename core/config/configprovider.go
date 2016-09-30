// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package config

import (
	"errors"
	"fmt"
)

type ConfigurationProvider interface {
	Name() string // the name of the provider (YAML, Env, etc)
	GetValue(key string) ConfigurationValue
	Scope(prefix string) ConfigurationProvider
}

type ConfigurationChangeCallback func(key string, provider string, configdata interface{})

type DynamicConfigurationProvider interface {
	ConfigurationProvider

	RegisterChangeCallback(key string, callback ConfigurationChangeCallback) string
	UnregisterChangeCallback(token string) bool
	Shutdown()
}

func keyNotFound(key string) error {
	return errors.New(fmt.Sprintf("Couldn't find key %q", key))
}

type scopedProvider struct {
	prefix string

	child ConfigurationProvider
}

func newScopedProvider(prefix string, provider ConfigurationProvider) ConfigurationProvider {
	return &scopedProvider{prefix, provider}
}

func (sp scopedProvider) Name() string {
	return sp.child.Name()
}

func (sp scopedProvider) GetValue(key string) ConfigurationValue {
	if sp.prefix != "" {
		key = fmt.Sprintf("%s.%s", sp.prefix, key)
	}
	return sp.child.GetValue(key)
}

func (sp scopedProvider) Scope(prefix string) ConfigurationProvider {
	return newScopedProvider(prefix, sp)
}
