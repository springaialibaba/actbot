// Copyright 2024-2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package actors

import (
	"github.com/ShyunnY/actbot/internal/options/dingtalk"
)

// Constant definitions related to GitHub labels
const (
	// HelpWantedLabel The value of the help wanted label has been defined
	HelpWantedLabel = "help wanted"
)

// Constant definitions related to GitHub comment reaction
const (
	// CommendReaction The value of the "+1 👍" reaction has been defined
	CommendReaction = "+1"

	// RocketReaction The value of the "rocket 🚀" reaction has been defined
	RocketReaction = "rocket"
)

type Actor interface {
	Handler() error

	Capture(event GenericEvent) bool

	Name() string
}

type GenericEvent struct {
	// This represents the actual GitHub events
	Event any
}

// Options GitHub Actor extension options.
type Options struct {
	*dingtalk.DingTalkClient
}
