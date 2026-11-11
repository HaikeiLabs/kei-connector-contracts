// @generated
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
