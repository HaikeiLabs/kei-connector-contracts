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
	"strings"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/contract"
)

const (
	guildID   = "81384788765712384"
	channelID = "175928847299117063"
	messageID = "1166474102939287632"
)

func discordPayloads() []Payload {
	return []Payload{
		DiscordGuildListPayload{}, DiscordGuildReadPayload{}, DiscordChannelListPayload{},
		DiscordChannelReadPayload{}, DiscordThreadListPayload{}, DiscordMessageListPayload{},
		DiscordMessageReadPayload{}, DiscordRoleListPayload{},
	}
}

// Every Discord capability in the contract has exactly one payload, and no
// payload names a capability the contract lacks.
func TestDiscordPayloadsCoverTheCapabilities(t *testing.T) {
	assertPayloadsCover(t, contract.ProviderDiscord, discordPayloads())
}

func TestDiscordResourceGrammar(t *testing.T) {
	valid := map[string]string{
		"guild.list":   "guilds",
		"guild.read":   "guilds/" + guildID,
		"channel.list": "guilds/" + guildID + "/channels",
		"channel.read": "channels/" + channelID,
		"thread.list":  "channels/" + channelID + "/threads",
		"message.list": "channels/" + channelID + "/messages",
		"message.read": "channels/" + channelID + "/messages/" + messageID,
		"role.list":    "guilds/" + guildID + "/roles",
	}
	for capability, resource := range valid {
		if err := ValidateDiscordResource(capability, resource); err != nil {
			t.Errorf("%s %q rejected: %v", capability, resource, err)
		}
	}
	for _, tc := range []struct{ capability, resource string }{
		{"guild.list", "guilds/" + guildID},
		{"guild.read", "guilds"},
		{"guild.read", "guilds/acme"},
		{"guild.read", "guilds/123"},
		{"guild.read", "guilds/" + guildID + "/"},
		{"channel.list", "channels/" + channelID + "/channels"},
		{"channel.read", "guilds/" + guildID + "/channels/" + channelID},
		{"thread.list", "channels/" + channelID + "/threads/" + messageID},
		{"message.list", "channels/../messages"},
		{"message.read", "channels/" + channelID + "/messages"},
		{"role.list", "guilds/" + guildID + "/members"},
		{"member.list", "guilds/" + guildID + "/members"},
		{"message.read", "/channels/" + channelID + "/messages/" + messageID},
		{"message.read", "channels/" + channelID + "/messages/" + messageID + "?x=1"},
	} {
		if err := ValidateDiscordResource(tc.capability, tc.resource); err == nil {
			t.Errorf("%s %q accepted", tc.capability, tc.resource)
		}
	}
}

func TestDiscordMessageListInputIsBounded(t *testing.T) {
	for _, p := range []DiscordMessageListPayload{
		{},
		{Limit: 1},
		{Limit: 100},
		{Limit: 50, Before: messageID},
		{After: messageID},
	} {
		if err := p.Validate(); err != nil {
			t.Errorf("%+v rejected: %v", p, err)
		}
	}
	for _, p := range []DiscordMessageListPayload{
		{Limit: -1},
		{Limit: 101},
		{Before: "yesterday"},
		{After: "12"},
		{Before: messageID, After: messageID},
	} {
		if err := p.Validate(); err == nil {
			t.Errorf("%+v accepted", p)
		}
	}
	for _, p := range discordPayloads() {
		v, ok := p.(interface{ Validate() error })
		if !ok {
			t.Fatalf("%T has no Validate", p)
		}
		if err := v.Validate(); err != nil {
			t.Errorf("zero %T rejected: %v", p, err)
		}
	}
}

func assertPayloadsCover(t *testing.T, provider contract.Provider, payloads []Payload) {
	t.Helper()
	seen := map[string]bool{}
	for _, p := range payloads {
		if seen[p.Capability()] {
			t.Errorf("%s: two payloads for %q", provider, p.Capability())
		}
		seen[p.Capability()] = true
	}
	var missing []string
	for _, c := range contract.CapabilitiesFor(provider) {
		if !seen[c.Name] {
			missing = append(missing, c.Name)
		}
		delete(seen, c.Name)
	}
	if len(missing) > 0 {
		t.Errorf("%s capabilities without a payload: %s", provider, strings.Join(missing, ", "))
	}
	for name := range seen {
		t.Errorf("%s payload %q is not a contract capability", provider, name)
	}
}
