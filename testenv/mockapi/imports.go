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

package mockapi

// Blank imports of every connector's generated package: each one registers
// its file descriptors in protoregistry.GlobalFiles, which is where the
// route table is discovered from. A new connector lands here with its
// manifest entry, and the mock picks it up automatically.
import (
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/airtable/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/atlassian/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/gmail/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/google_sheets/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/hubspot/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/notion/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/omni/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/power_bi/v1"
	_ "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/sentry/v1"
)
