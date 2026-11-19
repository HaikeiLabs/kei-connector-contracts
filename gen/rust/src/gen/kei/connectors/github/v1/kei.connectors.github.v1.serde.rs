// @generated
impl serde::Serialize for CheckRunSummary {
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
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.head_sha.is_empty() {
            len += 1;
        }
        if !self.status.is_empty() {
            len += 1;
        }
        if self.conclusion.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CheckRunSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.head_sha.is_empty() {
            struct_ser.serialize_field("headSha", &self.head_sha)?;
        }
        if !self.status.is_empty() {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if let Some(v) = self.conclusion.as_ref() {
            struct_ser.serialize_field("conclusion", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CheckRunSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
            "head_sha",
            "headSha",
            "status",
            "conclusion",
            "html_url",
            "htmlUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
            HeadSha,
            Status,
            Conclusion,
            HtmlUrl,
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
                            "headSha" | "head_sha" => Ok(GeneratedField::HeadSha),
                            "status" => Ok(GeneratedField::Status),
                            "conclusion" => Ok(GeneratedField::Conclusion),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CheckRunSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CheckRunSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CheckRunSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
                let mut head_sha__ = None;
                let mut status__ = None;
                let mut conclusion__ = None;
                let mut html_url__ = None;
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
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HeadSha => {
                            if head_sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("headSha"));
                            }
                            head_sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Conclusion => {
                            if conclusion__.is_some() {
                                return Err(serde::de::Error::duplicate_field("conclusion"));
                            }
                            conclusion__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CheckRunSummary {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    head_sha: head_sha__.unwrap_or_default(),
                    status: status__.unwrap_or_default(),
                    conclusion: conclusion__,
                    html_url: html_url__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CheckRunSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CommitSummary {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CommitSummary", len)?;
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CommitSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sha",
            "html_url",
            "htmlUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sha,
            HtmlUrl,
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
                            "sha" => Ok(GeneratedField::Sha),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CommitSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CommitSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CommitSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sha__ = None;
                let mut html_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CommitSummary {
                    sha: sha__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CommitSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ContentFile {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.path.is_empty() {
            len += 1;
        }
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.download_url.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ContentFile", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.path.is_empty() {
            struct_ser.serialize_field("path", &self.path)?;
        }
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.download_url.as_ref() {
            struct_ser.serialize_field("downloadUrl", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ContentFile {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "path",
            "sha",
            "html_url",
            "htmlUrl",
            "download_url",
            "downloadUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Path,
            Sha,
            HtmlUrl,
            DownloadUrl,
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
                            "name" => Ok(GeneratedField::Name),
                            "path" => Ok(GeneratedField::Path),
                            "sha" => Ok(GeneratedField::Sha),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "downloadUrl" | "download_url" => Ok(GeneratedField::DownloadUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ContentFile;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ContentFile")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ContentFile, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut path__ = None;
                let mut sha__ = None;
                let mut html_url__ = None;
                let mut download_url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Path => {
                            if path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("path"));
                            }
                            path__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::DownloadUrl => {
                            if download_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("downloadUrl"));
                            }
                            download_url__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ContentFile {
                    name: name__.unwrap_or_default(),
                    path: path__.unwrap_or_default(),
                    sha: sha__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                    download_url: download_url__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ContentFile", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateBlobRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.content.is_empty() {
            len += 1;
        }
        if !self.encoding.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateBlobRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.content.is_empty() {
            struct_ser.serialize_field("content", &self.content)?;
        }
        if !self.encoding.is_empty() {
            struct_ser.serialize_field("encoding", &self.encoding)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateBlobRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "content",
            "encoding",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Content,
            Encoding,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "content" => Ok(GeneratedField::Content),
                            "encoding" => Ok(GeneratedField::Encoding),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateBlobRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateBlobRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateBlobRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut content__ = None;
                let mut encoding__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Content => {
                            if content__.is_some() {
                                return Err(serde::de::Error::duplicate_field("content"));
                            }
                            content__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Encoding => {
                            if encoding__.is_some() {
                                return Err(serde::de::Error::duplicate_field("encoding"));
                            }
                            encoding__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateBlobRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    content: content__.unwrap_or_default(),
                    encoding: encoding__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateBlobRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateBlobResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateBlobResponse", len)?;
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateBlobResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sha",
            "url",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sha,
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
                            "sha" => Ok(GeneratedField::Sha),
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
            type Value = CreateBlobResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateBlobResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateBlobResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sha__ = None;
                let mut url__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
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
                Ok(CreateBlobResponse {
                    sha: sha__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateBlobResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateCommitRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        if !self.tree.is_empty() {
            len += 1;
        }
        if !self.parents.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateCommitRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        if !self.tree.is_empty() {
            struct_ser.serialize_field("tree", &self.tree)?;
        }
        if !self.parents.is_empty() {
            struct_ser.serialize_field("parents", &self.parents)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateCommitRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "message",
            "tree",
            "parents",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Message,
            Tree,
            Parents,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "message" => Ok(GeneratedField::Message),
                            "tree" => Ok(GeneratedField::Tree),
                            "parents" => Ok(GeneratedField::Parents),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateCommitRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateCommitRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateCommitRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut message__ = None;
                let mut tree__ = None;
                let mut parents__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Tree => {
                            if tree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tree"));
                            }
                            tree__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Parents => {
                            if parents__.is_some() {
                                return Err(serde::de::Error::duplicate_field("parents"));
                            }
                            parents__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateCommitRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                    tree: tree__.unwrap_or_default(),
                    parents: parents__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateCommitRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateCommitResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        if self.tree.is_some() {
            len += 1;
        }
        if !self.parents.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateCommitResponse", len)?;
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        if let Some(v) = self.tree.as_ref() {
            struct_ser.serialize_field("tree", v)?;
        }
        if !self.parents.is_empty() {
            struct_ser.serialize_field("parents", &self.parents)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateCommitResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sha",
            "html_url",
            "htmlUrl",
            "message",
            "tree",
            "parents",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sha,
            HtmlUrl,
            Message,
            Tree,
            Parents,
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
                            "sha" => Ok(GeneratedField::Sha),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "message" => Ok(GeneratedField::Message),
                            "tree" => Ok(GeneratedField::Tree),
                            "parents" => Ok(GeneratedField::Parents),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateCommitResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateCommitResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateCommitResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sha__ = None;
                let mut html_url__ = None;
                let mut message__ = None;
                let mut tree__ = None;
                let mut parents__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Tree => {
                            if tree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tree"));
                            }
                            tree__ = map_.next_value()?;
                        }
                        GeneratedField::Parents => {
                            if parents__.is_some() {
                                return Err(serde::de::Error::duplicate_field("parents"));
                            }
                            parents__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateCommitResponse {
                    sha: sha__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                    tree: tree__,
                    parents: parents__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateCommitResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateIssueCommentRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.issue_number != 0 {
            len += 1;
        }
        if !self.body.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateIssueCommentRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.issue_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("issueNumber", ToString::to_string(&self.issue_number).as_str())?;
        }
        if !self.body.is_empty() {
            struct_ser.serialize_field("body", &self.body)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateIssueCommentRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "issue_number",
            "issueNumber",
            "body",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            IssueNumber,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "issueNumber" | "issue_number" => Ok(GeneratedField::IssueNumber),
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
            type Value = CreateIssueCommentRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateIssueCommentRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueCommentRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut issue_number__ = None;
                let mut body__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IssueNumber => {
                            if issue_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueNumber"));
                            }
                            issue_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
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
                Ok(CreateIssueCommentRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    issue_number: issue_number__.unwrap_or_default(),
                    body: body__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateIssueCommentRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateIssueCommentResponse {
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
        if !self.body.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateIssueCommentResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.body.is_empty() {
            struct_ser.serialize_field("body", &self.body)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateIssueCommentResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "body",
            "html_url",
            "htmlUrl",
            "author",
            "user",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Body,
            HtmlUrl,
            Author,
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
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "user" | "author" => Ok(GeneratedField::Author),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateIssueCommentResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateIssueCommentResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueCommentResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut author__ = None;
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
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateIssueCommentResponse {
                    id: id__.unwrap_or_default(),
                    body: body__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                    author: author__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateIssueCommentResponse", FIELDS, GeneratedVisitor)
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
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateIssueRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
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
            "owner",
            "repository",
            "title",
            "body",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Title,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "title" => Ok(GeneratedField::Title),
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
            type Value = CreateIssueRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut title__ = None;
                let mut body__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateIssueRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    body: body__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateIssueRequest", FIELDS, GeneratedVisitor)
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
        if self.id != 0 {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateIssueResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
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
            "id",
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = CreateIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(CreateIssueResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreatePullRequestRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.head.is_empty() {
            len += 1;
        }
        if !self.base.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if self.draft.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreatePullRequestRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.head.is_empty() {
            struct_ser.serialize_field("head", &self.head)?;
        }
        if !self.base.is_empty() {
            struct_ser.serialize_field("base", &self.base)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if let Some(v) = self.draft.as_ref() {
            struct_ser.serialize_field("draft", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePullRequestRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "title",
            "head",
            "base",
            "body",
            "draft",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Title,
            Head,
            Base,
            Body,
            Draft,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "title" => Ok(GeneratedField::Title),
                            "head" => Ok(GeneratedField::Head),
                            "base" => Ok(GeneratedField::Base),
                            "body" => Ok(GeneratedField::Body),
                            "draft" => Ok(GeneratedField::Draft),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreatePullRequestRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreatePullRequestRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePullRequestRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut title__ = None;
                let mut head__ = None;
                let mut base__ = None;
                let mut body__ = None;
                let mut draft__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Head => {
                            if head__.is_some() {
                                return Err(serde::de::Error::duplicate_field("head"));
                            }
                            head__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::Draft => {
                            if draft__.is_some() {
                                return Err(serde::de::Error::duplicate_field("draft"));
                            }
                            draft__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreatePullRequestRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    head: head__.unwrap_or_default(),
                    base: base__.unwrap_or_default(),
                    body: body__,
                    draft: draft__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreatePullRequestRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreatePullRequestResponse {
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
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.draft {
            len += 1;
        }
        if self.head.is_some() {
            len += 1;
        }
        if self.base.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreatePullRequestResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if self.draft {
            struct_ser.serialize_field("draft", &self.draft)?;
        }
        if let Some(v) = self.head.as_ref() {
            struct_ser.serialize_field("head", v)?;
        }
        if let Some(v) = self.base.as_ref() {
            struct_ser.serialize_field("base", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePullRequestResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "number",
            "title",
            "state",
            "html_url",
            "htmlUrl",
            "draft",
            "head",
            "base",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            HtmlUrl,
            Draft,
            Head,
            Base,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "draft" => Ok(GeneratedField::Draft),
                            "head" => Ok(GeneratedField::Head),
                            "base" => Ok(GeneratedField::Base),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = CreatePullRequestResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreatePullRequestResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePullRequestResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut html_url__ = None;
                let mut draft__ = None;
                let mut head__ = None;
                let mut base__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Draft => {
                            if draft__.is_some() {
                                return Err(serde::de::Error::duplicate_field("draft"));
                            }
                            draft__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Head => {
                            if head__.is_some() {
                                return Err(serde::de::Error::duplicate_field("head"));
                            }
                            head__ = map_.next_value()?;
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(CreatePullRequestResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                    draft: draft__.unwrap_or_default(),
                    head: head__,
                    base: base__,
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreatePullRequestResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreatePullReviewRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.pull_number != 0 {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if self.event.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreatePullReviewRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.pull_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("pullNumber", ToString::to_string(&self.pull_number).as_str())?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if let Some(v) = self.event.as_ref() {
            struct_ser.serialize_field("event", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePullReviewRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "pull_number",
            "pullNumber",
            "body",
            "event",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PullNumber,
            Body,
            Event,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "pullNumber" | "pull_number" => Ok(GeneratedField::PullNumber),
                            "body" => Ok(GeneratedField::Body),
                            "event" => Ok(GeneratedField::Event),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreatePullReviewRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreatePullReviewRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePullReviewRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut pull_number__ = None;
                let mut body__ = None;
                let mut event__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PullNumber => {
                            if pull_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pullNumber"));
                            }
                            pull_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::Event => {
                            if event__.is_some() {
                                return Err(serde::de::Error::duplicate_field("event"));
                            }
                            event__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreatePullReviewRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    pull_number: pull_number__.unwrap_or_default(),
                    body: body__,
                    event: event__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreatePullReviewRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreatePullReviewResponse {
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
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.submitted_at.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreatePullReviewResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.submitted_at.as_ref() {
            struct_ser.serialize_field("submittedAt", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreatePullReviewResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "submitted_at",
            "submittedAt",
            "author",
            "user",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            State,
            Body,
            HtmlUrl,
            SubmittedAt,
            Author,
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
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "submittedAt" | "submitted_at" => Ok(GeneratedField::SubmittedAt),
                            "user" | "author" => Ok(GeneratedField::Author),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreatePullReviewResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreatePullReviewResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreatePullReviewResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut submitted_at__ = None;
                let mut author__ = None;
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
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SubmittedAt => {
                            if submitted_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("submittedAt"));
                            }
                            submitted_at__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreatePullReviewResponse {
                    id: id__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    submitted_at: submitted_at__,
                    author: author__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreatePullReviewResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateRefRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if !self.sha.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateRefRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateRefRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "ref",
            "sha",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Ref,
            Sha,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "ref" => Ok(GeneratedField::Ref),
                            "sha" => Ok(GeneratedField::Sha),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateRefRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateRefRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateRefRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut r#ref__ = None;
                let mut sha__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateRefRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    r#ref: r#ref__.unwrap_or_default(),
                    sha: sha__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateRefRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateRefResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if self.object.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateRefResponse", len)?;
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateRefResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ref",
            "object",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ref,
            Object,
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
                            "ref" => Ok(GeneratedField::Ref),
                            "object" => Ok(GeneratedField::Object),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateRefResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateRefResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateRefResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#ref__ = None;
                let mut object__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateRefResponse {
                    r#ref: r#ref__.unwrap_or_default(),
                    object: object__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateRefResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateTreeRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.tree.is_empty() {
            len += 1;
        }
        if self.base_tree.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateTreeRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.tree.is_empty() {
            struct_ser.serialize_field("tree", &self.tree)?;
        }
        if let Some(v) = self.base_tree.as_ref() {
            struct_ser.serialize_field("baseTree", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateTreeRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "tree",
            "base_tree",
            "baseTree",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Tree,
            BaseTree,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "tree" => Ok(GeneratedField::Tree),
                            "baseTree" | "base_tree" => Ok(GeneratedField::BaseTree),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateTreeRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateTreeRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateTreeRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut tree__ = None;
                let mut base_tree__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Tree => {
                            if tree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tree"));
                            }
                            tree__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BaseTree => {
                            if base_tree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseTree"));
                            }
                            base_tree__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateTreeRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    tree: tree__.unwrap_or_default(),
                    base_tree: base_tree__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateTreeRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateTreeResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.url.is_empty() {
            len += 1;
        }
        if !self.tree.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.CreateTreeResponse", len)?;
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.url.is_empty() {
            struct_ser.serialize_field("url", &self.url)?;
        }
        if !self.tree.is_empty() {
            struct_ser.serialize_field("tree", &self.tree)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateTreeResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sha",
            "url",
            "tree",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sha,
            Url,
            Tree,
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
                            "sha" => Ok(GeneratedField::Sha),
                            "url" => Ok(GeneratedField::Url),
                            "tree" => Ok(GeneratedField::Tree),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateTreeResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.CreateTreeResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateTreeResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sha__ = None;
                let mut url__ = None;
                let mut tree__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Url => {
                            if url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("url"));
                            }
                            url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Tree => {
                            if tree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tree"));
                            }
                            tree__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(CreateTreeResponse {
                    sha: sha__.unwrap_or_default(),
                    url: url__.unwrap_or_default(),
                    tree: tree__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.CreateTreeResponse", FIELDS, GeneratedVisitor)
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
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.issue_number != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetIssueRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.issue_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("issueNumber", ToString::to_string(&self.issue_number).as_str())?;
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
            "owner",
            "repository",
            "issue_number",
            "issueNumber",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            IssueNumber,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "issueNumber" | "issue_number" => Ok(GeneratedField::IssueNumber),
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
                formatter.write_str("struct kei.connectors.github.v1.GetIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut issue_number__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IssueNumber => {
                            if issue_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueNumber"));
                            }
                            issue_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetIssueRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    issue_number: issue_number__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetIssueRequest", FIELDS, GeneratedVisitor)
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
        if self.id != 0 {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetIssueResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
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
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = GetIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(GetIssueResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPullRequestRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.pull_number != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetPullRequestRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.pull_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("pullNumber", ToString::to_string(&self.pull_number).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPullRequestRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "pull_number",
            "pullNumber",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PullNumber,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "pullNumber" | "pull_number" => Ok(GeneratedField::PullNumber),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetPullRequestRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetPullRequestRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPullRequestRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut pull_number__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PullNumber => {
                            if pull_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pullNumber"));
                            }
                            pull_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetPullRequestRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    pull_number: pull_number__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetPullRequestRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPullRequestResponse {
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
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.draft {
            len += 1;
        }
        if self.merged_at.is_some() {
            len += 1;
        }
        if self.head.is_some() {
            len += 1;
        }
        if self.base.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetPullRequestResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if self.draft {
            struct_ser.serialize_field("draft", &self.draft)?;
        }
        if let Some(v) = self.merged_at.as_ref() {
            struct_ser.serialize_field("mergedAt", v)?;
        }
        if let Some(v) = self.head.as_ref() {
            struct_ser.serialize_field("head", v)?;
        }
        if let Some(v) = self.base.as_ref() {
            struct_ser.serialize_field("base", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPullRequestResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "draft",
            "merged_at",
            "mergedAt",
            "head",
            "base",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Draft,
            MergedAt,
            Head,
            Base,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "draft" => Ok(GeneratedField::Draft),
                            "mergedAt" | "merged_at" => Ok(GeneratedField::MergedAt),
                            "head" => Ok(GeneratedField::Head),
                            "base" => Ok(GeneratedField::Base),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = GetPullRequestResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetPullRequestResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPullRequestResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut draft__ = None;
                let mut merged_at__ = None;
                let mut head__ = None;
                let mut base__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Draft => {
                            if draft__.is_some() {
                                return Err(serde::de::Error::duplicate_field("draft"));
                            }
                            draft__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MergedAt => {
                            if merged_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mergedAt"));
                            }
                            merged_at__ = map_.next_value()?;
                        }
                        GeneratedField::Head => {
                            if head__.is_some() {
                                return Err(serde::de::Error::duplicate_field("head"));
                            }
                            head__ = map_.next_value()?;
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(GetPullRequestResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    draft: draft__.unwrap_or_default(),
                    merged_at: merged_at__,
                    head: head__,
                    base: base__,
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetPullRequestResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRefRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.r#ref.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetRefRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRefRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "ref",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Ref,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "ref" => Ok(GeneratedField::Ref),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRefRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetRefRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRefRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut r#ref__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRefRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    r#ref: r#ref__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetRefRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRefResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if self.object.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetRefResponse", len)?;
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRefResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ref",
            "object",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ref,
            Object,
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
                            "ref" => Ok(GeneratedField::Ref),
                            "object" => Ok(GeneratedField::Object),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRefResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetRefResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRefResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#ref__ = None;
                let mut object__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRefResponse {
                    r#ref: r#ref__.unwrap_or_default(),
                    object: object__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetRefResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRepositoryRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetRepositoryRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRepositoryRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRepositoryRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetRepositoryRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRepositoryRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRepositoryRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetRepositoryRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRepositoryResponse {
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
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.full_name.is_empty() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.private {
            len += 1;
        }
        if self.description.is_some() {
            len += 1;
        }
        if !self.default_branch.is_empty() {
            len += 1;
        }
        if self.owner.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GetRepositoryResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.full_name.is_empty() {
            struct_ser.serialize_field("fullName", &self.full_name)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if self.private {
            struct_ser.serialize_field("private", &self.private)?;
        }
        if let Some(v) = self.description.as_ref() {
            struct_ser.serialize_field("description", v)?;
        }
        if !self.default_branch.is_empty() {
            struct_ser.serialize_field("defaultBranch", &self.default_branch)?;
        }
        if let Some(v) = self.owner.as_ref() {
            struct_ser.serialize_field("owner", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRepositoryResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
            "full_name",
            "fullName",
            "html_url",
            "htmlUrl",
            "private",
            "description",
            "default_branch",
            "defaultBranch",
            "owner",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
            FullName,
            HtmlUrl,
            Private,
            Description,
            DefaultBranch,
            Owner,
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
                            "fullName" | "full_name" => Ok(GeneratedField::FullName),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "private" => Ok(GeneratedField::Private),
                            "description" => Ok(GeneratedField::Description),
                            "defaultBranch" | "default_branch" => Ok(GeneratedField::DefaultBranch),
                            "owner" => Ok(GeneratedField::Owner),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRepositoryResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GetRepositoryResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRepositoryResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
                let mut full_name__ = None;
                let mut html_url__ = None;
                let mut private__ = None;
                let mut description__ = None;
                let mut default_branch__ = None;
                let mut owner__ = None;
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
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::FullName => {
                            if full_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fullName"));
                            }
                            full_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Private => {
                            if private__.is_some() {
                                return Err(serde::de::Error::duplicate_field("private"));
                            }
                            private__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = map_.next_value()?;
                        }
                        GeneratedField::DefaultBranch => {
                            if default_branch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("defaultBranch"));
                            }
                            default_branch__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetRepositoryResponse {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    full_name: full_name__.unwrap_or_default(),
                    html_url: html_url__.unwrap_or_default(),
                    private: private__.unwrap_or_default(),
                    description: description__,
                    default_branch: default_branch__.unwrap_or_default(),
                    owner: owner__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GetRepositoryResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GitObject {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sha.is_empty() {
            len += 1;
        }
        if !self.r#type.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.GitObject", len)?;
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if !self.r#type.is_empty() {
            struct_ser.serialize_field("type", &self.r#type)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GitObject {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sha",
            "type",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sha,
            Type,
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
                            "sha" => Ok(GeneratedField::Sha),
                            "type" => Ok(GeneratedField::Type),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GitObject;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.GitObject")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GitObject, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sha__ = None;
                let mut r#type__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GitObject {
                    sha: sha__.unwrap_or_default(),
                    r#type: r#type__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.GitObject", FIELDS, GeneratedVisitor)
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
        if self.id != 0 {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.IssueSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
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
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = IssueSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.IssueSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<IssueSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(IssueSummary {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.IssueSummary", FIELDS, GeneratedVisitor)
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
        if self.id != 0 {
            len += 1;
        }
        if !self.name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.Label", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
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
                formatter.write_str("struct kei.connectors.github.v1.Label")
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
                            id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
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
        deserializer.deserialize_struct("kei.connectors.github.v1.Label", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListCheckRunsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        if self.page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListCheckRunsRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if let Some(v) = self.per_page.as_ref() {
            struct_ser.serialize_field("perPage", v)?;
        }
        if let Some(v) = self.page.as_ref() {
            struct_ser.serialize_field("page", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListCheckRunsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "ref",
            "per_page",
            "perPage",
            "page",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Ref,
            PerPage,
            Page,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "ref" => Ok(GeneratedField::Ref),
                            "perPage" | "per_page" => Ok(GeneratedField::PerPage),
                            "page" => Ok(GeneratedField::Page),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListCheckRunsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListCheckRunsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListCheckRunsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut r#ref__ = None;
                let mut per_page__ = None;
                let mut page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PerPage => {
                            if per_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("perPage"));
                            }
                            per_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Page => {
                            if page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("page"));
                            }
                            page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListCheckRunsRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    r#ref: r#ref__.unwrap_or_default(),
                    per_page: per_page__,
                    page: page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListCheckRunsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListCheckRunsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.check_runs.is_empty() {
            len += 1;
        }
        if self.total_count.is_some() {
            len += 1;
        }
        if self.next_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListCheckRunsResponse", len)?;
        if !self.check_runs.is_empty() {
            struct_ser.serialize_field("checkRuns", &self.check_runs)?;
        }
        if let Some(v) = self.total_count.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("totalCount", ToString::to_string(&v).as_str())?;
        }
        if let Some(v) = self.next_page.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nextPage", ToString::to_string(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListCheckRunsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "check_runs",
            "checkRuns",
            "total_count",
            "totalCount",
            "next_page",
            "nextPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            CheckRuns,
            TotalCount,
            NextPage,
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
                            "checkRuns" | "check_runs" => Ok(GeneratedField::CheckRuns),
                            "totalCount" | "total_count" => Ok(GeneratedField::TotalCount),
                            "nextPage" | "next_page" => Ok(GeneratedField::NextPage),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListCheckRunsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListCheckRunsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListCheckRunsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut check_runs__ = None;
                let mut total_count__ = None;
                let mut next_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::CheckRuns => {
                            if check_runs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("checkRuns"));
                            }
                            check_runs__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TotalCount => {
                            if total_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalCount"));
                            }
                            total_count__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::NextPage => {
                            if next_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPage"));
                            }
                            next_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListCheckRunsResponse {
                    check_runs: check_runs__.unwrap_or_default(),
                    total_count: total_count__,
                    next_page: next_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListCheckRunsResponse", FIELDS, GeneratedVisitor)
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
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.state.is_some() {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListIssuesRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if let Some(v) = self.state.as_ref() {
            struct_ser.serialize_field("state", v)?;
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
            "owner",
            "repository",
            "state",
            "per_page",
            "perPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            State,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "state" => Ok(GeneratedField::State),
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
                formatter.write_str("struct kei.connectors.github.v1.ListIssuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListIssuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut state__ = None;
                let mut per_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = map_.next_value()?;
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
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    state: state__,
                    per_page: per_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListIssuesRequest", FIELDS, GeneratedVisitor)
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListIssuesResponse", len)?;
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
                formatter.write_str("struct kei.connectors.github.v1.ListIssuesResponse")
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
        deserializer.deserialize_struct("kei.connectors.github.v1.ListIssuesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListPullReviewsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.pull_number != 0 {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        if self.page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListPullReviewsRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.pull_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("pullNumber", ToString::to_string(&self.pull_number).as_str())?;
        }
        if let Some(v) = self.per_page.as_ref() {
            struct_ser.serialize_field("perPage", v)?;
        }
        if let Some(v) = self.page.as_ref() {
            struct_ser.serialize_field("page", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListPullReviewsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "pull_number",
            "pullNumber",
            "per_page",
            "perPage",
            "page",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PullNumber,
            PerPage,
            Page,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "pullNumber" | "pull_number" => Ok(GeneratedField::PullNumber),
                            "perPage" | "per_page" => Ok(GeneratedField::PerPage),
                            "page" => Ok(GeneratedField::Page),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListPullReviewsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListPullReviewsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListPullReviewsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut pull_number__ = None;
                let mut per_page__ = None;
                let mut page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PullNumber => {
                            if pull_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pullNumber"));
                            }
                            pull_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PerPage => {
                            if per_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("perPage"));
                            }
                            per_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Page => {
                            if page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("page"));
                            }
                            page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListPullReviewsRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    pull_number: pull_number__.unwrap_or_default(),
                    per_page: per_page__,
                    page: page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListPullReviewsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListPullReviewsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.reviews.is_empty() {
            len += 1;
        }
        if self.next_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListPullReviewsResponse", len)?;
        if !self.reviews.is_empty() {
            struct_ser.serialize_field("reviews", &self.reviews)?;
        }
        if let Some(v) = self.next_page.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nextPage", ToString::to_string(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListPullReviewsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "reviews",
            "next_page",
            "nextPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Reviews,
            NextPage,
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
                            "reviews" => Ok(GeneratedField::Reviews),
                            "nextPage" | "next_page" => Ok(GeneratedField::NextPage),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListPullReviewsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListPullReviewsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListPullReviewsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut reviews__ = None;
                let mut next_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Reviews => {
                            if reviews__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reviews"));
                            }
                            reviews__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextPage => {
                            if next_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPage"));
                            }
                            next_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListPullReviewsResponse {
                    reviews: reviews__.unwrap_or_default(),
                    next_page: next_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListPullReviewsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListPullsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        if self.page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListPullsRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if let Some(v) = self.per_page.as_ref() {
            struct_ser.serialize_field("perPage", v)?;
        }
        if let Some(v) = self.page.as_ref() {
            struct_ser.serialize_field("page", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListPullsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "per_page",
            "perPage",
            "page",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PerPage,
            Page,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "perPage" | "per_page" => Ok(GeneratedField::PerPage),
                            "page" => Ok(GeneratedField::Page),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListPullsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListPullsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListPullsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut per_page__ = None;
                let mut page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PerPage => {
                            if per_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("perPage"));
                            }
                            per_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Page => {
                            if page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("page"));
                            }
                            page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListPullsRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    per_page: per_page__,
                    page: page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListPullsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListPullsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.pulls.is_empty() {
            len += 1;
        }
        if self.next_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListPullsResponse", len)?;
        if !self.pulls.is_empty() {
            struct_ser.serialize_field("pulls", &self.pulls)?;
        }
        if let Some(v) = self.next_page.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nextPage", ToString::to_string(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListPullsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "pulls",
            "next_page",
            "nextPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Pulls,
            NextPage,
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
                            "pulls" => Ok(GeneratedField::Pulls),
                            "nextPage" | "next_page" => Ok(GeneratedField::NextPage),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListPullsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListPullsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListPullsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut pulls__ = None;
                let mut next_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Pulls => {
                            if pulls__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pulls"));
                            }
                            pulls__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextPage => {
                            if next_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPage"));
                            }
                            next_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListPullsResponse {
                    pulls: pulls__.unwrap_or_default(),
                    next_page: next_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListPullsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListWorkflowRunsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.per_page.is_some() {
            len += 1;
        }
        if self.page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListWorkflowRunsRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if let Some(v) = self.per_page.as_ref() {
            struct_ser.serialize_field("perPage", v)?;
        }
        if let Some(v) = self.page.as_ref() {
            struct_ser.serialize_field("page", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListWorkflowRunsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "per_page",
            "perPage",
            "page",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PerPage,
            Page,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "perPage" | "per_page" => Ok(GeneratedField::PerPage),
                            "page" => Ok(GeneratedField::Page),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListWorkflowRunsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListWorkflowRunsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListWorkflowRunsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut per_page__ = None;
                let mut page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PerPage => {
                            if per_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("perPage"));
                            }
                            per_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Page => {
                            if page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("page"));
                            }
                            page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListWorkflowRunsRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    per_page: per_page__,
                    page: page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListWorkflowRunsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListWorkflowRunsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.workflow_runs.is_empty() {
            len += 1;
        }
        if self.next_page.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ListWorkflowRunsResponse", len)?;
        if !self.workflow_runs.is_empty() {
            struct_ser.serialize_field("workflowRuns", &self.workflow_runs)?;
        }
        if let Some(v) = self.next_page.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nextPage", ToString::to_string(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListWorkflowRunsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "workflow_runs",
            "workflowRuns",
            "next_page",
            "nextPage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WorkflowRuns,
            NextPage,
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
                            "workflowRuns" | "workflow_runs" => Ok(GeneratedField::WorkflowRuns),
                            "nextPage" | "next_page" => Ok(GeneratedField::NextPage),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListWorkflowRunsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ListWorkflowRunsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListWorkflowRunsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut workflow_runs__ = None;
                let mut next_page__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WorkflowRuns => {
                            if workflow_runs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("workflowRuns"));
                            }
                            workflow_runs__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextPage => {
                            if next_page__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPage"));
                            }
                            next_page__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ListWorkflowRunsResponse {
                    workflow_runs: workflow_runs__.unwrap_or_default(),
                    next_page: next_page__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ListWorkflowRunsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PullRef {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.r#ref.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.PullRef", len)?;
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PullRef {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ref",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ref,
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
                            "ref" => Ok(GeneratedField::Ref),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PullRef;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.PullRef")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PullRef, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#ref__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(PullRef {
                    r#ref: r#ref__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.PullRef", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PullRequestSummary {
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
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.draft {
            len += 1;
        }
        if self.head.is_some() {
            len += 1;
        }
        if self.base.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.PullRequestSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if self.draft {
            struct_ser.serialize_field("draft", &self.draft)?;
        }
        if let Some(v) = self.head.as_ref() {
            struct_ser.serialize_field("head", v)?;
        }
        if let Some(v) = self.base.as_ref() {
            struct_ser.serialize_field("base", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PullRequestSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "draft",
            "head",
            "base",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Draft,
            Head,
            Base,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "draft" => Ok(GeneratedField::Draft),
                            "head" => Ok(GeneratedField::Head),
                            "base" => Ok(GeneratedField::Base),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = PullRequestSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.PullRequestSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PullRequestSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut draft__ = None;
                let mut head__ = None;
                let mut base__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Draft => {
                            if draft__.is_some() {
                                return Err(serde::de::Error::duplicate_field("draft"));
                            }
                            draft__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Head => {
                            if head__.is_some() {
                                return Err(serde::de::Error::duplicate_field("head"));
                            }
                            head__ = map_.next_value()?;
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(PullRequestSummary {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    draft: draft__.unwrap_or_default(),
                    head: head__,
                    base: base__,
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.PullRequestSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PutContentsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.path.is_empty() {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        if !self.content.is_empty() {
            len += 1;
        }
        if self.branch.is_some() {
            len += 1;
        }
        if self.sha.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.PutContentsRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.path.is_empty() {
            struct_ser.serialize_field("path", &self.path)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        if !self.content.is_empty() {
            struct_ser.serialize_field("content", &self.content)?;
        }
        if let Some(v) = self.branch.as_ref() {
            struct_ser.serialize_field("branch", v)?;
        }
        if let Some(v) = self.sha.as_ref() {
            struct_ser.serialize_field("sha", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PutContentsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "path",
            "message",
            "content",
            "branch",
            "sha",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Path,
            Message,
            Content,
            Branch,
            Sha,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "path" => Ok(GeneratedField::Path),
                            "message" => Ok(GeneratedField::Message),
                            "content" => Ok(GeneratedField::Content),
                            "branch" => Ok(GeneratedField::Branch),
                            "sha" => Ok(GeneratedField::Sha),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PutContentsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.PutContentsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PutContentsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut path__ = None;
                let mut message__ = None;
                let mut content__ = None;
                let mut branch__ = None;
                let mut sha__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Path => {
                            if path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("path"));
                            }
                            path__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Content => {
                            if content__.is_some() {
                                return Err(serde::de::Error::duplicate_field("content"));
                            }
                            content__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Branch => {
                            if branch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("branch"));
                            }
                            branch__ = map_.next_value()?;
                        }
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(PutContentsRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    path: path__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                    content: content__.unwrap_or_default(),
                    branch: branch__,
                    sha: sha__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.PutContentsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PutContentsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.content.is_some() {
            len += 1;
        }
        if self.commit.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.PutContentsResponse", len)?;
        if let Some(v) = self.content.as_ref() {
            struct_ser.serialize_field("content", v)?;
        }
        if let Some(v) = self.commit.as_ref() {
            struct_ser.serialize_field("commit", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PutContentsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "content",
            "commit",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Content,
            Commit,
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
                            "content" => Ok(GeneratedField::Content),
                            "commit" => Ok(GeneratedField::Commit),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PutContentsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.PutContentsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PutContentsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut content__ = None;
                let mut commit__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Content => {
                            if content__.is_some() {
                                return Err(serde::de::Error::duplicate_field("content"));
                            }
                            content__ = map_.next_value()?;
                        }
                        GeneratedField::Commit => {
                            if commit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commit"));
                            }
                            commit__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(PutContentsResponse {
                    content: content__,
                    commit: commit__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.PutContentsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ReviewSummary {
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
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.submitted_at.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.ReviewSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.submitted_at.as_ref() {
            struct_ser.serialize_field("submittedAt", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ReviewSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "submitted_at",
            "submittedAt",
            "author",
            "user",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            State,
            Body,
            HtmlUrl,
            SubmittedAt,
            Author,
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
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "submittedAt" | "submitted_at" => Ok(GeneratedField::SubmittedAt),
                            "user" | "author" => Ok(GeneratedField::Author),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ReviewSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.ReviewSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ReviewSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut submitted_at__ = None;
                let mut author__ = None;
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
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SubmittedAt => {
                            if submitted_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("submittedAt"));
                            }
                            submitted_at__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ReviewSummary {
                    id: id__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    submitted_at: submitted_at__,
                    author: author__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.ReviewSummary", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TreeEntry {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.path.is_empty() {
            len += 1;
        }
        if !self.mode.is_empty() {
            len += 1;
        }
        if !self.r#type.is_empty() {
            len += 1;
        }
        if !self.sha.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.TreeEntry", len)?;
        if !self.path.is_empty() {
            struct_ser.serialize_field("path", &self.path)?;
        }
        if !self.mode.is_empty() {
            struct_ser.serialize_field("mode", &self.mode)?;
        }
        if !self.r#type.is_empty() {
            struct_ser.serialize_field("type", &self.r#type)?;
        }
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TreeEntry {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "path",
            "mode",
            "type",
            "sha",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Path,
            Mode,
            Type,
            Sha,
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
                            "path" => Ok(GeneratedField::Path),
                            "mode" => Ok(GeneratedField::Mode),
                            "type" => Ok(GeneratedField::Type),
                            "sha" => Ok(GeneratedField::Sha),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TreeEntry;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.TreeEntry")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TreeEntry, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut path__ = None;
                let mut mode__ = None;
                let mut r#type__ = None;
                let mut sha__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Path => {
                            if path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("path"));
                            }
                            path__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Mode => {
                            if mode__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mode"));
                            }
                            mode__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(TreeEntry {
                    path: path__.unwrap_or_default(),
                    mode: mode__.unwrap_or_default(),
                    r#type: r#type__.unwrap_or_default(),
                    sha: sha__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.TreeEntry", FIELDS, GeneratedVisitor)
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
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.issue_number != 0 {
            len += 1;
        }
        if self.title.is_some() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if self.state.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdateIssueRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.issue_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("issueNumber", ToString::to_string(&self.issue_number).as_str())?;
        }
        if let Some(v) = self.title.as_ref() {
            struct_ser.serialize_field("title", v)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if let Some(v) = self.state.as_ref() {
            struct_ser.serialize_field("state", v)?;
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
            "owner",
            "repository",
            "issue_number",
            "issueNumber",
            "title",
            "body",
            "state",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            IssueNumber,
            Title,
            Body,
            State,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "issueNumber" | "issue_number" => Ok(GeneratedField::IssueNumber),
                            "title" => Ok(GeneratedField::Title),
                            "body" => Ok(GeneratedField::Body),
                            "state" => Ok(GeneratedField::State),
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
                formatter.write_str("struct kei.connectors.github.v1.UpdateIssueRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateIssueRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut issue_number__ = None;
                let mut title__ = None;
                let mut body__ = None;
                let mut state__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IssueNumber => {
                            if issue_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("issueNumber"));
                            }
                            issue_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = map_.next_value()?;
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateIssueRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    issue_number: issue_number__.unwrap_or_default(),
                    title: title__,
                    body: body__,
                    state: state__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdateIssueRequest", FIELDS, GeneratedVisitor)
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
        if self.id != 0 {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdateIssueResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
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
            "id",
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = UpdateIssueResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.UpdateIssueResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateIssueResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(UpdateIssueResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdateIssueResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdatePullRequestRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if self.pull_number != 0 {
            len += 1;
        }
        if self.title.is_some() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if self.state.is_some() {
            len += 1;
        }
        if self.base.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdatePullRequestRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if self.pull_number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("pullNumber", ToString::to_string(&self.pull_number).as_str())?;
        }
        if let Some(v) = self.title.as_ref() {
            struct_ser.serialize_field("title", v)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if let Some(v) = self.state.as_ref() {
            struct_ser.serialize_field("state", v)?;
        }
        if let Some(v) = self.base.as_ref() {
            struct_ser.serialize_field("base", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdatePullRequestRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "pull_number",
            "pullNumber",
            "title",
            "body",
            "state",
            "base",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            PullNumber,
            Title,
            Body,
            State,
            Base,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "pullNumber" | "pull_number" => Ok(GeneratedField::PullNumber),
                            "title" => Ok(GeneratedField::Title),
                            "body" => Ok(GeneratedField::Body),
                            "state" => Ok(GeneratedField::State),
                            "base" => Ok(GeneratedField::Base),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdatePullRequestRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.UpdatePullRequestRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdatePullRequestRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut pull_number__ = None;
                let mut title__ = None;
                let mut body__ = None;
                let mut state__ = None;
                let mut base__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PullNumber => {
                            if pull_number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pullNumber"));
                            }
                            pull_number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = map_.next_value()?;
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = map_.next_value()?;
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdatePullRequestRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    pull_number: pull_number__.unwrap_or_default(),
                    title: title__,
                    body: body__,
                    state: state__,
                    base: base__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdatePullRequestRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdatePullRequestResponse {
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
        if self.number != 0 {
            len += 1;
        }
        if !self.title.is_empty() {
            len += 1;
        }
        if !self.state.is_empty() {
            len += 1;
        }
        if self.body.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        if self.draft {
            len += 1;
        }
        if self.merged_at.is_some() {
            len += 1;
        }
        if self.head.is_some() {
            len += 1;
        }
        if self.base.is_some() {
            len += 1;
        }
        if self.author.is_some() {
            len += 1;
        }
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdatePullRequestResponse", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.number != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("number", ToString::to_string(&self.number).as_str())?;
        }
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        if !self.state.is_empty() {
            struct_ser.serialize_field("state", &self.state)?;
        }
        if let Some(v) = self.body.as_ref() {
            struct_ser.serialize_field("body", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        if self.draft {
            struct_ser.serialize_field("draft", &self.draft)?;
        }
        if let Some(v) = self.merged_at.as_ref() {
            struct_ser.serialize_field("mergedAt", v)?;
        }
        if let Some(v) = self.head.as_ref() {
            struct_ser.serialize_field("head", v)?;
        }
        if let Some(v) = self.base.as_ref() {
            struct_ser.serialize_field("base", v)?;
        }
        if let Some(v) = self.author.as_ref() {
            struct_ser.serialize_field("user", v)?;
        }
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdatePullRequestResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "number",
            "title",
            "state",
            "body",
            "html_url",
            "htmlUrl",
            "draft",
            "merged_at",
            "mergedAt",
            "head",
            "base",
            "author",
            "user",
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Number,
            Title,
            State,
            Body,
            HtmlUrl,
            Draft,
            MergedAt,
            Head,
            Base,
            Author,
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
                            "number" => Ok(GeneratedField::Number),
                            "title" => Ok(GeneratedField::Title),
                            "state" => Ok(GeneratedField::State),
                            "body" => Ok(GeneratedField::Body),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            "draft" => Ok(GeneratedField::Draft),
                            "mergedAt" | "merged_at" => Ok(GeneratedField::MergedAt),
                            "head" => Ok(GeneratedField::Head),
                            "base" => Ok(GeneratedField::Base),
                            "user" | "author" => Ok(GeneratedField::Author),
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
            type Value = UpdatePullRequestResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.UpdatePullRequestResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdatePullRequestResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut number__ = None;
                let mut title__ = None;
                let mut state__ = None;
                let mut body__ = None;
                let mut html_url__ = None;
                let mut draft__ = None;
                let mut merged_at__ = None;
                let mut head__ = None;
                let mut base__ = None;
                let mut author__ = None;
                let mut labels__ = None;
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
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Body => {
                            if body__.is_some() {
                                return Err(serde::de::Error::duplicate_field("body"));
                            }
                            body__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Draft => {
                            if draft__.is_some() {
                                return Err(serde::de::Error::duplicate_field("draft"));
                            }
                            draft__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MergedAt => {
                            if merged_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mergedAt"));
                            }
                            merged_at__ = map_.next_value()?;
                        }
                        GeneratedField::Head => {
                            if head__.is_some() {
                                return Err(serde::de::Error::duplicate_field("head"));
                            }
                            head__ = map_.next_value()?;
                        }
                        GeneratedField::Base => {
                            if base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("base"));
                            }
                            base__ = map_.next_value()?;
                        }
                        GeneratedField::Author => {
                            if author__.is_some() {
                                return Err(serde::de::Error::duplicate_field("user"));
                            }
                            author__ = map_.next_value()?;
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
                Ok(UpdatePullRequestResponse {
                    id: id__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    title: title__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    body: body__,
                    html_url: html_url__.unwrap_or_default(),
                    draft: draft__.unwrap_or_default(),
                    merged_at: merged_at__,
                    head: head__,
                    base: base__,
                    author: author__,
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdatePullRequestResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateRefRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.owner.is_empty() {
            len += 1;
        }
        if !self.repository.is_empty() {
            len += 1;
        }
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if !self.sha.is_empty() {
            len += 1;
        }
        if self.force.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdateRefRequest", len)?;
        if !self.owner.is_empty() {
            struct_ser.serialize_field("owner", &self.owner)?;
        }
        if !self.repository.is_empty() {
            struct_ser.serialize_field("repository", &self.repository)?;
        }
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if !self.sha.is_empty() {
            struct_ser.serialize_field("sha", &self.sha)?;
        }
        if let Some(v) = self.force.as_ref() {
            struct_ser.serialize_field("force", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateRefRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "owner",
            "repository",
            "ref",
            "sha",
            "force",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Owner,
            Repository,
            Ref,
            Sha,
            Force,
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
                            "owner" => Ok(GeneratedField::Owner),
                            "repository" => Ok(GeneratedField::Repository),
                            "ref" => Ok(GeneratedField::Ref),
                            "sha" => Ok(GeneratedField::Sha),
                            "force" => Ok(GeneratedField::Force),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdateRefRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.UpdateRefRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateRefRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut owner__ = None;
                let mut repository__ = None;
                let mut r#ref__ = None;
                let mut sha__ = None;
                let mut force__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Repository => {
                            if repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("repository"));
                            }
                            repository__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Sha => {
                            if sha__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sha"));
                            }
                            sha__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Force => {
                            if force__.is_some() {
                                return Err(serde::de::Error::duplicate_field("force"));
                            }
                            force__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateRefRequest {
                    owner: owner__.unwrap_or_default(),
                    repository: repository__.unwrap_or_default(),
                    r#ref: r#ref__.unwrap_or_default(),
                    sha: sha__.unwrap_or_default(),
                    force: force__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdateRefRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateRefResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.r#ref.is_empty() {
            len += 1;
        }
        if self.object.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.UpdateRefResponse", len)?;
        if !self.r#ref.is_empty() {
            struct_ser.serialize_field("ref", &self.r#ref)?;
        }
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateRefResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ref",
            "object",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ref,
            Object,
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
                            "ref" => Ok(GeneratedField::Ref),
                            "object" => Ok(GeneratedField::Object),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdateRefResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.UpdateRefResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateRefResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#ref__ = None;
                let mut object__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ref => {
                            if r#ref__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ref"));
                            }
                            r#ref__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateRefResponse {
                    r#ref: r#ref__.unwrap_or_default(),
                    object: object__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.UpdateRefResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for User {
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
        if !self.login.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.User", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.login.is_empty() {
            struct_ser.serialize_field("login", &self.login)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for User {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "login",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Login,
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
                            "login" => Ok(GeneratedField::Login),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = User;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.User")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<User, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut login__ = None;
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
                        GeneratedField::Login => {
                            if login__.is_some() {
                                return Err(serde::de::Error::duplicate_field("login"));
                            }
                            login__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(User {
                    id: id__.unwrap_or_default(),
                    login: login__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.User", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for WorkflowRunSummary {
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
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.head_branch.is_empty() {
            len += 1;
        }
        if !self.status.is_empty() {
            len += 1;
        }
        if self.conclusion.is_some() {
            len += 1;
        }
        if !self.html_url.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.github.v1.WorkflowRunSummary", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.head_branch.is_empty() {
            struct_ser.serialize_field("headBranch", &self.head_branch)?;
        }
        if !self.status.is_empty() {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if let Some(v) = self.conclusion.as_ref() {
            struct_ser.serialize_field("conclusion", v)?;
        }
        if !self.html_url.is_empty() {
            struct_ser.serialize_field("htmlUrl", &self.html_url)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for WorkflowRunSummary {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "name",
            "head_branch",
            "headBranch",
            "status",
            "conclusion",
            "html_url",
            "htmlUrl",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Name,
            HeadBranch,
            Status,
            Conclusion,
            HtmlUrl,
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
                            "headBranch" | "head_branch" => Ok(GeneratedField::HeadBranch),
                            "status" => Ok(GeneratedField::Status),
                            "conclusion" => Ok(GeneratedField::Conclusion),
                            "htmlUrl" | "html_url" => Ok(GeneratedField::HtmlUrl),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = WorkflowRunSummary;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.github.v1.WorkflowRunSummary")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<WorkflowRunSummary, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut name__ = None;
                let mut head_branch__ = None;
                let mut status__ = None;
                let mut conclusion__ = None;
                let mut html_url__ = None;
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
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HeadBranch => {
                            if head_branch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("headBranch"));
                            }
                            head_branch__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Conclusion => {
                            if conclusion__.is_some() {
                                return Err(serde::de::Error::duplicate_field("conclusion"));
                            }
                            conclusion__ = map_.next_value()?;
                        }
                        GeneratedField::HtmlUrl => {
                            if html_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("htmlUrl"));
                            }
                            html_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(WorkflowRunSummary {
                    id: id__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    head_branch: head_branch__.unwrap_or_default(),
                    status: status__.unwrap_or_default(),
                    conclusion: conclusion__,
                    html_url: html_url__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.github.v1.WorkflowRunSummary", FIELDS, GeneratedVisitor)
    }
}
