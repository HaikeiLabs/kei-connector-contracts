// @generated
impl serde::Serialize for ConfluencePageSummary {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.space_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.ConfluencePageSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.space_id.is_empty() {
            struct_ser.serialize_field("spaceId", &self.space_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ConfluencePageSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "title",
            "space_id",
            "spaceId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Title,
            SpaceId,
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
                            "title" => Ok(GeneratedField::Title),
                            "spaceId" | "space_id" => Ok(GeneratedField::SpaceId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ConfluencePageSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.ConfluencePageSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ConfluencePageSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut title__ = None;
                let mut space_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SpaceId => {
                            if space_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spaceId"));
                            }
                            space_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ConfluencePageSummary {
                    id: id__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    space_id: space_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.ConfluencePageSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetConfluencePageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.page_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.GetConfluencePageRequest", len)?;
        if self.page_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("pageId", ToString::to_string(&self.page_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetConfluencePageRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "page_id",
            "pageId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PageId,
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
                            "pageId" | "page_id" => Ok(GeneratedField::PageId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetConfluencePageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.GetConfluencePageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetConfluencePageRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut page_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PageId => {
                            if page_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageId"));
                            }
                            page_id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetConfluencePageRequest {
                    page_id: page_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.GetConfluencePageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetConfluencePageResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.space_id.is_empty() {
            len += 1;
        }
        if self.parent_id.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.GetConfluencePageResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.space_id.is_empty() {
            struct_ser.serialize_field("spaceId", &self.space_id)?;
        }
        if let Some(v) = self.parent_id.as_ref() {
            struct_ser.serialize_field("parentId", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetConfluencePageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "title",
            "space_id",
            "spaceId",
            "parent_id",
            "parentId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Title,
            SpaceId,
            ParentId,
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
                            "title" => Ok(GeneratedField::Title),
                            "spaceId" | "space_id" => Ok(GeneratedField::SpaceId),
                            "parentId" | "parent_id" => Ok(GeneratedField::ParentId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetConfluencePageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.GetConfluencePageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetConfluencePageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut title__ = None;
                let mut space_id__ = None;
                let mut parent_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SpaceId => {
                            if space_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spaceId"));
                            }
                            space_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ParentId => {
                            if parent_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("parentId"));
                            }
                            parent_id__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetConfluencePageResponse {
                    id: id__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    space_id: space_id__.unwrap_or_default(),
                    parent_id: parent_id__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.GetConfluencePageResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetJiraIssueRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.issue_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.GetJiraIssueRequest", len)?;
        if !self.issue_key.is_empty() {
            struct_ser.serialize_field("issueKey", &self.issue_key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetJiraIssueRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "issue_key",
            "issueKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            IssueKey,
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
                            "issueKey" | "issue_key" => Ok(GeneratedField::IssueKey),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetJiraIssueRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.GetJiraIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetJiraIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut issue_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::IssueKey => {
                            if issue_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueKey"));
                            }
                            issue_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetJiraIssueRequest {
                    issue_key: issue_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.GetJiraIssueRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetJiraIssueResponse {
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
        if !self.key.is_empty() {
            len += 1;
        }
        if self.fields.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.GetJiraIssueResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.key.is_empty() {
            struct_ser.serialize_field("key", &self.key)?;
        }
        if let Some(v) = self.fields.as_ref() {
            struct_ser.serialize_field("fields", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetJiraIssueResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "key",
            "fields",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Key,
            Fields,
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
                            "key" => Ok(GeneratedField::Key),
                            "fields" => Ok(GeneratedField::Fields),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetJiraIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.GetJiraIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetJiraIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut key__ = None;
                let mut fields__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Fields => {
                            if fields__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fields"));
                            }
                            fields__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetJiraIssueResponse {
                    id: id__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    fields: fields__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.GetJiraIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for JiraIssueFields {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.summary.is_empty() {
            len += 1;
        }
        if self.description.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.JiraIssueFields", len)?;
        if !self.summary.is_empty() {
            struct_ser.serialize_field("summary", &self.summary)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for JiraIssueFields {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "summary",
            "description",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Summary,
            Description,
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
                            "summary" => Ok(GeneratedField::Summary),
                            "description" => Ok(GeneratedField::Description),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = JiraIssueFields;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.JiraIssueFields")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<JiraIssueFields, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut summary__ = None;
                let mut description__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Summary => {
                            if summary__.is_some() {
                                return Err(serde::de::Error::duplicate_field("summary"));
                            }
                            summary__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(JiraIssueFields {
                    summary: summary__.unwrap_or_default(),
                    description: description__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.JiraIssueFields", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for JiraIssueSummary {
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
        if !self.key.is_empty() {
            len += 1;
        }
        if self.fields.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.JiraIssueSummary", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.key.is_empty() {
            struct_ser.serialize_field("key", &self.key)?;
        }
        if let Some(v) = self.fields.as_ref() {
            struct_ser.serialize_field("fields", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for JiraIssueSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "key",
            "fields",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Key,
            Fields,
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
                            "key" => Ok(GeneratedField::Key),
                            "fields" => Ok(GeneratedField::Fields),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = JiraIssueSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.JiraIssueSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<JiraIssueSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut key__ = None;
                let mut fields__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Fields => {
                            if fields__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fields"));
                            }
                            fields__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(JiraIssueSummary {
                    id: id__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    fields: fields__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.JiraIssueSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListConfluencePagesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.space_id != 0 {
            len += 1;
        }
        if self.limit.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.ListConfluencePagesRequest", len)?;
        if self.space_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("spaceId", ToString::to_string(&self.space_id).as_str())?;
        }
        if let Some(v) = self.limit.as_ref() {
            struct_ser.serialize_field("limit", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListConfluencePagesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "space_id",
            "spaceId",
            "limit",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpaceId,
            Limit,
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
                            "spaceId" | "space_id" => Ok(GeneratedField::SpaceId),
                            "limit" => Ok(GeneratedField::Limit),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListConfluencePagesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.ListConfluencePagesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListConfluencePagesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut space_id__ = None;
                let mut limit__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpaceId => {
                            if space_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spaceId"));
                            }
                            space_id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Limit => {
                            if limit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("limit"));
                            }
                            limit__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListConfluencePagesRequest {
                    space_id: space_id__.unwrap_or_default(),
                    limit: limit__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.ListConfluencePagesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListConfluencePagesResponse {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.ListConfluencePagesResponse", len)?;
        if !self.results.is_empty() {
            struct_ser.serialize_field("results", &self.results)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListConfluencePagesResponse {
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
            type Value = ListConfluencePagesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.ListConfluencePagesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListConfluencePagesResponse, V::Error>
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
                Ok(ListConfluencePagesResponse {
                    results: results__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.ListConfluencePagesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SearchJiraIssuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.jql.is_some() {
            len += 1;
        }
        if self.next_page_token.is_some() {
            len += 1;
        }
        if self.max_results.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.SearchJiraIssuesRequest", len)?;
        if let Some(v) = self.jql.as_ref() {
            struct_ser.serialize_field("jql", v)?;
        }
        if let Some(v) = self.next_page_token.as_ref() {
            struct_ser.serialize_field("nextPageToken", v)?;
        }
        if let Some(v) = self.max_results.as_ref() {
            struct_ser.serialize_field("maxResults", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SearchJiraIssuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "jql",
            "next_page_token",
            "nextPageToken",
            "max_results",
            "maxResults",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Jql,
            NextPageToken,
            MaxResults,
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
                            "jql" => Ok(GeneratedField::Jql),
                            "nextPageToken" | "next_page_token" => Ok(GeneratedField::NextPageToken),
                            "maxResults" | "max_results" => Ok(GeneratedField::MaxResults),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SearchJiraIssuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.SearchJiraIssuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SearchJiraIssuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut jql__ = None;
                let mut next_page_token__ = None;
                let mut max_results__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Jql => {
                            if jql__.is_some() {
                                return Err(serde::de::Error::duplicate_field("jql"));
                            }
                            jql__ = map_.next_value()?;
                        }
                        GeneratedField::NextPageToken => {
                            if next_page_token__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPageToken"));
                            }
                            next_page_token__ = map_.next_value()?;
                        }
                        GeneratedField::MaxResults => {
                            if max_results__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxResults"));
                            }
                            max_results__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SearchJiraIssuesRequest {
                    jql: jql__,
                    next_page_token: next_page_token__,
                    max_results: max_results__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.SearchJiraIssuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SearchJiraIssuesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.issues.is_empty() {
            len += 1;
        }
        if self.next_page_token.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.atlassian.v1.SearchJiraIssuesResponse", len)?;
        if !self.issues.is_empty() {
            struct_ser.serialize_field("issues", &self.issues)?;
        }
        if let Some(v) = self.next_page_token.as_ref() {
            struct_ser.serialize_field("nextPageToken", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SearchJiraIssuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "issues",
            "next_page_token",
            "nextPageToken",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Issues,
            NextPageToken,
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
                            "issues" => Ok(GeneratedField::Issues),
                            "nextPageToken" | "next_page_token" => Ok(GeneratedField::NextPageToken),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SearchJiraIssuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.atlassian.v1.SearchJiraIssuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SearchJiraIssuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut issues__ = None;
                let mut next_page_token__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Issues => {
                            if issues__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issues"));
                            }
                            issues__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextPageToken => {
                            if next_page_token__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPageToken"));
                            }
                            next_page_token__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SearchJiraIssuesResponse {
                    issues: issues__.unwrap_or_default(),
                    next_page_token: next_page_token__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.atlassian.v1.SearchJiraIssuesResponse", FIELDS, GeneratedVisitor)
    }
}
