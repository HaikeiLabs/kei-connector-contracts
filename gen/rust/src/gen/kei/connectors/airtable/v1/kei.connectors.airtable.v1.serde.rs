// @generated
impl serde::Serialize for GetRecordRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.base_id.is_empty() {
            len += 1;
        }
        if !self.table_id_or_name.is_empty() {
            len += 1;
        }
        if !self.record_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.airtable.v1.GetRecordRequest", len)?;
        if !self.base_id.is_empty() {
            struct_ser.serialize_field("baseId", &self.base_id)?;
        }
        if !self.table_id_or_name.is_empty() {
            struct_ser.serialize_field("tableIdOrName", &self.table_id_or_name)?;
        }
        if !self.record_id.is_empty() {
            struct_ser.serialize_field("recordId", &self.record_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRecordRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "base_id",
            "baseId",
            "table_id_or_name",
            "tableIdOrName",
            "record_id",
            "recordId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            BaseId,
            TableIdOrName,
            RecordId,
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
                            "baseId" | "base_id" => Ok(GeneratedField::BaseId),
                            "tableIdOrName" | "table_id_or_name" => Ok(GeneratedField::TableIdOrName),
                            "recordId" | "record_id" => Ok(GeneratedField::RecordId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRecordRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.airtable.v1.GetRecordRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRecordRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut base_id__ = None;
                let mut table_id_or_name__ = None;
                let mut record_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::BaseId => {
                            if base_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseId"));
                            }
                            base_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TableIdOrName => {
                            if table_id_or_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tableIdOrName"));
                            }
                            table_id_or_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::RecordId => {
                            if record_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("recordId"));
                            }
                            record_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRecordRequest {
                    base_id: base_id__.unwrap_or_default(),
                    table_id_or_name: table_id_or_name__.unwrap_or_default(),
                    record_id: record_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.airtable.v1.GetRecordRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRecordResponse {
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
        if self.fields.is_some() {
            len += 1;
        }
        if !self.created_time.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.airtable.v1.GetRecordResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if let Some(v) = self.fields.as_ref() {
            struct_ser.serialize_field("fields", v)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRecordResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "fields",
            "created_time",
            "createdTime",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Fields,
            CreatedTime,
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
                            "fields" => Ok(GeneratedField::Fields),
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRecordResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.airtable.v1.GetRecordResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRecordResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut fields__ = None;
                let mut created_time__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Fields => {
                            if fields__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fields"));
                            }
                            fields__ = map_.next_value()?;
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRecordResponse {
                    id: id__.unwrap_or_default(),
                    fields: fields__,
                    created_time: created_time__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.airtable.v1.GetRecordResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListRecordsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.base_id.is_empty() {
            len += 1;
        }
        if !self.table_id_or_name.is_empty() {
            len += 1;
        }
        if self.offset.is_some() {
            len += 1;
        }
        if self.page_size.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.airtable.v1.ListRecordsRequest", len)?;
        if !self.base_id.is_empty() {
            struct_ser.serialize_field("baseId", &self.base_id)?;
        }
        if !self.table_id_or_name.is_empty() {
            struct_ser.serialize_field("tableIdOrName", &self.table_id_or_name)?;
        }
        if let Some(v) = self.offset.as_ref() {
            struct_ser.serialize_field("offset", v)?;
        }
        if let Some(v) = self.page_size.as_ref() {
            struct_ser.serialize_field("pageSize", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListRecordsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "base_id",
            "baseId",
            "table_id_or_name",
            "tableIdOrName",
            "offset",
            "page_size",
            "pageSize",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            BaseId,
            TableIdOrName,
            Offset,
            PageSize,
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
                            "baseId" | "base_id" => Ok(GeneratedField::BaseId),
                            "tableIdOrName" | "table_id_or_name" => Ok(GeneratedField::TableIdOrName),
                            "offset" => Ok(GeneratedField::Offset),
                            "pageSize" | "page_size" => Ok(GeneratedField::PageSize),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListRecordsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.airtable.v1.ListRecordsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListRecordsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut base_id__ = None;
                let mut table_id_or_name__ = None;
                let mut offset__ = None;
                let mut page_size__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::BaseId => {
                            if base_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseId"));
                            }
                            base_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TableIdOrName => {
                            if table_id_or_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tableIdOrName"));
                            }
                            table_id_or_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Offset => {
                            if offset__.is_some() {
                                return Err(serde::de::Error::duplicate_field("offset"));
                            }
                            offset__ = map_.next_value()?;
                        }
                        GeneratedField::PageSize => {
                            if page_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageSize"));
                            }
                            page_size__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListRecordsRequest {
                    base_id: base_id__.unwrap_or_default(),
                    table_id_or_name: table_id_or_name__.unwrap_or_default(),
                    offset: offset__,
                    page_size: page_size__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.airtable.v1.ListRecordsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListRecordsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.records.is_empty() {
            len += 1;
        }
        if self.offset.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.airtable.v1.ListRecordsResponse", len)?;
        if !self.records.is_empty() {
            struct_ser.serialize_field("records", &self.records)?;
        }
        if let Some(v) = self.offset.as_ref() {
            struct_ser.serialize_field("offset", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListRecordsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "records",
            "offset",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Records,
            Offset,
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
                            "records" => Ok(GeneratedField::Records),
                            "offset" => Ok(GeneratedField::Offset),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListRecordsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.airtable.v1.ListRecordsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListRecordsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut records__ = None;
                let mut offset__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Records => {
                            if records__.is_some() {
                                return Err(serde::de::Error::duplicate_field("records"));
                            }
                            records__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Offset => {
                            if offset__.is_some() {
                                return Err(serde::de::Error::duplicate_field("offset"));
                            }
                            offset__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListRecordsResponse {
                    records: records__.unwrap_or_default(),
                    offset: offset__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.airtable.v1.ListRecordsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RecordSummary {
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
        if self.fields.is_some() {
            len += 1;
        }
        if !self.created_time.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.airtable.v1.RecordSummary", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if let Some(v) = self.fields.as_ref() {
            struct_ser.serialize_field("fields", v)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RecordSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "fields",
            "created_time",
            "createdTime",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Fields,
            CreatedTime,
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
                            "fields" => Ok(GeneratedField::Fields),
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RecordSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.airtable.v1.RecordSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RecordSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut fields__ = None;
                let mut created_time__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Fields => {
                            if fields__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fields"));
                            }
                            fields__ = map_.next_value()?;
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(RecordSummary {
                    id: id__.unwrap_or_default(),
                    fields: fields__,
                    created_time: created_time__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.airtable.v1.RecordSummary", FIELDS, GeneratedVisitor)
    }
}
