// @generated
impl serde::Serialize for Column {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.database.v1.Column", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.r#type.is_empty() {
            struct_ser.serialize_field("type", &self.r#type)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Column {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "type",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
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
                            "name" => Ok(GeneratedField::Name),
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
            type Value = Column;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.database.v1.Column")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Column, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut r#type__ = None;
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
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(Column {
                    name: name__.unwrap_or_default(),
                    r#type: r#type__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.database.v1.Column", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Dialect {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "DIALECT_UNSPECIFIED",
            Self::Postgres => "DIALECT_POSTGRES",
            Self::Mysql => "DIALECT_MYSQL",
            Self::Mssql => "DIALECT_MSSQL",
            Self::Snowflake => "DIALECT_SNOWFLAKE",
            Self::Supabase => "DIALECT_SUPABASE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for Dialect {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "DIALECT_UNSPECIFIED",
            "DIALECT_POSTGRES",
            "DIALECT_MYSQL",
            "DIALECT_MSSQL",
            "DIALECT_SNOWFLAKE",
            "DIALECT_SUPABASE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Dialect;

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
                    "DIALECT_UNSPECIFIED" => Ok(Dialect::Unspecified),
                    "DIALECT_POSTGRES" => Ok(Dialect::Postgres),
                    "DIALECT_MYSQL" => Ok(Dialect::Mysql),
                    "DIALECT_MSSQL" => Ok(Dialect::Mssql),
                    "DIALECT_SNOWFLAKE" => Ok(Dialect::Snowflake),
                    "DIALECT_SUPABASE" => Ok(Dialect::Supabase),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for RunQueryRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.dialect != 0 {
            len += 1;
        }
        if !self.sql.is_empty() {
            len += 1;
        }
        if !self.params.is_empty() {
            len += 1;
        }
        if self.max_rows.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.database.v1.RunQueryRequest", len)?;
        if self.dialect != 0 {
            let v = Dialect::try_from(self.dialect)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.dialect)))?;
            struct_ser.serialize_field("dialect", &v)?;
        }
        if !self.sql.is_empty() {
            struct_ser.serialize_field("sql", &self.sql)?;
        }
        if !self.params.is_empty() {
            struct_ser.serialize_field("params", &self.params)?;
        }
        if let Some(v) = self.max_rows.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("maxRows", ToString::to_string(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RunQueryRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "dialect",
            "sql",
            "params",
            "max_rows",
            "maxRows",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Dialect,
            Sql,
            Params,
            MaxRows,
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
                            "dialect" => Ok(GeneratedField::Dialect),
                            "sql" => Ok(GeneratedField::Sql),
                            "params" => Ok(GeneratedField::Params),
                            "maxRows" | "max_rows" => Ok(GeneratedField::MaxRows),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RunQueryRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.database.v1.RunQueryRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RunQueryRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut dialect__ = None;
                let mut sql__ = None;
                let mut params__ = None;
                let mut max_rows__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Dialect => {
                            if dialect__.is_some() {
                                return Err(serde::de::Error::duplicate_field("dialect"));
                            }
                            dialect__ = Some(map_.next_value::<Dialect>()? as i32);
                        }
                        GeneratedField::Sql => {
                            if sql__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sql"));
                            }
                            sql__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Params => {
                            if params__.is_some() {
                                return Err(serde::de::Error::duplicate_field("params"));
                            }
                            params__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MaxRows => {
                            if max_rows__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxRows"));
                            }
                            max_rows__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(RunQueryRequest {
                    dialect: dialect__.unwrap_or_default(),
                    sql: sql__.unwrap_or_default(),
                    params: params__.unwrap_or_default(),
                    max_rows: max_rows__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.database.v1.RunQueryRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RunQueryResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.columns.is_empty() {
            len += 1;
        }
        if !self.rows.is_empty() {
            len += 1;
        }
        if self.truncated {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.database.v1.RunQueryResponse", len)?;
        if !self.columns.is_empty() {
            struct_ser.serialize_field("columns", &self.columns)?;
        }
        if !self.rows.is_empty() {
            struct_ser.serialize_field("rows", &self.rows)?;
        }
        if self.truncated {
            struct_ser.serialize_field("truncated", &self.truncated)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RunQueryResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "columns",
            "rows",
            "truncated",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Columns,
            Rows,
            Truncated,
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
                            "columns" => Ok(GeneratedField::Columns),
                            "rows" => Ok(GeneratedField::Rows),
                            "truncated" => Ok(GeneratedField::Truncated),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RunQueryResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.database.v1.RunQueryResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RunQueryResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut columns__ = None;
                let mut rows__ = None;
                let mut truncated__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Columns => {
                            if columns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("columns"));
                            }
                            columns__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Rows => {
                            if rows__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rows"));
                            }
                            rows__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Truncated => {
                            if truncated__.is_some() {
                                return Err(serde::de::Error::duplicate_field("truncated"));
                            }
                            truncated__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(RunQueryResponse {
                    columns: columns__.unwrap_or_default(),
                    rows: rows__.unwrap_or_default(),
                    truncated: truncated__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.database.v1.RunQueryResponse", FIELDS, GeneratedVisitor)
    }
}
