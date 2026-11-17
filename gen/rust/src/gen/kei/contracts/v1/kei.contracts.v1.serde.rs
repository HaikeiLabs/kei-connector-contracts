// @generated
impl serde::Serialize for AccountModel {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "ACCOUNT_MODEL_UNSPECIFIED",
            Self::PerUser => "ACCOUNT_MODEL_PER_USER",
            Self::Shared => "ACCOUNT_MODEL_SHARED",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for AccountModel {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ACCOUNT_MODEL_UNSPECIFIED",
            "ACCOUNT_MODEL_PER_USER",
            "ACCOUNT_MODEL_SHARED",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AccountModel;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "ACCOUNT_MODEL_UNSPECIFIED" => Ok(AccountModel::Unspecified),
                    "ACCOUNT_MODEL_PER_USER" => Ok(AccountModel::PerUser),
                    "ACCOUNT_MODEL_SHARED" => Ok(AccountModel::Shared),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for Connection {
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
        if self.kind != 0 {
            len += 1;
        }
        if !self.auth_method.is_empty() {
            len += 1;
        }
        if self.account_model != 0 {
            len += 1;
        }
        if !self.base_url.is_empty() {
            len += 1;
        }
        if !self.base_url_field.is_empty() {
            len += 1;
        }
        if !self.setup_fields.is_empty() {
            len += 1;
        }
        if !self.scopes.is_empty() {
            len += 1;
        }
        if !self.engine.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.Connection", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if self.kind != 0 {
            let v = ConnectionKind::try_from(self.kind)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.kind)))?;
            struct_ser.serialize_field("kind", &v)?;
        }
        if !self.auth_method.is_empty() {
            struct_ser.serialize_field("authMethod", &self.auth_method)?;
        }
        if self.account_model != 0 {
            let v = AccountModel::try_from(self.account_model)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.account_model)))?;
            struct_ser.serialize_field("accountModel", &v)?;
        }
        if !self.base_url.is_empty() {
            struct_ser.serialize_field("baseUrl", &self.base_url)?;
        }
        if !self.base_url_field.is_empty() {
            struct_ser.serialize_field("baseUrlField", &self.base_url_field)?;
        }
        if !self.setup_fields.is_empty() {
            struct_ser.serialize_field("setupFields", &self.setup_fields)?;
        }
        if !self.scopes.is_empty() {
            struct_ser.serialize_field("scopes", &self.scopes)?;
        }
        if !self.engine.is_empty() {
            struct_ser.serialize_field("engine", &self.engine)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Connection {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "kind",
            "auth_method",
            "authMethod",
            "account_model",
            "accountModel",
            "base_url",
            "baseUrl",
            "base_url_field",
            "baseUrlField",
            "setup_fields",
            "setupFields",
            "scopes",
            "engine",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Kind,
            AuthMethod,
            AccountModel,
            BaseUrl,
            BaseUrlField,
            SetupFields,
            Scopes,
            Engine,
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
                            "kind" => Ok(GeneratedField::Kind),
                            "authMethod" | "auth_method" => Ok(GeneratedField::AuthMethod),
                            "accountModel" | "account_model" => Ok(GeneratedField::AccountModel),
                            "baseUrl" | "base_url" => Ok(GeneratedField::BaseUrl),
                            "baseUrlField" | "base_url_field" => Ok(GeneratedField::BaseUrlField),
                            "setupFields" | "setup_fields" => Ok(GeneratedField::SetupFields),
                            "scopes" => Ok(GeneratedField::Scopes),
                            "engine" => Ok(GeneratedField::Engine),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Connection;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.Connection")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Connection, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut kind__ = None;
                let mut auth_method__ = None;
                let mut account_model__ = None;
                let mut base_url__ = None;
                let mut base_url_field__ = None;
                let mut setup_fields__ = None;
                let mut scopes__ = None;
                let mut engine__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Kind => {
                            if kind__.is_some() {
                                return Err(serde::de::Error::duplicate_field("kind"));
                            }
                            kind__ = Some(map_.next_value::<ConnectionKind>()? as i32);
                        }
                        GeneratedField::AuthMethod => {
                            if auth_method__.is_some() {
                                return Err(serde::de::Error::duplicate_field("authMethod"));
                            }
                            auth_method__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AccountModel => {
                            if account_model__.is_some() {
                                return Err(serde::de::Error::duplicate_field("accountModel"));
                            }
                            account_model__ = Some(map_.next_value::<AccountModel>()? as i32);
                        }
                        GeneratedField::BaseUrl => {
                            if base_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseUrl"));
                            }
                            base_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BaseUrlField => {
                            if base_url_field__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseUrlField"));
                            }
                            base_url_field__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SetupFields => {
                            if setup_fields__.is_some() {
                                return Err(serde::de::Error::duplicate_field("setupFields"));
                            }
                            setup_fields__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Scopes => {
                            if scopes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("scopes"));
                            }
                            scopes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Engine => {
                            if engine__.is_some() {
                                return Err(serde::de::Error::duplicate_field("engine"));
                            }
                            engine__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Connection {
                    name: name__.unwrap_or_default(),
                    kind: kind__.unwrap_or_default(),
                    auth_method: auth_method__.unwrap_or_default(),
                    account_model: account_model__.unwrap_or_default(),
                    base_url: base_url__.unwrap_or_default(),
                    base_url_field: base_url_field__.unwrap_or_default(),
                    setup_fields: setup_fields__.unwrap_or_default(),
                    scopes: scopes__.unwrap_or_default(),
                    engine: engine__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.Connection", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ConnectionKind {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "CONNECTION_KIND_UNSPECIFIED",
            Self::HttpBearer => "CONNECTION_KIND_HTTP_BEARER",
            Self::HttpHeader => "CONNECTION_KIND_HTTP_HEADER",
            Self::HttpBasic => "CONNECTION_KIND_HTTP_BASIC",
            Self::SqlDsn => "CONNECTION_KIND_SQL_DSN",
            Self::KeyPairJwt => "CONNECTION_KIND_KEY_PAIR_JWT",
            Self::TokenExchange => "CONNECTION_KIND_TOKEN_EXCHANGE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ConnectionKind {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "CONNECTION_KIND_UNSPECIFIED",
            "CONNECTION_KIND_HTTP_BEARER",
            "CONNECTION_KIND_HTTP_HEADER",
            "CONNECTION_KIND_HTTP_BASIC",
            "CONNECTION_KIND_SQL_DSN",
            "CONNECTION_KIND_KEY_PAIR_JWT",
            "CONNECTION_KIND_TOKEN_EXCHANGE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ConnectionKind;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "CONNECTION_KIND_UNSPECIFIED" => Ok(ConnectionKind::Unspecified),
                    "CONNECTION_KIND_HTTP_BEARER" => Ok(ConnectionKind::HttpBearer),
                    "CONNECTION_KIND_HTTP_HEADER" => Ok(ConnectionKind::HttpHeader),
                    "CONNECTION_KIND_HTTP_BASIC" => Ok(ConnectionKind::HttpBasic),
                    "CONNECTION_KIND_SQL_DSN" => Ok(ConnectionKind::SqlDsn),
                    "CONNECTION_KIND_KEY_PAIR_JWT" => Ok(ConnectionKind::KeyPairJwt),
                    "CONNECTION_KIND_TOKEN_EXCHANGE" => Ok(ConnectionKind::TokenExchange),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for Connector {
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
        if !self.base_url.is_empty() {
            len += 1;
        }
        if !self.connections.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.Connector", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.base_url.is_empty() {
            struct_ser.serialize_field("baseUrl", &self.base_url)?;
        }
        if !self.connections.is_empty() {
            struct_ser.serialize_field("connections", &self.connections)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Connector {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "base_url",
            "baseUrl",
            "connections",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            BaseUrl,
            Connections,
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
                            "baseUrl" | "base_url" => Ok(GeneratedField::BaseUrl),
                            "connections" => Ok(GeneratedField::Connections),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Connector;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.Connector")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Connector, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut base_url__ = None;
                let mut connections__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BaseUrl => {
                            if base_url__.is_some() {
                                return Err(serde::de::Error::duplicate_field("baseUrl"));
                            }
                            base_url__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Connections => {
                            if connections__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connections"));
                            }
                            connections__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Connector {
                    id: id__.unwrap_or_default(),
                    base_url: base_url__.unwrap_or_default(),
                    connections: connections__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.Connector", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CredentialPlacement {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "CREDENTIAL_PLACEMENT_UNSPECIFIED",
            Self::HeaderBearer => "CREDENTIAL_PLACEMENT_HEADER_BEARER",
            Self::HeaderCustom => "CREDENTIAL_PLACEMENT_HEADER_CUSTOM",
            Self::Query => "CREDENTIAL_PLACEMENT_QUERY",
            Self::Basic => "CREDENTIAL_PLACEMENT_BASIC",
            Self::None => "CREDENTIAL_PLACEMENT_NONE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for CredentialPlacement {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "CREDENTIAL_PLACEMENT_UNSPECIFIED",
            "CREDENTIAL_PLACEMENT_HEADER_BEARER",
            "CREDENTIAL_PLACEMENT_HEADER_CUSTOM",
            "CREDENTIAL_PLACEMENT_QUERY",
            "CREDENTIAL_PLACEMENT_BASIC",
            "CREDENTIAL_PLACEMENT_NONE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CredentialPlacement;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "CREDENTIAL_PLACEMENT_UNSPECIFIED" => Ok(CredentialPlacement::Unspecified),
                    "CREDENTIAL_PLACEMENT_HEADER_BEARER" => Ok(CredentialPlacement::HeaderBearer),
                    "CREDENTIAL_PLACEMENT_HEADER_CUSTOM" => Ok(CredentialPlacement::HeaderCustom),
                    "CREDENTIAL_PLACEMENT_QUERY" => Ok(CredentialPlacement::Query),
                    "CREDENTIAL_PLACEMENT_BASIC" => Ok(CredentialPlacement::Basic),
                    "CREDENTIAL_PLACEMENT_NONE" => Ok(CredentialPlacement::None),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for ErrorKind {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "ERROR_KIND_UNSPECIFIED",
            Self::NotFound => "ERROR_KIND_NOT_FOUND",
            Self::Unauthorized => "ERROR_KIND_UNAUTHORIZED",
            Self::Forbidden => "ERROR_KIND_FORBIDDEN",
            Self::RateLimited => "ERROR_KIND_RATE_LIMITED",
            Self::Conflict => "ERROR_KIND_CONFLICT",
            Self::Invalid => "ERROR_KIND_INVALID",
            Self::Unavailable => "ERROR_KIND_UNAVAILABLE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ErrorKind {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ERROR_KIND_UNSPECIFIED",
            "ERROR_KIND_NOT_FOUND",
            "ERROR_KIND_UNAUTHORIZED",
            "ERROR_KIND_FORBIDDEN",
            "ERROR_KIND_RATE_LIMITED",
            "ERROR_KIND_CONFLICT",
            "ERROR_KIND_INVALID",
            "ERROR_KIND_UNAVAILABLE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ErrorKind;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "ERROR_KIND_UNSPECIFIED" => Ok(ErrorKind::Unspecified),
                    "ERROR_KIND_NOT_FOUND" => Ok(ErrorKind::NotFound),
                    "ERROR_KIND_UNAUTHORIZED" => Ok(ErrorKind::Unauthorized),
                    "ERROR_KIND_FORBIDDEN" => Ok(ErrorKind::Forbidden),
                    "ERROR_KIND_RATE_LIMITED" => Ok(ErrorKind::RateLimited),
                    "ERROR_KIND_CONFLICT" => Ok(ErrorKind::Conflict),
                    "ERROR_KIND_INVALID" => Ok(ErrorKind::Invalid),
                    "ERROR_KIND_UNAVAILABLE" => Ok(ErrorKind::Unavailable),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for ErrorMapping {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.status != 0 {
            len += 1;
        }
        if self.kind != 0 {
            len += 1;
        }
        if self.retryable {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.ErrorMapping", len)?;
        if self.status != 0 {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if self.kind != 0 {
            let v = ErrorKind::try_from(self.kind)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.kind)))?;
            struct_ser.serialize_field("kind", &v)?;
        }
        if self.retryable {
            struct_ser.serialize_field("retryable", &self.retryable)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ErrorMapping {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "status",
            "kind",
            "retryable",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Status,
            Kind,
            Retryable,
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
                            "status" => Ok(GeneratedField::Status),
                            "kind" => Ok(GeneratedField::Kind),
                            "retryable" => Ok(GeneratedField::Retryable),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ErrorMapping;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.ErrorMapping")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ErrorMapping, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut status__ = None;
                let mut kind__ = None;
                let mut retryable__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Kind => {
                            if kind__.is_some() {
                                return Err(serde::de::Error::duplicate_field("kind"));
                            }
                            kind__ = Some(map_.next_value::<ErrorKind>()? as i32);
                        }
                        GeneratedField::Retryable => {
                            if retryable__.is_some() {
                                return Err(serde::de::Error::duplicate_field("retryable"));
                            }
                            retryable__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ErrorMapping {
                    status: status__.unwrap_or_default(),
                    kind: kind__.unwrap_or_default(),
                    retryable: retryable__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.ErrorMapping", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Operation {
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
        if !self.connector.is_empty() {
            len += 1;
        }
        if self.class != 0 {
            len += 1;
        }
        if self.destructive {
            len += 1;
        }
        if !self.http_method.is_empty() {
            len += 1;
        }
        if !self.path_template.is_empty() {
            len += 1;
        }
        if !self.scopes.is_empty() {
            len += 1;
        }
        if self.credential_placement != 0 {
            len += 1;
        }
        if !self.method_name.is_empty() {
            len += 1;
        }
        if !self.connections.is_empty() {
            len += 1;
        }
        if !self.resource_template.is_empty() {
            len += 1;
        }
        if !self.errors.is_empty() {
            len += 1;
        }
        if self.deprecated {
            len += 1;
        }
        if !self.replaced_by.is_empty() {
            len += 1;
        }
        if !self.graphql_document.is_empty() {
            len += 1;
        }
        if self.positional_helper {
            len += 1;
        }
        if self.semantics.is_some() {
            len += 1;
        }
        if !self.cursor_field.is_empty() {
            len += 1;
        }
        if !self.next_cursor_field.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.Operation", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.connector.is_empty() {
            struct_ser.serialize_field("connector", &self.connector)?;
        }
        if self.class != 0 {
            let v = OperationClass::try_from(self.class)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.class)))?;
            struct_ser.serialize_field("class", &v)?;
        }
        if self.destructive {
            struct_ser.serialize_field("destructive", &self.destructive)?;
        }
        if !self.http_method.is_empty() {
            struct_ser.serialize_field("httpMethod", &self.http_method)?;
        }
        if !self.path_template.is_empty() {
            struct_ser.serialize_field("pathTemplate", &self.path_template)?;
        }
        if !self.scopes.is_empty() {
            struct_ser.serialize_field("scopes", &self.scopes)?;
        }
        if self.credential_placement != 0 {
            let v = CredentialPlacement::try_from(self.credential_placement)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.credential_placement)))?;
            struct_ser.serialize_field("credentialPlacement", &v)?;
        }
        if !self.method_name.is_empty() {
            struct_ser.serialize_field("methodName", &self.method_name)?;
        }
        if !self.connections.is_empty() {
            struct_ser.serialize_field("connections", &self.connections)?;
        }
        if !self.resource_template.is_empty() {
            struct_ser.serialize_field("resourceTemplate", &self.resource_template)?;
        }
        if !self.errors.is_empty() {
            struct_ser.serialize_field("errors", &self.errors)?;
        }
        if self.deprecated {
            struct_ser.serialize_field("deprecated", &self.deprecated)?;
        }
        if !self.replaced_by.is_empty() {
            struct_ser.serialize_field("replacedBy", &self.replaced_by)?;
        }
        if !self.graphql_document.is_empty() {
            struct_ser.serialize_field("graphqlDocument", &self.graphql_document)?;
        }
        if self.positional_helper {
            struct_ser.serialize_field("positionalHelper", &self.positional_helper)?;
        }
        if let Some(v) = self.semantics.as_ref() {
            struct_ser.serialize_field("semantics", v)?;
        }
        if !self.cursor_field.is_empty() {
            struct_ser.serialize_field("cursorField", &self.cursor_field)?;
        }
        if !self.next_cursor_field.is_empty() {
            struct_ser.serialize_field("nextCursorField", &self.next_cursor_field)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Operation {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "connector",
            "class",
            "destructive",
            "http_method",
            "httpMethod",
            "path_template",
            "pathTemplate",
            "scopes",
            "credential_placement",
            "credentialPlacement",
            "method_name",
            "methodName",
            "connections",
            "resource_template",
            "resourceTemplate",
            "errors",
            "deprecated",
            "replaced_by",
            "replacedBy",
            "graphql_document",
            "graphqlDocument",
            "positional_helper",
            "positionalHelper",
            "semantics",
            "cursor_field",
            "cursorField",
            "next_cursor_field",
            "nextCursorField",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Connector,
            Class,
            Destructive,
            HttpMethod,
            PathTemplate,
            Scopes,
            CredentialPlacement,
            MethodName,
            Connections,
            ResourceTemplate,
            Errors,
            Deprecated,
            ReplacedBy,
            GraphqlDocument,
            PositionalHelper,
            Semantics,
            CursorField,
            NextCursorField,
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
                            "connector" => Ok(GeneratedField::Connector),
                            "class" => Ok(GeneratedField::Class),
                            "destructive" => Ok(GeneratedField::Destructive),
                            "httpMethod" | "http_method" => Ok(GeneratedField::HttpMethod),
                            "pathTemplate" | "path_template" => Ok(GeneratedField::PathTemplate),
                            "scopes" => Ok(GeneratedField::Scopes),
                            "credentialPlacement" | "credential_placement" => Ok(GeneratedField::CredentialPlacement),
                            "methodName" | "method_name" => Ok(GeneratedField::MethodName),
                            "connections" => Ok(GeneratedField::Connections),
                            "resourceTemplate" | "resource_template" => Ok(GeneratedField::ResourceTemplate),
                            "errors" => Ok(GeneratedField::Errors),
                            "deprecated" => Ok(GeneratedField::Deprecated),
                            "replacedBy" | "replaced_by" => Ok(GeneratedField::ReplacedBy),
                            "graphqlDocument" | "graphql_document" => Ok(GeneratedField::GraphqlDocument),
                            "positionalHelper" | "positional_helper" => Ok(GeneratedField::PositionalHelper),
                            "semantics" => Ok(GeneratedField::Semantics),
                            "cursorField" | "cursor_field" => Ok(GeneratedField::CursorField),
                            "nextCursorField" | "next_cursor_field" => Ok(GeneratedField::NextCursorField),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Operation;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.Operation")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Operation, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut connector__ = None;
                let mut class__ = None;
                let mut destructive__ = None;
                let mut http_method__ = None;
                let mut path_template__ = None;
                let mut scopes__ = None;
                let mut credential_placement__ = None;
                let mut method_name__ = None;
                let mut connections__ = None;
                let mut resource_template__ = None;
                let mut errors__ = None;
                let mut deprecated__ = None;
                let mut replaced_by__ = None;
                let mut graphql_document__ = None;
                let mut positional_helper__ = None;
                let mut semantics__ = None;
                let mut cursor_field__ = None;
                let mut next_cursor_field__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Connector => {
                            if connector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connector"));
                            }
                            connector__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Class => {
                            if class__.is_some() {
                                return Err(serde::de::Error::duplicate_field("class"));
                            }
                            class__ = Some(map_.next_value::<OperationClass>()? as i32);
                        }
                        GeneratedField::Destructive => {
                            if destructive__.is_some() {
                                return Err(serde::de::Error::duplicate_field("destructive"));
                            }
                            destructive__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HttpMethod => {
                            if http_method__.is_some() {
                                return Err(serde::de::Error::duplicate_field("httpMethod"));
                            }
                            http_method__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PathTemplate => {
                            if path_template__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pathTemplate"));
                            }
                            path_template__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Scopes => {
                            if scopes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("scopes"));
                            }
                            scopes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CredentialPlacement => {
                            if credential_placement__.is_some() {
                                return Err(serde::de::Error::duplicate_field("credentialPlacement"));
                            }
                            credential_placement__ = Some(map_.next_value::<CredentialPlacement>()? as i32);
                        }
                        GeneratedField::MethodName => {
                            if method_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("methodName"));
                            }
                            method_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Connections => {
                            if connections__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connections"));
                            }
                            connections__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ResourceTemplate => {
                            if resource_template__.is_some() {
                                return Err(serde::de::Error::duplicate_field("resourceTemplate"));
                            }
                            resource_template__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Errors => {
                            if errors__.is_some() {
                                return Err(serde::de::Error::duplicate_field("errors"));
                            }
                            errors__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Deprecated => {
                            if deprecated__.is_some() {
                                return Err(serde::de::Error::duplicate_field("deprecated"));
                            }
                            deprecated__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ReplacedBy => {
                            if replaced_by__.is_some() {
                                return Err(serde::de::Error::duplicate_field("replacedBy"));
                            }
                            replaced_by__ = Some(map_.next_value()?);
                        }
                        GeneratedField::GraphqlDocument => {
                            if graphql_document__.is_some() {
                                return Err(serde::de::Error::duplicate_field("graphqlDocument"));
                            }
                            graphql_document__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PositionalHelper => {
                            if positional_helper__.is_some() {
                                return Err(serde::de::Error::duplicate_field("positionalHelper"));
                            }
                            positional_helper__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Semantics => {
                            if semantics__.is_some() {
                                return Err(serde::de::Error::duplicate_field("semantics"));
                            }
                            semantics__ = map_.next_value()?;
                        }
                        GeneratedField::CursorField => {
                            if cursor_field__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cursorField"));
                            }
                            cursor_field__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextCursorField => {
                            if next_cursor_field__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextCursorField"));
                            }
                            next_cursor_field__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Operation {
                    id: id__.unwrap_or_default(),
                    connector: connector__.unwrap_or_default(),
                    class: class__.unwrap_or_default(),
                    destructive: destructive__.unwrap_or_default(),
                    http_method: http_method__.unwrap_or_default(),
                    path_template: path_template__.unwrap_or_default(),
                    scopes: scopes__.unwrap_or_default(),
                    credential_placement: credential_placement__.unwrap_or_default(),
                    method_name: method_name__.unwrap_or_default(),
                    connections: connections__.unwrap_or_default(),
                    resource_template: resource_template__.unwrap_or_default(),
                    errors: errors__.unwrap_or_default(),
                    deprecated: deprecated__.unwrap_or_default(),
                    replaced_by: replaced_by__.unwrap_or_default(),
                    graphql_document: graphql_document__.unwrap_or_default(),
                    positional_helper: positional_helper__.unwrap_or_default(),
                    semantics: semantics__,
                    cursor_field: cursor_field__.unwrap_or_default(),
                    next_cursor_field: next_cursor_field__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.Operation", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for OperationClass {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "OPERATION_CLASS_UNSPECIFIED",
            Self::Read => "OPERATION_CLASS_READ",
            Self::Write => "OPERATION_CLASS_WRITE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for OperationClass {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "OPERATION_CLASS_UNSPECIFIED",
            "OPERATION_CLASS_READ",
            "OPERATION_CLASS_WRITE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = OperationClass;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "OPERATION_CLASS_UNSPECIFIED" => Ok(OperationClass::Unspecified),
                    "OPERATION_CLASS_READ" => Ok(OperationClass::Read),
                    "OPERATION_CLASS_WRITE" => Ok(OperationClass::Write),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for Param {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.r#in != 0 {
            len += 1;
        }
        if !self.name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.Param", len)?;
        if self.r#in != 0 {
            let v = ParamIn::try_from(self.r#in)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.r#in)))?;
            struct_ser.serialize_field("in", &v)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Param {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "in",
            "name",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            In,
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
                            "in" => Ok(GeneratedField::In),
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
            type Value = Param;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.Param")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Param, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#in__ = None;
                let mut name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::In => {
                            if r#in__.is_some() {
                                return Err(serde::de::Error::duplicate_field("in"));
                            }
                            r#in__ = Some(map_.next_value::<ParamIn>()? as i32);
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
                Ok(Param {
                    r#in: r#in__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.Param", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ParamIn {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "PARAM_IN_UNSPECIFIED",
            Self::Path => "PARAM_IN_PATH",
            Self::Query => "PARAM_IN_QUERY",
            Self::Body => "PARAM_IN_BODY",
            Self::Header => "PARAM_IN_HEADER",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ParamIn {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "PARAM_IN_UNSPECIFIED",
            "PARAM_IN_PATH",
            "PARAM_IN_QUERY",
            "PARAM_IN_BODY",
            "PARAM_IN_HEADER",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ParamIn;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "PARAM_IN_UNSPECIFIED" => Ok(ParamIn::Unspecified),
                    "PARAM_IN_PATH" => Ok(ParamIn::Path),
                    "PARAM_IN_QUERY" => Ok(ParamIn::Query),
                    "PARAM_IN_BODY" => Ok(ParamIn::Body),
                    "PARAM_IN_HEADER" => Ok(ParamIn::Header),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for SemanticAction {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "SEMANTIC_ACTION_UNSPECIFIED",
            Self::Read => "SEMANTIC_ACTION_READ",
            Self::Create => "SEMANTIC_ACTION_CREATE",
            Self::Update => "SEMANTIC_ACTION_UPDATE",
            Self::Delete => "SEMANTIC_ACTION_DELETE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for SemanticAction {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "SEMANTIC_ACTION_UNSPECIFIED",
            "SEMANTIC_ACTION_READ",
            "SEMANTIC_ACTION_CREATE",
            "SEMANTIC_ACTION_UPDATE",
            "SEMANTIC_ACTION_DELETE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SemanticAction;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "SEMANTIC_ACTION_UNSPECIFIED" => Ok(SemanticAction::Unspecified),
                    "SEMANTIC_ACTION_READ" => Ok(SemanticAction::Read),
                    "SEMANTIC_ACTION_CREATE" => Ok(SemanticAction::Create),
                    "SEMANTIC_ACTION_UPDATE" => Ok(SemanticAction::Update),
                    "SEMANTIC_ACTION_DELETE" => Ok(SemanticAction::Delete),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for Semantics {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.action != 0 {
            len += 1;
        }
        if !self.resource_type.is_empty() {
            len += 1;
        }
        if self.version != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.Semantics", len)?;
        if self.action != 0 {
            let v = SemanticAction::try_from(self.action)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.action)))?;
            struct_ser.serialize_field("action", &v)?;
        }
        if !self.resource_type.is_empty() {
            struct_ser.serialize_field("resourceType", &self.resource_type)?;
        }
        if self.version != 0 {
            struct_ser.serialize_field("version", &self.version)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Semantics {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "action",
            "resource_type",
            "resourceType",
            "version",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Action,
            ResourceType,
            Version,
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
                            "action" => Ok(GeneratedField::Action),
                            "resourceType" | "resource_type" => Ok(GeneratedField::ResourceType),
                            "version" => Ok(GeneratedField::Version),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Semantics;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.Semantics")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Semantics, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut action__ = None;
                let mut resource_type__ = None;
                let mut version__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Action => {
                            if action__.is_some() {
                                return Err(serde::de::Error::duplicate_field("action"));
                            }
                            action__ = Some(map_.next_value::<SemanticAction>()? as i32);
                        }
                        GeneratedField::ResourceType => {
                            if resource_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("resourceType"));
                            }
                            resource_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Semantics {
                    action: action__.unwrap_or_default(),
                    resource_type: resource_type__.unwrap_or_default(),
                    version: version__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.Semantics", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SetupField {
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
        if !self.r#type.is_empty() {
            len += 1;
        }
        if !self.pattern.is_empty() {
            len += 1;
        }
        if self.secret {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.contracts.v1.SetupField", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.r#type.is_empty() {
            struct_ser.serialize_field("type", &self.r#type)?;
        }
        if !self.pattern.is_empty() {
            struct_ser.serialize_field("pattern", &self.pattern)?;
        }
        if self.secret {
            struct_ser.serialize_field("secret", &self.secret)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SetupField {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "type",
            "pattern",
            "secret",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Type,
            Pattern,
            Secret,
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
                            "type" => Ok(GeneratedField::Type),
                            "pattern" => Ok(GeneratedField::Pattern),
                            "secret" => Ok(GeneratedField::Secret),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SetupField;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.contracts.v1.SetupField")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SetupField, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut r#type__ = None;
                let mut pattern__ = None;
                let mut secret__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Pattern => {
                            if pattern__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pattern"));
                            }
                            pattern__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Secret => {
                            if secret__.is_some() {
                                return Err(serde::de::Error::duplicate_field("secret"));
                            }
                            secret__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SetupField {
                    name: name__.unwrap_or_default(),
                    r#type: r#type__.unwrap_or_default(),
                    pattern: pattern__.unwrap_or_default(),
                    secret: secret__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.contracts.v1.SetupField", FIELDS, GeneratedVisitor)
    }
}
