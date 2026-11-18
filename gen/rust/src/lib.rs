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

//! Kei connector operation contracts, generated from the protobuf sources in
//! `proto/` via buf's prost and prost-serde plugins.
//!
//! The types are plain data: fill them by JSON unmarshal; nothing here
//! executes connector calls.

/// `kei.contracts.v1` — the operation contract options shared by connectors.
pub mod kei {
    /// Contracts.
    pub mod contracts {
        /// Version 1.
        pub mod v1 {
            // The generated file includes the matching .serde.rs itself.
            include!("gen/kei/contracts/v1/kei.contracts.v1.rs");
        }
    }
    /// Connectors.
    pub mod connectors {
        /// GitHub connector contract.
        pub mod github {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/github/v1/kei.connectors.github.v1.rs");
            }
        }
        /// Linear connector contract.
        pub mod linear {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/linear/v1/kei.connectors.linear.v1.rs");
            }
        }
        /// Generic HTTP API connector contract.
        pub mod http_api {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/http_api/v1/kei.connectors.http_api.v1.rs");
            }
        }
        /// Database connector contract.
        pub mod database {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/database/v1/kei.connectors.database.v1.rs");
            }
        }
        /// Atlassian connector contract: Jira and Confluence reads.
        pub mod atlassian {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/atlassian/v1/kei.connectors.atlassian.v1.rs");
            }
        }
        /// Airtable connector contract.
        pub mod airtable {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/airtable/v1/kei.connectors.airtable.v1.rs");
            }
        }
        /// HubSpot connector contract.
        pub mod hubspot {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/hubspot/v1/kei.connectors.hubspot.v1.rs");
            }
        }
        /// Notion connector contract.
        pub mod notion {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/notion/v1/kei.connectors.notion.v1.rs");
            }
        }
        /// Sentry connector contract.
        pub mod sentry {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/sentry/v1/kei.connectors.sentry.v1.rs");
            }
        }
        /// MCP connector contract: tool calls against remote MCP servers.
        pub mod mcp {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/mcp/v1/kei.connectors.mcp.v1.rs");
            }
        }
    }
}
