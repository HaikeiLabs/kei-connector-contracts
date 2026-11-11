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

//! The same sample JSON the Go registry test exercises: int64 fields must
//! survive proto-JSON unmarshal exactly (2^53 + 1 would lose precision if a
//! float64 were involved), unknown fields are discarded, and a JSON null
//! clears an optional field.

use haikeilabs_connector_contracts::kei::connectors::github::v1::GetIssueResponse;

const SAMPLE_JSON: &str = r#"{"id": 9007199254740993, "number": 7, "title": "t", "state": "open", "body": null, "html_url": "https://example.invalid/x", "extra_field": 1}"#;

#[test]
fn sample_json_round_trips_int64_exactly() {
    let msg: GetIssueResponse =
        serde_json::from_str(SAMPLE_JSON).expect("deserialize sample JSON");

    assert_eq!(msg.id, 9007199254740993);
    assert_eq!(msg.number, 7);
    assert_eq!(msg.title, "t");
    assert_eq!(msg.state, "open");
    assert_eq!(msg.body, None);

    let back = serde_json::to_string(&msg).expect("serialize sample message");
    let reparsed: GetIssueResponse =
        serde_json::from_str(&back).expect("reparse serialized message");
    assert_eq!(reparsed.id, 9007199254740993);
    assert_eq!(reparsed.number, 7);
}
