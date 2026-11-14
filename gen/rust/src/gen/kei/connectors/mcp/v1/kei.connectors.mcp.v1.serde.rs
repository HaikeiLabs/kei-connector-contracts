// @generated
impl serde::Serialize for ToolCall {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.server.is_empty() {
            len += 1;
        }
        if !self.tool.is_empty() {
            len += 1;
        }
        if self.arguments.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.mcp.v1.ToolCall", len)?;
        if !self.server.is_empty() {
            struct_ser.serialize_field("server", &self.server)?;
        }
        if !self.tool.is_empty() {
            struct_ser.serialize_field("tool", &self.tool)?;
        }
        if let Some(v) = self.arguments.as_ref() {
            struct_ser.serialize_field("arguments", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ToolCall {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "server",
            "tool",
            "arguments",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Server,
            Tool,
            Arguments,
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
                            "server" => Ok(GeneratedField::Server),
                            "tool" => Ok(GeneratedField::Tool),
                            "arguments" => Ok(GeneratedField::Arguments),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ToolCall;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.mcp.v1.ToolCall")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ToolCall, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut server__ = None;
                let mut tool__ = None;
                let mut arguments__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Server => {
                            if server__.is_some() {
                                return Err(serde::de::Error::duplicate_field("server"));
                            }
                            server__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Tool => {
                            if tool__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tool"));
                            }
                            tool__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Arguments => {
                            if arguments__.is_some() {
                                return Err(serde::de::Error::duplicate_field("arguments"));
                            }
                            arguments__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ToolCall {
                    server: server__.unwrap_or_default(),
                    tool: tool__.unwrap_or_default(),
                    arguments: arguments__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.mcp.v1.ToolCall", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ToolResult {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.content.is_empty() {
            len += 1;
        }
        if self.is_error {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.mcp.v1.ToolResult", len)?;
        if !self.content.is_empty() {
            struct_ser.serialize_field("content", &self.content)?;
        }
        if self.is_error {
            struct_ser.serialize_field("isError", &self.is_error)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ToolResult {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "content",
            "is_error",
            "isError",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Content,
            IsError,
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
                            "isError" | "is_error" => Ok(GeneratedField::IsError),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ToolResult;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.mcp.v1.ToolResult")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ToolResult, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut content__ = None;
                let mut is_error__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Content => {
                            if content__.is_some() {
                                return Err(serde::de::Error::duplicate_field("content"));
                            }
                            content__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IsError => {
                            if is_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isError"));
                            }
                            is_error__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(ToolResult {
                    content: content__.unwrap_or_default(),
                    is_error: is_error__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.mcp.v1.ToolResult", FIELDS, GeneratedVisitor)
    }
}
