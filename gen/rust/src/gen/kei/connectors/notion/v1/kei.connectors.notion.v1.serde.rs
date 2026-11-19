// @generated
impl serde::Serialize for CreatePageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.parent.is_some() {
            len += 1;
        }
        if self.properties.is_some() {
            len += 1;
        }
        if !self.children.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.CreatePageRequest", len)?;
        if let Some(v) = self.parent.as_ref() {
            struct_ser.serialize_field("parent", v)?;
        }
        if let Some(v) = self.properties.as_ref() {
            struct_ser.serialize_field("properties", v)?;
        }
        if !self.children.is_empty() {
            struct_ser.serialize_field("children", &self.children)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePageRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "parent",
            "properties",
            "children",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Parent,
            Properties,
            Children,
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
                            "parent" => Ok(GeneratedField::Parent),
                            "properties" => Ok(GeneratedField::Properties),
                            "children" => Ok(GeneratedField::Children),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreatePageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.CreatePageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePageRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut parent__ = None;
                let mut properties__ = None;
                let mut children__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Parent => {
                            if parent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("parent"));
                            }
                            parent__ = map_.next_value()?;
                        }
                        GeneratedField::Properties => {
                            if properties__.is_some() {
                                return Err(serde::de::Error::duplicate_field("properties"));
                            }
                            properties__ = map_.next_value()?;
                        }
                        GeneratedField::Children => {
                            if children__.is_some() {
                                return Err(serde::de::Error::duplicate_field("children"));
                            }
                            children__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreatePageRequest {
                    parent: parent__,
                    properties: properties__,
                    children: children__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.CreatePageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreatePageResponse {
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
        if !self.created_time.is_empty() {
            len += 1;
        }
        if !self.last_edited_time.is_empty() {
            len += 1;
        }
        if self.archived {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if self.public_url.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.CreatePageResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        if !self.last_edited_time.is_empty() {
            struct_ser.serialize_field("lastEditedTime", &self.last_edited_time)?;
        }
        if self.archived {
            struct_ser.serialize_field("archived", &self.archived)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if let Some(v) = self.public_url.as_ref() {
            struct_ser.serialize_field("publicUrl", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "created_time",
            "createdTime",
            "last_edited_time",
            "lastEditedTime",
            "archived",
            "url",
            "public_url",
            "publicUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            CreatedTime,
            LastEditedTime,
            Archived,
            Url,
            PublicUrl,
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
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            "lastEditedTime" | "last_edited_time" => Ok(GeneratedField::LastEditedTime),
                            "archived" => Ok(GeneratedField::Archived),
                            "url" => Ok(GeneratedField::Url),
                            "publicUrl" | "public_url" => Ok(GeneratedField::PublicUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreatePageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.CreatePageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut created_time__ = None;
                let mut last_edited_time__ = None;
                let mut archived__ = None;
                let mut url__ = None;
                let mut public_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LastEditedTime => {
                            if last_edited_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastEditedTime"));
                            }
                            last_edited_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Archived => {
                            if archived__.is_some() {
                                return Err(serde::de::Error::duplicate_field("archived"));
                            }
                            archived__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicUrl => {
                            if public_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicUrl"));
                            }
                            public_url__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreatePageResponse {
                    id: id__.unwrap_or_default(),
                    created_time: created_time__.unwrap_or_default(),
                    last_edited_time: last_edited_time__.unwrap_or_default(),
                    archived: archived__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                    public_url: public_url__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.CreatePageResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.page_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.GetPageRequest", len)?;
        if !self.page_id.is_empty() {
            struct_ser.serialize_field("pageId", &self.page_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPageRequest {
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
            type Value = GetPageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.GetPageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPageRequest, V::Error>
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
                            page_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetPageRequest {
                    page_id: page_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.GetPageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPageResponse {
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
        if !self.created_time.is_empty() {
            len += 1;
        }
        if !self.last_edited_time.is_empty() {
            len += 1;
        }
        if self.archived {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if self.public_url.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.GetPageResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        if !self.last_edited_time.is_empty() {
            struct_ser.serialize_field("lastEditedTime", &self.last_edited_time)?;
        }
        if self.archived {
            struct_ser.serialize_field("archived", &self.archived)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if let Some(v) = self.public_url.as_ref() {
            struct_ser.serialize_field("publicUrl", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "created_time",
            "createdTime",
            "last_edited_time",
            "lastEditedTime",
            "archived",
            "url",
            "public_url",
            "publicUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            CreatedTime,
            LastEditedTime,
            Archived,
            Url,
            PublicUrl,
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
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            "lastEditedTime" | "last_edited_time" => Ok(GeneratedField::LastEditedTime),
                            "archived" => Ok(GeneratedField::Archived),
                            "url" => Ok(GeneratedField::Url),
                            "publicUrl" | "public_url" => Ok(GeneratedField::PublicUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetPageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.GetPageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut created_time__ = None;
                let mut last_edited_time__ = None;
                let mut archived__ = None;
                let mut url__ = None;
                let mut public_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LastEditedTime => {
                            if last_edited_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastEditedTime"));
                            }
                            last_edited_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Archived => {
                            if archived__.is_some() {
                                return Err(serde::de::Error::duplicate_field("archived"));
                            }
                            archived__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicUrl => {
                            if public_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicUrl"));
                            }
                            public_url__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetPageResponse {
                    id: id__.unwrap_or_default(),
                    created_time: created_time__.unwrap_or_default(),
                    last_edited_time: last_edited_time__.unwrap_or_default(),
                    archived: archived__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                    public_url: public_url__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.GetPageResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PageSummary {
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
        if !self.created_time.is_empty() {
            len += 1;
        }
        if !self.last_edited_time.is_empty() {
            len += 1;
        }
        if self.archived {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if self.public_url.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.PageSummary", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        if !self.last_edited_time.is_empty() {
            struct_ser.serialize_field("lastEditedTime", &self.last_edited_time)?;
        }
        if self.archived {
            struct_ser.serialize_field("archived", &self.archived)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if let Some(v) = self.public_url.as_ref() {
            struct_ser.serialize_field("publicUrl", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PageSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "created_time",
            "createdTime",
            "last_edited_time",
            "lastEditedTime",
            "archived",
            "url",
            "public_url",
            "publicUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            CreatedTime,
            LastEditedTime,
            Archived,
            Url,
            PublicUrl,
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
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            "lastEditedTime" | "last_edited_time" => Ok(GeneratedField::LastEditedTime),
                            "archived" => Ok(GeneratedField::Archived),
                            "url" => Ok(GeneratedField::Url),
                            "publicUrl" | "public_url" => Ok(GeneratedField::PublicUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PageSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.PageSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PageSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut created_time__ = None;
                let mut last_edited_time__ = None;
                let mut archived__ = None;
                let mut url__ = None;
                let mut public_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LastEditedTime => {
                            if last_edited_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastEditedTime"));
                            }
                            last_edited_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Archived => {
                            if archived__.is_some() {
                                return Err(serde::de::Error::duplicate_field("archived"));
                            }
                            archived__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicUrl => {
                            if public_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicUrl"));
                            }
                            public_url__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(PageSummary {
                    id: id__.unwrap_or_default(),
                    created_time: created_time__.unwrap_or_default(),
                    last_edited_time: last_edited_time__.unwrap_or_default(),
                    archived: archived__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                    public_url: public_url__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.PageSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for QueryDatabaseRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.database_id.is_empty() {
            len += 1;
        }
        if self.start_cursor.is_some() {
            len += 1;
        }
        if self.page_size.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.QueryDatabaseRequest", len)?;
        if !self.database_id.is_empty() {
            struct_ser.serialize_field("databaseId", &self.database_id)?;
        }
        if let Some(v) = self.start_cursor.as_ref() {
            struct_ser.serialize_field("startCursor", v)?;
        }
        if let Some(v) = self.page_size.as_ref() {
            struct_ser.serialize_field("pageSize", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for QueryDatabaseRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "database_id",
            "databaseId",
            "start_cursor",
            "startCursor",
            "page_size",
            "pageSize",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            DatabaseId,
            StartCursor,
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
                            "databaseId" | "database_id" => Ok(GeneratedField::DatabaseId),
                            "startCursor" | "start_cursor" => Ok(GeneratedField::StartCursor),
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
            type Value = QueryDatabaseRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.QueryDatabaseRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<QueryDatabaseRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut database_id__ = None;
                let mut start_cursor__ = None;
                let mut page_size__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::DatabaseId => {
                            if database_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("databaseId"));
                            }
                            database_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::StartCursor => {
                            if start_cursor__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startCursor"));
                            }
                            start_cursor__ = map_.next_value()?;
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
                Ok(QueryDatabaseRequest {
                    database_id: database_id__.unwrap_or_default(),
                    start_cursor: start_cursor__,
                    page_size: page_size__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.QueryDatabaseRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for QueryDatabaseResponse {
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
        if self.next_cursor.is_some() {
            len += 1;
        }
        if self.has_more {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.QueryDatabaseResponse", len)?;
        if !self.results.is_empty() {
            struct_ser.serialize_field("results", &self.results)?;
        }
        if let Some(v) = self.next_cursor.as_ref() {
            struct_ser.serialize_field("nextCursor", v)?;
        }
        if self.has_more {
            struct_ser.serialize_field("hasMore", &self.has_more)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for QueryDatabaseResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "results",
            "next_cursor",
            "nextCursor",
            "has_more",
            "hasMore",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Results,
            NextCursor,
            HasMore,
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
                            "nextCursor" | "next_cursor" => Ok(GeneratedField::NextCursor),
                            "hasMore" | "has_more" => Ok(GeneratedField::HasMore),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = QueryDatabaseResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.QueryDatabaseResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<QueryDatabaseResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut results__ = None;
                let mut next_cursor__ = None;
                let mut has_more__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Results => {
                            if results__.is_some() {
                                return Err(serde::de::Error::duplicate_field("results"));
                            }
                            results__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextCursor => {
                            if next_cursor__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextCursor"));
                            }
                            next_cursor__ = map_.next_value()?;
                        }
                        GeneratedField::HasMore => {
                            if has_more__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hasMore"));
                            }
                            has_more__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(QueryDatabaseResponse {
                    results: results__.unwrap_or_default(),
                    next_cursor: next_cursor__,
                    has_more: has_more__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.QueryDatabaseResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdatePageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.page_id.is_empty() {
            len += 1;
        }
        if self.properties.is_some() {
            len += 1;
        }
        if self.archived.is_some() {
            len += 1;
        }
        if self.in_trash.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.UpdatePageRequest", len)?;
        if !self.page_id.is_empty() {
            struct_ser.serialize_field("pageId", &self.page_id)?;
        }
        if let Some(v) = self.properties.as_ref() {
            struct_ser.serialize_field("properties", v)?;
        }
        if let Some(v) = self.archived.as_ref() {
            struct_ser.serialize_field("archived", v)?;
        }
        if let Some(v) = self.in_trash.as_ref() {
            struct_ser.serialize_field("inTrash", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdatePageRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "page_id",
            "pageId",
            "properties",
            "archived",
            "in_trash",
            "inTrash",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PageId,
            Properties,
            Archived,
            InTrash,
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
                            "properties" => Ok(GeneratedField::Properties),
                            "archived" => Ok(GeneratedField::Archived),
                            "inTrash" | "in_trash" => Ok(GeneratedField::InTrash),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdatePageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.UpdatePageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdatePageRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut page_id__ = None;
                let mut properties__ = None;
                let mut archived__ = None;
                let mut in_trash__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PageId => {
                            if page_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageId"));
                            }
                            page_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Properties => {
                            if properties__.is_some() {
                                return Err(serde::de::Error::duplicate_field("properties"));
                            }
                            properties__ = map_.next_value()?;
                        }
                        GeneratedField::Archived => {
                            if archived__.is_some() {
                                return Err(serde::de::Error::duplicate_field("archived"));
                            }
                            archived__ = map_.next_value()?;
                        }
                        GeneratedField::InTrash => {
                            if in_trash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("inTrash"));
                            }
                            in_trash__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdatePageRequest {
                    page_id: page_id__.unwrap_or_default(),
                    properties: properties__,
                    archived: archived__,
                    in_trash: in_trash__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.UpdatePageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdatePageResponse {
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
        if !self.created_time.is_empty() {
            len += 1;
        }
        if !self.last_edited_time.is_empty() {
            len += 1;
        }
        if self.archived {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if self.public_url.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.notion.v1.UpdatePageResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.created_time.is_empty() {
            struct_ser.serialize_field("createdTime", &self.created_time)?;
        }
        if !self.last_edited_time.is_empty() {
            struct_ser.serialize_field("lastEditedTime", &self.last_edited_time)?;
        }
        if self.archived {
            struct_ser.serialize_field("archived", &self.archived)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if let Some(v) = self.public_url.as_ref() {
            struct_ser.serialize_field("publicUrl", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdatePageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "created_time",
            "createdTime",
            "last_edited_time",
            "lastEditedTime",
            "archived",
            "url",
            "public_url",
            "publicUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            CreatedTime,
            LastEditedTime,
            Archived,
            Url,
            PublicUrl,
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
                            "createdTime" | "created_time" => Ok(GeneratedField::CreatedTime),
                            "lastEditedTime" | "last_edited_time" => Ok(GeneratedField::LastEditedTime),
                            "archived" => Ok(GeneratedField::Archived),
                            "url" => Ok(GeneratedField::Url),
                            "publicUrl" | "public_url" => Ok(GeneratedField::PublicUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdatePageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.notion.v1.UpdatePageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdatePageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut created_time__ = None;
                let mut last_edited_time__ = None;
                let mut archived__ = None;
                let mut url__ = None;
                let mut public_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedTime => {
                            if created_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdTime"));
                            }
                            created_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LastEditedTime => {
                            if last_edited_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastEditedTime"));
                            }
                            last_edited_time__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Archived => {
                            if archived__.is_some() {
                                return Err(serde::de::Error::duplicate_field("archived"));
                            }
                            archived__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicUrl => {
                            if public_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicUrl"));
                            }
                            public_url__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdatePageResponse {
                    id: id__.unwrap_or_default(),
                    created_time: created_time__.unwrap_or_default(),
                    last_edited_time: last_edited_time__.unwrap_or_default(),
                    archived: archived__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                    public_url: public_url__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.notion.v1.UpdatePageResponse", FIELDS, GeneratedVisitor)
    }
}
