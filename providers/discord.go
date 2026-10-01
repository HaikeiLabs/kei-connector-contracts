// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package providers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Discord is a contract-only provider: the payloads and resource grammar
// below are what the tenant runtime executes against the Discord API with a
// bot token resolved from the connector's opaque_ref. There is no backend
// seam here; the runtime implements the reads directly.
//
// Resources (snowflake = Discord's 17-20 digit ID):
//
//	guild.list    guilds
//	guild.read    guilds/<snowflake>
//	channel.list  guilds/<snowflake>/channels
//	channel.read  channels/<snowflake>
//	thread.list   channels/<snowflake>/threads
//	message.list  channels/<snowflake>/messages
//	message.read  channels/<snowflake>/messages/<snowflake>
//	role.list     guilds/<snowflake>/roles
//
// There is no member.* capability: member lists are personal data.

// DiscordMaxMessageListLimit is the largest page message.list may request,
// which is also Discord's own maximum.
const DiscordMaxMessageListLimit = 100

var snowflake = regexp.MustCompile(`^[0-9]{17,20}$`)

// discordResources maps each capability to its resource segments; "*" is a
// snowflake.
var discordResources = map[string][]string{
	"guild.list":   {"guilds"},
	"guild.read":   {"guilds", "*"},
	"channel.list": {"guilds", "*", "channels"},
	"channel.read": {"channels", "*"},
	"thread.list":  {"channels", "*", "threads"},
	"message.list": {"channels", "*", "messages"},
	"message.read": {"channels", "*", "messages", "*"},
	"role.list":    {"guilds", "*", "roles"},
}

// ValidateDiscordResource checks that resource is the resource shape of a
// Discord capability.
func ValidateDiscordResource(capability, resource string) error {
	return validateResource("discord", discordResources, capability, resource, snowflake.MatchString)
}

type DiscordGuildListPayload struct{}

func (DiscordGuildListPayload) Capability() string { return "guild.list" }
func (DiscordGuildListPayload) Validate() error    { return nil }

type DiscordGuildReadPayload struct{}

func (DiscordGuildReadPayload) Capability() string { return "guild.read" }
func (DiscordGuildReadPayload) Validate() error    { return nil }

type DiscordChannelListPayload struct{}

func (DiscordChannelListPayload) Capability() string { return "channel.list" }
func (DiscordChannelListPayload) Validate() error    { return nil }

type DiscordChannelReadPayload struct{}

func (DiscordChannelReadPayload) Capability() string { return "channel.read" }
func (DiscordChannelReadPayload) Validate() error    { return nil }

type DiscordThreadListPayload struct{}

func (DiscordThreadListPayload) Capability() string { return "thread.list" }
func (DiscordThreadListPayload) Validate() error    { return nil }

// DiscordMessageListPayload pages through a channel's messages. Limit 0 is
// the provider default; Before and After are snowflake cursors and at most
// one may be set.
type DiscordMessageListPayload struct {
	Limit  int    `json:"limit,omitempty"`
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

func (DiscordMessageListPayload) Capability() string { return "message.list" }

func (p DiscordMessageListPayload) Validate() error {
	if p.Limit < 0 || p.Limit > DiscordMaxMessageListLimit {
		return fmt.Errorf("message.list limit must be between 1 and %d", DiscordMaxMessageListLimit)
	}
	if p.Before != "" && p.After != "" {
		return errors.New("message.list accepts before or after, not both")
	}
	for name, cursor := range map[string]string{"before": p.Before, "after": p.After} {
		if cursor != "" && !snowflake.MatchString(cursor) {
			return fmt.Errorf("message.list %s must be a snowflake", name)
		}
	}
	return nil
}

type DiscordMessageReadPayload struct{}

func (DiscordMessageReadPayload) Capability() string { return "message.read" }
func (DiscordMessageReadPayload) Validate() error    { return nil }

type DiscordRoleListPayload struct{}

func (DiscordRoleListPayload) Capability() string { return "role.list" }
func (DiscordRoleListPayload) Validate() error    { return nil }

// validateResource matches resource against the capability's segment shape,
// where "*" is an identifier accepted by id.
func validateResource(provider string, grammar map[string][]string, capability, resource string, id func(string) bool) error {
	shape, ok := grammar[capability]
	if !ok {
		return fmt.Errorf("capability %q is not defined for provider %q", capability, provider)
	}
	parts := strings.Split(resource, "/")
	if len(parts) != len(shape) {
		return fmt.Errorf("resource must be %s", strings.Join(shape, "/"))
	}
	for i, want := range shape {
		if (want == "*" && !id(parts[i])) || (want != "*" && parts[i] != want) {
			return fmt.Errorf("resource must be %s", strings.Join(shape, "/"))
		}
	}
	return nil
}
