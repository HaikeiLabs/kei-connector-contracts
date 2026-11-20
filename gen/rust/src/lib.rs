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
//! The types are plain data: fill them by JSON unmarshal. Generated
//! connector operations run on the hand-written `connection` runtime, which
//! takes its credentials as arguments and never fetches, stores or logs
//! them; SQL connectors currently run through the caller's own transport.

/// The hand-written runtime that generated connector operations run on:
/// the authenticated connection type, the path renderer, canonical
/// request/response JSON helpers and the error mapping. It takes
/// credentials as arguments and never fetches, stores or logs them, and it
/// reaches the network only through the `reqwest::Client` the caller passes
/// in. SQL connectors have no driver dependency: they run through the
/// caller's own transport until a Rust consumer needs one.
pub mod connection;

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
                /// Generated GitHub connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/github/v1/ops.rs");
                }
            }
        }
        /// Linear connector contract.
        pub mod linear {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/linear/v1/kei.connectors.linear.v1.rs");
                /// Generated Linear connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/linear/v1/ops.rs");
                }
            }
        }
        /// Generic HTTP API connector contract.
        pub mod http_api {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/http_api/v1/kei.connectors.http_api.v1.rs");
                /// Generated HTTP API connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/http_api/v1/ops.rs");
                }
            }
        }
        /// Omni connector contract.
        pub mod omni {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/omni/v1/kei.connectors.omni.v1.rs");
                /// Generated Omni connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/omni/v1/ops.rs");
                }
            }
        }
        /// Power BI connector contract.
        pub mod power_bi {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/power_bi/v1/kei.connectors.power_bi.v1.rs");
                /// Generated Power BI connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/power_bi/v1/ops.rs");
                }
            }
        }
        /// Database connector contract.
        pub mod database {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/database/v1/kei.connectors.database.v1.rs");
                /// Generated database connector operations, driven by the
                /// caller's own SQL transport.
                pub mod ops {
                    include!("kei/connectors/database/v1/ops.rs");
                }
            }
        }
        /// Gmail connector contract.
        pub mod gmail {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/gmail/v1/kei.connectors.gmail.v1.rs");
                /// Generated Gmail connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/gmail/v1/ops.rs");
                }
            }
        }
        /// Google Calendar connector contract.
        pub mod google_calendar {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/google_calendar/v1/kei.connectors.google_calendar.v1.rs");
                /// Generated Google Calendar connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/google_calendar/v1/ops.rs");
                }
            }
        }
        /// Google Sheets connector contract.
        pub mod google_sheets {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/google_sheets/v1/kei.connectors.google_sheets.v1.rs");
                /// Generated Google Sheets connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/google_sheets/v1/ops.rs");
                }
            }
        }
        /// Atlassian connector contract: Jira and Confluence reads.
        pub mod atlassian {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/atlassian/v1/kei.connectors.atlassian.v1.rs");
                /// Generated Atlassian connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/atlassian/v1/ops.rs");
                }
            }
        }
        /// Airtable connector contract.
        pub mod airtable {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/airtable/v1/kei.connectors.airtable.v1.rs");
                /// Generated Airtable connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/airtable/v1/ops.rs");
                }
            }
        }
        /// HubSpot connector contract.
        pub mod hubspot {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/hubspot/v1/kei.connectors.hubspot.v1.rs");
                /// Generated HubSpot connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/hubspot/v1/ops.rs");
                }
            }
        }
        /// Notion connector contract.
        pub mod notion {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/notion/v1/kei.connectors.notion.v1.rs");
                /// Generated Notion connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/notion/v1/ops.rs");
                }
            }
        }
        /// Sentry connector contract.
        pub mod sentry {
            /// Version 1.
            pub mod v1 {
                // The generated file includes the matching .serde.rs itself.
                include!("gen/kei/connectors/sentry/v1/kei.connectors.sentry.v1.rs");
                /// Generated Sentry connector operations, driven by the
                /// hand-written connection runtime.
                pub mod ops {
                    include!("kei/connectors/sentry/v1/ops.rs");
                }
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
