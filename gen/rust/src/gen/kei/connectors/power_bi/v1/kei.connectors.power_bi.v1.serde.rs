// @generated
impl serde::Serialize for DatasetSummary {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if !self.name.is_empty() {
            len += 1;
        }
        if self.configured_by.is_some() {
            len += 1;
        }
        if self.target_storage_mode.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.DatasetSummary", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if let Some(v) = self.configured_by.as_ref() {
            struct_ser.serialize_field("ConfiguredBy", v)?;
        }
        if let Some(v) = self.target_storage_mode.as_ref() {
            struct_ser.serialize_field("TargetStorageMode", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DatasetSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
            "configured_by",
            "ConfiguredBy",
            "target_storage_mode",
            "TargetStorageMode",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
            ConfiguredBy,
            TargetStorageMode,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "id" => Ok(GeneratedField::Id),
                            "name" => Ok(GeneratedField::Name),
                            "ConfiguredBy" | "configured_by" => Ok(GeneratedField::ConfiguredBy),
                            "TargetStorageMode" | "target_storage_mode" => Ok(GeneratedField::TargetStorageMode),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DatasetSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.DatasetSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DatasetSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
                let mut configured_by__ = None;
                let mut target_storage_mode__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ConfiguredBy => {
                            if configured_by__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ConfiguredBy"));
                            }
                            configured_by__ = map_.next_value()?;
                        }
                        GeneratedField::TargetStorageMode => {
                            if target_storage_mode__.is_some() {
                                return Err(serde::de::Error::duplicate_field("TargetStorageMode"));
                            }
                            target_storage_mode__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(DatasetSummary {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    configured_by: configured_by__,
                    target_storage_mode: target_storage_mode__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.DatasetSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ExecuteQueriesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.group_id.is_empty() {
            len += 1;
        }
        if !self.dataset_id.is_empty() {
            len += 1;
        }
        if !self.queries.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.ExecuteQueriesRequest", len)?;
        if !self.group_id.is_empty() {
            struct_ser.serialize_field("groupId", &self.group_id)?;
        }
        if !self.dataset_id.is_empty() {
            struct_ser.serialize_field("datasetId", &self.dataset_id)?;
        }
        if !self.queries.is_empty() {
            struct_ser.serialize_field("queries", &self.queries)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ExecuteQueriesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "group_id",
            "groupId",
            "dataset_id",
            "datasetId",
            "queries",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            GroupId,
            DatasetId,
            Queries,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "groupId" | "group_id" => Ok(GeneratedField::GroupId),
                            "datasetId" | "dataset_id" => Ok(GeneratedField::DatasetId),
                            "queries" => Ok(GeneratedField::Queries),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ExecuteQueriesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.ExecuteQueriesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ExecuteQueriesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut group_id__ = None;
                let mut dataset_id__ = None;
                let mut queries__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::GroupId => {
                            if group_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("groupId"));
                            }
                            group_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::DatasetId => {
                            if dataset_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("datasetId"));
                            }
                            dataset_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Queries => {
                            if queries__.is_some() {
                                return Err(serde::de::Error::duplicate_field("queries"));
                            }
                            queries__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ExecuteQueriesRequest {
                    group_id: group_id__.unwrap_or_default(),
                    dataset_id: dataset_id__.unwrap_or_default(),
                    queries: queries__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.ExecuteQueriesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ExecuteQueriesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.results.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.ExecuteQueriesResponse", len)?;
        if !self.results.is_empty() {
            struct_ser.serialize_field("results", &self.results)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ExecuteQueriesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "results",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Results,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "results" => Ok(GeneratedField::Results),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ExecuteQueriesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.ExecuteQueriesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ExecuteQueriesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut results__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Results => {
                            if results__.is_some() {
                                return Err(serde::de::Error::duplicate_field("results"));
                            }
                            results__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ExecuteQueriesResponse {
                    results: results__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.ExecuteQueriesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListDatasetsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.group_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.ListDatasetsRequest", len)?;
        if !self.group_id.is_empty() {
            struct_ser.serialize_field("groupId", &self.group_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListDatasetsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "group_id",
            "groupId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            GroupId,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "groupId" | "group_id" => Ok(GeneratedField::GroupId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListDatasetsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.ListDatasetsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListDatasetsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut group_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::GroupId => {
                            if group_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("groupId"));
                            }
                            group_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListDatasetsRequest {
                    group_id: group_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.ListDatasetsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListDatasetsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.value.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.ListDatasetsResponse", len)?;
        if !self.value.is_empty() {
            struct_ser.serialize_field("value", &self.value)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListDatasetsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "value",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Value,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "value" => Ok(GeneratedField::Value),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListDatasetsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.ListDatasetsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListDatasetsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut value__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListDatasetsResponse {
                    value: value__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.ListDatasetsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for QueryDefinition {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.query.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.QueryDefinition", len)?;
        if !self.query.is_empty() {
            struct_ser.serialize_field("query", &self.query)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for QueryDefinition {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "query",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Query,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "query" => Ok(GeneratedField::Query),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = QueryDefinition;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.QueryDefinition")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<QueryDefinition, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut query__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Query => {
                            if query__.is_some() {
                                return Err(serde::de::Error::duplicate_field("query"));
                            }
                            query__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(QueryDefinition {
                    query: query__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.QueryDefinition", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for QueryResult {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.tables.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.QueryResult", len)?;
        if !self.tables.is_empty() {
            struct_ser.serialize_field("tables", &self.tables)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for QueryResult {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "tables",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Tables,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "tables" => Ok(GeneratedField::Tables),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = QueryResult;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.QueryResult")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<QueryResult, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut tables__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Tables => {
                            if tables__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tables"));
                            }
                            tables__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(QueryResult {
                    tables: tables__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.QueryResult", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for QueryTable {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.rows.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.power_bi.v1.QueryTable", len)?;
        if !self.rows.is_empty() {
            struct_ser.serialize_field("rows", &self.rows)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for QueryTable {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "rows",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Rows,
            __SkipField__,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "rows" => Ok(GeneratedField::Rows),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = QueryTable;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.power_bi.v1.QueryTable")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<QueryTable, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut rows__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Rows => {
                            if rows__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rows"));
                            }
                            rows__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(QueryTable {
                    rows: rows__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.power_bi.v1.QueryTable", FIELDS, GeneratedVisitor)
    }
}
