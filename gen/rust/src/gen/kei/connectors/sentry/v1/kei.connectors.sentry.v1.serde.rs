// @generated
impl serde::Serialize for GetIssueRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.organization_slug.is_empty() {
            len += 1;
        }
        if !self.issue_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.sentry.v1.GetIssueRequest", len)?;
        if !self.organization_slug.is_empty() {
            struct_ser.serialize_field("organizationSlug", &self.organization_slug)?;
        }
        if !self.issue_id.is_empty() {
            struct_ser.serialize_field("issueId", &self.issue_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetIssueRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "organization_slug",
            "organizationSlug",
            "issue_id",
            "issueId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            OrganizationSlug,
            IssueId,
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
                            "organizationSlug" | "organization_slug" => Ok(GeneratedField::OrganizationSlug),
                            "issueId" | "issue_id" => Ok(GeneratedField::IssueId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetIssueRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.sentry.v1.GetIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut organization_slug__ = None;
                let mut issue_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::OrganizationSlug => {
                            if organization_slug__.is_some() {
                                return Err(serde::de::Error::duplicate_field("organizationSlug"));
                            }
                            organization_slug__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IssueId => {
                            if issue_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueId"));
                            }
                            issue_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetIssueRequest {
                    organization_slug: organization_slug__.unwrap_or_default(),
                    issue_id: issue_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.sentry.v1.GetIssueRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetIssueResponse {
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
        if !self.short_id.is_empty() {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.status.is_empty() {
            len += 1;
        }
        if !self.level.is_empty() {
            len += 1;
        }
        if self.share_id.is_some() {
            len += 1;
        }
        if !self.permalink.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.sentry.v1.GetIssueResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.short_id.is_empty() {
            struct_ser.serialize_field("shortId", &self.short_id)?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.status.is_empty() {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if !self.level.is_empty() {
            struct_ser.serialize_field("level", &self.level)?;
        }
        if let Some(v) = self.share_id.as_ref() {
            struct_ser.serialize_field("shareId", v)?;
        }
        if !self.permalink.is_empty() {
            struct_ser.serialize_field("permalink", &self.permalink)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetIssueResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "short_id",
            "shortId",
            "title",
            "status",
            "level",
            "share_id",
            "shareId",
            "permalink",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            ShortId,
            Title,
            Status,
            Level,
            ShareId,
            Permalink,
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
                            "shortId" | "short_id" => Ok(GeneratedField::ShortId),
                            "title" => Ok(GeneratedField::Title),
                            "status" => Ok(GeneratedField::Status),
                            "level" => Ok(GeneratedField::Level),
                            "shareId" | "share_id" => Ok(GeneratedField::ShareId),
                            "permalink" => Ok(GeneratedField::Permalink),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.sentry.v1.GetIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut short_id__ = None;
                let mut title__ = None;
                let mut status__ = None;
                let mut level__ = None;
                let mut share_id__ = None;
                let mut permalink__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ShortId => {
                            if short_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shortId"));
                            }
                            short_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Level => {
                            if level__.is_some() {
                                return Err(serde::de::Error::duplicate_field("level"));
                            }
                            level__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ShareId => {
                            if share_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shareId"));
                            }
                            share_id__ = map_.next_value()?;
                        }
                        GeneratedField::Permalink => {
                            if permalink__.is_some() {
                                return Err(serde::de::Error::duplicate_field("permalink"));
                            }
                            permalink__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetIssueResponse {
                    id: id__.unwrap_or_default(),
                    short_id: short_id__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    status: status__.unwrap_or_default(),
                    level: level__.unwrap_or_default(),
                    share_id: share_id__,
                    permalink: permalink__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.sentry.v1.GetIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for IssueSummary {
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
        if !self.short_id.is_empty() {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.status.is_empty() {
            len += 1;
        }
        if !self.level.is_empty() {
            len += 1;
        }
        if self.share_id.is_some() {
            len += 1;
        }
        if !self.permalink.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.sentry.v1.IssueSummary", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.short_id.is_empty() {
            struct_ser.serialize_field("shortId", &self.short_id)?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.status.is_empty() {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if !self.level.is_empty() {
            struct_ser.serialize_field("level", &self.level)?;
        }
        if let Some(v) = self.share_id.as_ref() {
            struct_ser.serialize_field("shareId", v)?;
        }
        if !self.permalink.is_empty() {
            struct_ser.serialize_field("permalink", &self.permalink)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for IssueSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "short_id",
            "shortId",
            "title",
            "status",
            "level",
            "share_id",
            "shareId",
            "permalink",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            ShortId,
            Title,
            Status,
            Level,
            ShareId,
            Permalink,
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
                            "shortId" | "short_id" => Ok(GeneratedField::ShortId),
                            "title" => Ok(GeneratedField::Title),
                            "status" => Ok(GeneratedField::Status),
                            "level" => Ok(GeneratedField::Level),
                            "shareId" | "share_id" => Ok(GeneratedField::ShareId),
                            "permalink" => Ok(GeneratedField::Permalink),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = IssueSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.sentry.v1.IssueSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<IssueSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut short_id__ = None;
                let mut title__ = None;
                let mut status__ = None;
                let mut level__ = None;
                let mut share_id__ = None;
                let mut permalink__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ShortId => {
                            if short_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shortId"));
                            }
                            short_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Level => {
                            if level__.is_some() {
                                return Err(serde::de::Error::duplicate_field("level"));
                            }
                            level__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ShareId => {
                            if share_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shareId"));
                            }
                            share_id__ = map_.next_value()?;
                        }
                        GeneratedField::Permalink => {
                            if permalink__.is_some() {
                                return Err(serde::de::Error::duplicate_field("permalink"));
                            }
                            permalink__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(IssueSummary {
                    id: id__.unwrap_or_default(),
                    short_id: short_id__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    status: status__.unwrap_or_default(),
                    level: level__.unwrap_or_default(),
                    share_id: share_id__,
                    permalink: permalink__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.sentry.v1.IssueSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListIssuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.organization_slug.is_empty() {
            len += 1;
        }
        if self.query.is_some() {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.sentry.v1.ListIssuesRequest", len)?;
        if !self.organization_slug.is_empty() {
            struct_ser.serialize_field("organizationSlug", &self.organization_slug)?;
        }
        if let Some(v) = self.query.as_ref() {
            struct_ser.serialize_field("query", v)?;
        }
        if let Some(v) = self.per_page.as_ref() {
            struct_ser.serialize_field("perPage", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListIssuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "organization_slug",
            "organizationSlug",
            "query",
            "per_page",
            "perPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            OrganizationSlug,
            Query,
            PerPage,
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
                            "organizationSlug" | "organization_slug" => Ok(GeneratedField::OrganizationSlug),
                            "query" => Ok(GeneratedField::Query),
                            "perPage" | "per_page" => Ok(GeneratedField::PerPage),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListIssuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.sentry.v1.ListIssuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListIssuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut organization_slug__ = None;
                let mut query__ = None;
                let mut per_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::OrganizationSlug => {
                            if organization_slug__.is_some() {
                                return Err(serde::de::Error::duplicate_field("organizationSlug"));
                            }
                            organization_slug__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Query => {
                            if query__.is_some() {
                                return Err(serde::de::Error::duplicate_field("query"));
                            }
                            query__ = map_.next_value()?;
                        }
                        GeneratedField::PerPage => {
                            if per_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("perPage"));
                            }
                            per_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListIssuesRequest {
                    organization_slug: organization_slug__.unwrap_or_default(),
                    query: query__,
                    per_page: per_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.sentry.v1.ListIssuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListIssuesResponse {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.sentry.v1.ListIssuesResponse", len)?;
        if !self.issues.is_empty() {
            struct_ser.serialize_field("issues", &self.issues)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListIssuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "issues",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Issues,
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
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListIssuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.sentry.v1.ListIssuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListIssuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut issues__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Issues => {
                            if issues__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issues"));
                            }
                            issues__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListIssuesResponse {
                    issues: issues__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.sentry.v1.ListIssuesResponse", FIELDS, GeneratedVisitor)
    }
}
