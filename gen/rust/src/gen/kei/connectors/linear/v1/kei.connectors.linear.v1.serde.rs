// @generated
impl serde::Serialize for Comment {
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
        if !self.body.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Comment", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.body.is_empty() {
            struct_ser.serialize_field("body", &self.body)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Comment {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "body",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Body,
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
                            "body" => Ok(GeneratedField::Body),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Comment;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Comment")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Comment, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut body__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Comment {
                    id: id__.unwrap_or_default(),
                    body: body__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Comment", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateCommentRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.issue_id.is_empty() {
            len += 1;
        }
        if !self.body.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.CreateCommentRequest", len)?;
        if !self.issue_id.is_empty() {
            struct_ser.serialize_field("issueId", &self.issue_id)?;
        }
        if !self.body.is_empty() {
            struct_ser.serialize_field("body", &self.body)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateCommentRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "issue_id",
            "issueId",
            "body",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            IssueId,
            Body,
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
                            "issueId" | "issue_id" => Ok(GeneratedField::IssueId),
                            "body" => Ok(GeneratedField::Body),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateCommentRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.CreateCommentRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateCommentRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut issue_id__ = None;
                let mut body__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::IssueId => {
                            if issue_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueId"));
                            }
                            issue_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateCommentRequest {
                    issue_id: issue_id__.unwrap_or_default(),
                    body: body__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.CreateCommentRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateCommentResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.success {
            len += 1;
        }
        if self.comment.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.CreateCommentResponse", len)?;
        if self.success {
            struct_ser.serialize_field("success", &self.success)?;
        }
        if let Some(v) = self.comment.as_ref() {
            struct_ser.serialize_field("comment", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateCommentResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "success",
            "comment",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Success,
            Comment,
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
                            "success" => Ok(GeneratedField::Success),
                            "comment" => Ok(GeneratedField::Comment),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateCommentResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.CreateCommentResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateCommentResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut success__ = None;
                let mut comment__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Success => {
                            if success__.is_some() {
                                return Err(serde::de::Error::duplicate_field("success"));
                            }
                            success__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Comment => {
                            if comment__.is_some() {
                                return Err(serde::de::Error::duplicate_field("comment"));
                            }
                            comment__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateCommentResponse {
                    success: success__.unwrap_or_default(),
                    comment: comment__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.CreateCommentResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateIssueRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.team_id.is_empty() {
            len += 1;
        }
        if self.description.is_some() {
            len += 1;
        }
        if self.priority.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.CreateIssueRequest", len)?;
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.team_id.is_empty() {
            struct_ser.serialize_field("teamId", &self.team_id)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        if let Some(v) = self.priority.as_ref() {
            struct_ser.serialize_field("priority", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateIssueRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "title",
            "team_id",
            "teamId",
            "description",
            "priority",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Title,
            TeamId,
            Description,
            Priority,
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
                            "title" => Ok(GeneratedField::Title),
                            "teamId" | "team_id" => Ok(GeneratedField::TeamId),
                            "description" => Ok(GeneratedField::Description),
                            "priority" => Ok(GeneratedField::Priority),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateIssueRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.CreateIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut title__ = None;
                let mut team_id__ = None;
                let mut description__ = None;
                let mut priority__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TeamId => {
                            if team_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("teamId"));
                            }
                            team_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = map_.next_value()?;
                        }
                        GeneratedField::Priority => {
                            if priority__.is_some() {
                                return Err(serde::de::Error::duplicate_field("priority"));
                            }
                            priority__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateIssueRequest {
                    title: title__.unwrap_or_default(),
                    team_id: team_id__.unwrap_or_default(),
                    description: description__,
                    priority: priority__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.CreateIssueRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateIssueResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.success {
            len += 1;
        }
        if self.issue.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.CreateIssueResponse", len)?;
        if self.success {
            struct_ser.serialize_field("success", &self.success)?;
        }
        if let Some(v) = self.issue.as_ref() {
            struct_ser.serialize_field("issue", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateIssueResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "success",
            "issue",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Success,
            Issue,
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
                            "success" => Ok(GeneratedField::Success),
                            "issue" => Ok(GeneratedField::Issue),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.CreateIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut success__ = None;
                let mut issue__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Success => {
                            if success__.is_some() {
                                return Err(serde::de::Error::duplicate_field("success"));
                            }
                            success__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Issue => {
                            if issue__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issue"));
                            }
                            issue__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateIssueResponse {
                    success: success__.unwrap_or_default(),
                    issue: issue__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.CreateIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Cycle {
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
        if self.number != 0 {
            len += 1;
        }
        if !self.starts_at.is_empty() {
            len += 1;
        }
        if !self.ends_at.is_empty() {
            len += 1;
        }
        if self.is_active {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Cycle", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if self.number != 0 {
            struct_ser.serialize_field("number", &self.number)?;
        }
        if !self.starts_at.is_empty() {
            struct_ser.serialize_field("startsAt", &self.starts_at)?;
        }
        if !self.ends_at.is_empty() {
            struct_ser.serialize_field("endsAt", &self.ends_at)?;
        }
        if self.is_active {
            struct_ser.serialize_field("isActive", &self.is_active)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Cycle {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "number",
            "starts_at",
            "startsAt",
            "ends_at",
            "endsAt",
            "is_active",
            "isActive",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            StartsAt,
            EndsAt,
            IsActive,
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
                            "number" => Ok(GeneratedField::Number),
                            "startsAt" | "starts_at" => Ok(GeneratedField::StartsAt),
                            "endsAt" | "ends_at" => Ok(GeneratedField::EndsAt),
                            "isActive" | "is_active" => Ok(GeneratedField::IsActive),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Cycle;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Cycle")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Cycle, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut starts_at__ = None;
                let mut ends_at__ = None;
                let mut is_active__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::StartsAt => {
                            if starts_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startsAt"));
                            }
                            starts_at__ = Some(map_.next_value()?);
                        }
                        GeneratedField::EndsAt => {
                            if ends_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("endsAt"));
                            }
                            ends_at__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IsActive => {
                            if is_active__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isActive"));
                            }
                            is_active__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Cycle {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    starts_at: starts_at__.unwrap_or_default(),
                    ends_at: ends_at__.unwrap_or_default(),
                    is_active: is_active__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Cycle", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetIssueRequest {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.GetIssueRequest", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
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
            "id",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
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
                formatter.write_str("struct kei.connectors.linear.v1.GetIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetIssueRequest {
                    id: id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.GetIssueRequest", FIELDS, GeneratedVisitor)
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
        if self.issue.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.GetIssueResponse", len)?;
        if let Some(v) = self.issue.as_ref() {
            struct_ser.serialize_field("issue", v)?;
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
            "issue",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Issue,
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
                            "issue" => Ok(GeneratedField::Issue),
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
                formatter.write_str("struct kei.connectors.linear.v1.GetIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut issue__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Issue => {
                            if issue__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issue"));
                            }
                            issue__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetIssueResponse {
                    issue: issue__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.GetIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTeamRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.team_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.GetTeamRequest", len)?;
        if !self.team_id.is_empty() {
            struct_ser.serialize_field("teamId", &self.team_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTeamRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "team_id",
            "teamId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TeamId,
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
                            "teamId" | "team_id" => Ok(GeneratedField::TeamId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTeamRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.GetTeamRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTeamRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut team_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TeamId => {
                            if team_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("teamId"));
                            }
                            team_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetTeamRequest {
                    team_id: team_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.GetTeamRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTeamResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.team.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.GetTeamResponse", len)?;
        if let Some(v) = self.team.as_ref() {
            struct_ser.serialize_field("team", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTeamResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "team",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Team,
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
                            "team" => Ok(GeneratedField::Team),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTeamResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.GetTeamResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTeamResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut team__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Team => {
                            if team__.is_some() {
                                return Err(serde::de::Error::duplicate_field("team"));
                            }
                            team__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetTeamResponse {
                    team: team__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.GetTeamResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Issue {
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
        if !self.identifier.is_empty() {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if self.description.is_some() {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if !self.created_at.is_empty() {
            len += 1;
        }
        if !self.updated_at.is_empty() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Issue", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.identifier.is_empty() {
            struct_ser.serialize_field("identifier", &self.identifier)?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if !self.created_at.is_empty() {
            struct_ser.serialize_field("createdAt", &self.created_at)?;
        }
        if !self.updated_at.is_empty() {
            struct_ser.serialize_field("updatedAt", &self.updated_at)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Issue {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "identifier",
            "title",
            "description",
            "url",
            "created_at",
            "createdAt",
            "updated_at",
            "updatedAt",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Identifier,
            Title,
            Description,
            Url,
            CreatedAt,
            UpdatedAt,
            Labels,
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
                            "identifier" => Ok(GeneratedField::Identifier),
                            "title" => Ok(GeneratedField::Title),
                            "description" => Ok(GeneratedField::Description),
                            "url" => Ok(GeneratedField::Url),
                            "createdAt" | "created_at" => Ok(GeneratedField::CreatedAt),
                            "updatedAt" | "updated_at" => Ok(GeneratedField::UpdatedAt),
                            "labels" => Ok(GeneratedField::Labels),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Issue;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Issue")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Issue, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut identifier__ = None;
                let mut title__ = None;
                let mut description__ = None;
                let mut url__ = None;
                let mut created_at__ = None;
                let mut updated_at__ = None;
                let mut labels__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = map_.next_value()?;
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedAt => {
                            if created_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdAt"));
                            }
                            created_at__ = Some(map_.next_value()?);
                        }
                        GeneratedField::UpdatedAt => {
                            if updated_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedAt"));
                            }
                            updated_at__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Labels => {
                            if labels__.is_some() {
                                return Err(serde::de::Error::duplicate_field("labels"));
                            }
                            labels__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Issue {
                    id: id__.unwrap_or_default(),
                    identifier: identifier__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    description: description__,
                    url: url__.unwrap_or_default(),
                    created_at: created_at__.unwrap_or_default(),
                    updated_at: updated_at__.unwrap_or_default(),
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Issue", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Label {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Label", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Label {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
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
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Label;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Label")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Label, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
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
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Label {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Label", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListCyclesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.team_id.is_empty() {
            len += 1;
        }
        if self.first.is_some() {
            len += 1;
        }
        if self.after.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.ListCyclesRequest", len)?;
        if !self.team_id.is_empty() {
            struct_ser.serialize_field("teamId", &self.team_id)?;
        }
        if let Some(v) = self.first.as_ref() {
            struct_ser.serialize_field("first", v)?;
        }
        if let Some(v) = self.after.as_ref() {
            struct_ser.serialize_field("after", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListCyclesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "team_id",
            "teamId",
            "first",
            "after",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TeamId,
            First,
            After,
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
                            "teamId" | "team_id" => Ok(GeneratedField::TeamId),
                            "first" => Ok(GeneratedField::First),
                            "after" => Ok(GeneratedField::After),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListCyclesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.ListCyclesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListCyclesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut team_id__ = None;
                let mut first__ = None;
                let mut after__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TeamId => {
                            if team_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("teamId"));
                            }
                            team_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::First => {
                            if first__.is_some() {
                                return Err(serde::de::Error::duplicate_field("first"));
                            }
                            first__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::After => {
                            if after__.is_some() {
                                return Err(serde::de::Error::duplicate_field("after"));
                            }
                            after__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListCyclesRequest {
                    team_id: team_id__.unwrap_or_default(),
                    first: first__,
                    after: after__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.ListCyclesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListCyclesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.nodes.is_empty() {
            len += 1;
        }
        if self.page_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.ListCyclesResponse", len)?;
        if !self.nodes.is_empty() {
            struct_ser.serialize_field("nodes", &self.nodes)?;
        }
        if let Some(v) = self.page_info.as_ref() {
            struct_ser.serialize_field("pageInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListCyclesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "nodes",
            "page_info",
            "pageInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Nodes,
            PageInfo,
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
                            "nodes" => Ok(GeneratedField::Nodes),
                            "pageInfo" | "page_info" => Ok(GeneratedField::PageInfo),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListCyclesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.ListCyclesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListCyclesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut nodes__ = None;
                let mut page_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Nodes => {
                            if nodes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodes"));
                            }
                            nodes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PageInfo => {
                            if page_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageInfo"));
                            }
                            page_info__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListCyclesResponse {
                    nodes: nodes__.unwrap_or_default(),
                    page_info: page_info__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.ListCyclesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListProjectsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.first.is_some() {
            len += 1;
        }
        if self.after.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.ListProjectsRequest", len)?;
        if let Some(v) = self.first.as_ref() {
            struct_ser.serialize_field("first", v)?;
        }
        if let Some(v) = self.after.as_ref() {
            struct_ser.serialize_field("after", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListProjectsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "first",
            "after",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            First,
            After,
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
                            "first" => Ok(GeneratedField::First),
                            "after" => Ok(GeneratedField::After),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListProjectsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.ListProjectsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListProjectsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut first__ = None;
                let mut after__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::First => {
                            if first__.is_some() {
                                return Err(serde::de::Error::duplicate_field("first"));
                            }
                            first__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::After => {
                            if after__.is_some() {
                                return Err(serde::de::Error::duplicate_field("after"));
                            }
                            after__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListProjectsRequest {
                    first: first__,
                    after: after__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.ListProjectsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListProjectsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.nodes.is_empty() {
            len += 1;
        }
        if self.page_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.ListProjectsResponse", len)?;
        if !self.nodes.is_empty() {
            struct_ser.serialize_field("nodes", &self.nodes)?;
        }
        if let Some(v) = self.page_info.as_ref() {
            struct_ser.serialize_field("pageInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListProjectsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "nodes",
            "page_info",
            "pageInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Nodes,
            PageInfo,
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
                            "nodes" => Ok(GeneratedField::Nodes),
                            "pageInfo" | "page_info" => Ok(GeneratedField::PageInfo),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListProjectsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.ListProjectsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListProjectsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut nodes__ = None;
                let mut page_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Nodes => {
                            if nodes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodes"));
                            }
                            nodes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PageInfo => {
                            if page_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageInfo"));
                            }
                            page_info__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListProjectsResponse {
                    nodes: nodes__.unwrap_or_default(),
                    page_info: page_info__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.ListProjectsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PageInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.has_next_page {
            len += 1;
        }
        if self.end_cursor.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.PageInfo", len)?;
        if self.has_next_page {
            struct_ser.serialize_field("hasNextPage", &self.has_next_page)?;
        }
        if let Some(v) = self.end_cursor.as_ref() {
            struct_ser.serialize_field("endCursor", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PageInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "has_next_page",
            "hasNextPage",
            "end_cursor",
            "endCursor",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            HasNextPage,
            EndCursor,
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
                            "hasNextPage" | "has_next_page" => Ok(GeneratedField::HasNextPage),
                            "endCursor" | "end_cursor" => Ok(GeneratedField::EndCursor),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PageInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.PageInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PageInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut has_next_page__ = None;
                let mut end_cursor__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::HasNextPage => {
                            if has_next_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hasNextPage"));
                            }
                            has_next_page__ = Some(map_.next_value()?);
                        }
                        GeneratedField::EndCursor => {
                            if end_cursor__.is_some() {
                                return Err(serde::de::Error::duplicate_field("endCursor"));
                            }
                            end_cursor__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(PageInfo {
                    has_next_page: has_next_page__.unwrap_or_default(),
                    end_cursor: end_cursor__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.PageInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Project {
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
        if self.description.is_some() {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Project", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Project {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
            "description",
            "url",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
            Description,
            Url,
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
                            "description" => Ok(GeneratedField::Description),
                            "url" => Ok(GeneratedField::Url),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Project;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Project")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Project, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
                let mut description__ = None;
                let mut url__ = None;
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
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = map_.next_value()?;
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Project {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    description: description__,
                    url: url__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Project", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SearchIssuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.term.is_empty() {
            len += 1;
        }
        if self.first.is_some() {
            len += 1;
        }
        if self.after.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.SearchIssuesRequest", len)?;
        if !self.term.is_empty() {
            struct_ser.serialize_field("term", &self.term)?;
        }
        if let Some(v) = self.first.as_ref() {
            struct_ser.serialize_field("first", v)?;
        }
        if let Some(v) = self.after.as_ref() {
            struct_ser.serialize_field("after", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SearchIssuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "term",
            "first",
            "after",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Term,
            First,
            After,
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
                            "term" => Ok(GeneratedField::Term),
                            "first" => Ok(GeneratedField::First),
                            "after" => Ok(GeneratedField::After),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SearchIssuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.SearchIssuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SearchIssuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut term__ = None;
                let mut first__ = None;
                let mut after__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Term => {
                            if term__.is_some() {
                                return Err(serde::de::Error::duplicate_field("term"));
                            }
                            term__ = Some(map_.next_value()?);
                        }
                        GeneratedField::First => {
                            if first__.is_some() {
                                return Err(serde::de::Error::duplicate_field("first"));
                            }
                            first__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::After => {
                            if after__.is_some() {
                                return Err(serde::de::Error::duplicate_field("after"));
                            }
                            after__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SearchIssuesRequest {
                    term: term__.unwrap_or_default(),
                    first: first__,
                    after: after__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.SearchIssuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SearchIssuesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.nodes.is_empty() {
            len += 1;
        }
        if self.page_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.SearchIssuesResponse", len)?;
        if !self.nodes.is_empty() {
            struct_ser.serialize_field("nodes", &self.nodes)?;
        }
        if let Some(v) = self.page_info.as_ref() {
            struct_ser.serialize_field("pageInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SearchIssuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "nodes",
            "page_info",
            "pageInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Nodes,
            PageInfo,
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
                            "nodes" => Ok(GeneratedField::Nodes),
                            "pageInfo" | "page_info" => Ok(GeneratedField::PageInfo),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SearchIssuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.SearchIssuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SearchIssuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut nodes__ = None;
                let mut page_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Nodes => {
                            if nodes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodes"));
                            }
                            nodes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PageInfo => {
                            if page_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageInfo"));
                            }
                            page_info__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SearchIssuesResponse {
                    nodes: nodes__.unwrap_or_default(),
                    page_info: page_info__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.SearchIssuesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Team {
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
        if !self.name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.Team", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.key.is_empty() {
            struct_ser.serialize_field("key", &self.key)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Team {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "key",
            "name",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Key,
            Name,
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
                            "name" => Ok(GeneratedField::Name),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Team;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.Team")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Team, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut key__ = None;
                let mut name__ = None;
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
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Team {
                    id: id__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.Team", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateIssueRequest {
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
        if self.title.is_some() {
            len += 1;
        }
        if self.description.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.UpdateIssueRequest", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if let Some(v) = self.title.as_ref() {
            struct_ser.serialize_field("title", v)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateIssueRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "title",
            "description",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Title,
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
                            "id" => Ok(GeneratedField::Id),
                            "title" => Ok(GeneratedField::Title),
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
            type Value = UpdateIssueRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.UpdateIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut title__ = None;
                let mut description__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = map_.next_value()?;
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
                Ok(UpdateIssueRequest {
                    id: id__.unwrap_or_default(),
                    title: title__,
                    description: description__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.UpdateIssueRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateIssueResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.success {
            len += 1;
        }
        if self.issue.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.linear.v1.UpdateIssueResponse", len)?;
        if self.success {
            struct_ser.serialize_field("success", &self.success)?;
        }
        if let Some(v) = self.issue.as_ref() {
            struct_ser.serialize_field("issue", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateIssueResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "success",
            "issue",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Success,
            Issue,
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
                            "success" => Ok(GeneratedField::Success),
                            "issue" => Ok(GeneratedField::Issue),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdateIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.linear.v1.UpdateIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut success__ = None;
                let mut issue__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Success => {
                            if success__.is_some() {
                                return Err(serde::de::Error::duplicate_field("success"));
                            }
                            success__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Issue => {
                            if issue__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issue"));
                            }
                            issue__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateIssueResponse {
                    success: success__.unwrap_or_default(),
                    issue: issue__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.linear.v1.UpdateIssueResponse", FIELDS, GeneratedVisitor)
    }
}
