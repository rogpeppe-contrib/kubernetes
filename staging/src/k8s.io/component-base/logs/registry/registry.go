/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

import (
	"fmt"
	"sort"

	"github.com/go-logr/logr"
)

// LogRegistry is new init LogFormatRegistry struct
var LogRegistry = NewLogFormatRegistry()

// LogFormatRegistry store klog format registry
type LogFormatRegistry struct {
	registry map[string]LogFormatFactory
	frozen   bool
}

// LogFormatFactory provides support for a certain additional,
// non-default log format.
type LogFormatFactory interface {
	// Create returns a logger.
	Create() logr.Logger
}

// NewLogFormatRegistry return new init LogFormatRegistry struct
func NewLogFormatRegistry() *LogFormatRegistry {
	return &LogFormatRegistry{
		registry: make(map[string]LogFormatFactory),
		frozen:   false,
	}
}

// Register new log format registry to global logRegistry.
// nil is valid and selects the default klog output.
func (lfr *LogFormatRegistry) Register(name string, factory LogFormatFactory) error {
	if lfr.frozen {
		return fmt.Errorf("log format is frozen, unable to register log format")
	}
	if _, ok := lfr.registry[name]; ok {
		return fmt.Errorf("log format: %s already exists", name)
	}
	lfr.registry[name] = factory
	return nil
}

// Get specified log format logger
func (lfr *LogFormatRegistry) Get(name string) (LogFormatFactory, error) {
	re, ok := lfr.registry[name]
	if !ok {
		return nil, fmt.Errorf("log format: %s does not exists", name)
	}
	return re, nil
}

// Set specified log format logger
func (lfr *LogFormatRegistry) Set(name string, factory LogFormatFactory) error {
	if lfr.frozen {
		return fmt.Errorf("log format is frozen, unable to set log format")
	}

	lfr.registry[name] = factory
	return nil
}

// Delete specified log format logger
func (lfr *LogFormatRegistry) Delete(name string) error {
	if lfr.frozen {
		return fmt.Errorf("log format is frozen, unable to delete log format")
	}

	delete(lfr.registry, name)
	return nil
}

// List names of registered log formats (sorted)
func (lfr *LogFormatRegistry) List() []string {
	formats := make([]string, 0, len(lfr.registry))
	for f := range lfr.registry {
		formats = append(formats, f)
	}
	sort.Strings(formats)
	return formats
}

// Freeze freezes the log format registry
func (lfr *LogFormatRegistry) Freeze() {
	lfr.frozen = true
}
