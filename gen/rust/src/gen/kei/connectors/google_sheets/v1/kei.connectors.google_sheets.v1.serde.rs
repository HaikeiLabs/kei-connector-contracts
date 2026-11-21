// @generated
impl serde::Serialize for AppendValuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if !self.range.is_empty() {
            len += 1;
        }
        if self.value_input_option.is_some() {
            len += 1;
        }
        if self.insert_data_option.is_some() {
            len += 1;
        }
        if !self.values.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.AppendValuesRequest", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if !self.range.is_empty() {
            struct_ser.serialize_field("range", &self.range)?;
        }
        if let Some(v) = self.value_input_option.as_ref() {
            struct_ser.serialize_field("valueInputOption", v)?;
        }
        if let Some(v) = self.insert_data_option.as_ref() {
            struct_ser.serialize_field("insertDataOption", v)?;
        }
        if !self.values.is_empty() {
            struct_ser.serialize_field("values", &self.values)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AppendValuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "range",
            "value_input_option",
            "valueInputOption",
            "insert_data_option",
            "insertDataOption",
            "values",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            Range,
            ValueInputOption,
            InsertDataOption,
            Values,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "range" => Ok(GeneratedField::Range),
                            "valueInputOption" | "value_input_option" => Ok(GeneratedField::ValueInputOption),
                            "insertDataOption" | "insert_data_option" => Ok(GeneratedField::InsertDataOption),
                            "values" => Ok(GeneratedField::Values),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AppendValuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.AppendValuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AppendValuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut range__ = None;
                let mut value_input_option__ = None;
                let mut insert_data_option__ = None;
                let mut values__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ValueInputOption => {
                            if value_input_option__.is_some() {
                                return Err(serde::de::Error::duplicate_field("valueInputOption"));
                            }
                            value_input_option__ = map_.next_value()?;
                        }
                        GeneratedField::InsertDataOption => {
                            if insert_data_option__.is_some() {
                                return Err(serde::de::Error::duplicate_field("insertDataOption"));
                            }
                            insert_data_option__ = map_.next_value()?;
                        }
                        GeneratedField::Values => {
                            if values__.is_some() {
                                return Err(serde::de::Error::duplicate_field("values"));
                            }
                            values__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(AppendValuesRequest {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    range: range__.unwrap_or_default(),
                    value_input_option: value_input_option__,
                    insert_data_option: insert_data_option__,
                    values: values__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.AppendValuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for AppendValuesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if self.table_range.is_some() {
            len += 1;
        }
        if self.updates.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.AppendValuesResponse", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if let Some(v) = self.table_range.as_ref() {
            struct_ser.serialize_field("tableRange", v)?;
        }
        if let Some(v) = self.updates.as_ref() {
            struct_ser.serialize_field("updates", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AppendValuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "table_range",
            "tableRange",
            "updates",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            TableRange,
            Updates,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "tableRange" | "table_range" => Ok(GeneratedField::TableRange),
                            "updates" => Ok(GeneratedField::Updates),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AppendValuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.AppendValuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AppendValuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut table_range__ = None;
                let mut updates__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TableRange => {
                            if table_range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tableRange"));
                            }
                            table_range__ = map_.next_value()?;
                        }
                        GeneratedField::Updates => {
                            if updates__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updates"));
                            }
                            updates__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(AppendValuesResponse {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    table_range: table_range__,
                    updates: updates__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.AppendValuesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetSpreadsheetRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.GetSpreadsheetRequest", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetSpreadsheetRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetSpreadsheetRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.GetSpreadsheetRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetSpreadsheetRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetSpreadsheetRequest {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.GetSpreadsheetRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetSpreadsheetResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if self.properties.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.GetSpreadsheetResponse", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if let Some(v) = self.properties.as_ref() {
            struct_ser.serialize_field("properties", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetSpreadsheetResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "properties",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            Properties,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "properties" => Ok(GeneratedField::Properties),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetSpreadsheetResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.GetSpreadsheetResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetSpreadsheetResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut properties__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Properties => {
                            if properties__.is_some() {
                                return Err(serde::de::Error::duplicate_field("properties"));
                            }
                            properties__ = map_.next_value()?;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetSpreadsheetResponse {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    properties: properties__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.GetSpreadsheetResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if !self.range.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.GetValuesRequest", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if !self.range.is_empty() {
            struct_ser.serialize_field("range", &self.range)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "range",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            Range,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "range" => Ok(GeneratedField::Range),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.GetValuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetValuesRequest {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    range: range__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.GetValuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValuesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.range.is_empty() {
            len += 1;
        }
        if !self.values.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.GetValuesResponse", len)?;
        if !self.range.is_empty() {
            struct_ser.serialize_field("range", &self.range)?;
        }
        if !self.values.is_empty() {
            struct_ser.serialize_field("values", &self.values)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "range",
            "values",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Range,
            Values,
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
                            "range" => Ok(GeneratedField::Range),
                            "values" => Ok(GeneratedField::Values),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.GetValuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut range__ = None;
                let mut values__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Values => {
                            if values__.is_some() {
                                return Err(serde::de::Error::duplicate_field("values"));
                            }
                            values__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(GetValuesResponse {
                    range: range__.unwrap_or_default(),
                    values: values__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.GetValuesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SpreadsheetProperties {
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
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.SpreadsheetProperties", len)?;
        if !self.title.is_empty() {
            struct_ser.serialize_field("title", &self.title)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SpreadsheetProperties {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "title",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Title,
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
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SpreadsheetProperties;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.SpreadsheetProperties")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SpreadsheetProperties, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut title__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Title => {
                            if title__.is_some() {
                                return Err(serde::de::Error::duplicate_field("title"));
                            }
                            title__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(SpreadsheetProperties {
                    title: title__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.SpreadsheetProperties", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateValuesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if !self.range.is_empty() {
            len += 1;
        }
        if !self.values.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.UpdateValuesRequest", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if !self.range.is_empty() {
            struct_ser.serialize_field("range", &self.range)?;
        }
        if !self.values.is_empty() {
            struct_ser.serialize_field("values", &self.values)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateValuesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "range",
            "values",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            Range,
            Values,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "range" => Ok(GeneratedField::Range),
                            "values" => Ok(GeneratedField::Values),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdateValuesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.UpdateValuesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateValuesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut range__ = None;
                let mut values__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Values => {
                            if values__.is_some() {
                                return Err(serde::de::Error::duplicate_field("values"));
                            }
                            values__ = Some(map_.next_value()?);
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateValuesRequest {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    range: range__.unwrap_or_default(),
                    values: values__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.UpdateValuesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UpdateValuesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.spreadsheet_id.is_empty() {
            len += 1;
        }
        if !self.updated_range.is_empty() {
            len += 1;
        }
        if self.updated_rows.is_some() {
            len += 1;
        }
        if self.updated_columns.is_some() {
            len += 1;
        }
        if self.updated_cells.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("kei.connectors.google_sheets.v1.UpdateValuesResponse", len)?;
        if !self.spreadsheet_id.is_empty() {
            struct_ser.serialize_field("spreadsheetId", &self.spreadsheet_id)?;
        }
        if !self.updated_range.is_empty() {
            struct_ser.serialize_field("updatedRange", &self.updated_range)?;
        }
        if let Some(v) = self.updated_rows.as_ref() {
            struct_ser.serialize_field("updatedRows", v)?;
        }
        if let Some(v) = self.updated_columns.as_ref() {
            struct_ser.serialize_field("updatedColumns", v)?;
        }
        if let Some(v) = self.updated_cells.as_ref() {
            struct_ser.serialize_field("updatedCells", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UpdateValuesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "spreadsheet_id",
            "spreadsheetId",
            "updated_range",
            "updatedRange",
            "updated_rows",
            "updatedRows",
            "updated_columns",
            "updatedColumns",
            "updated_cells",
            "updatedCells",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SpreadsheetId,
            UpdatedRange,
            UpdatedRows,
            UpdatedColumns,
            UpdatedCells,
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
                            "spreadsheetId" | "spreadsheet_id" => Ok(GeneratedField::SpreadsheetId),
                            "updatedRange" | "updated_range" => Ok(GeneratedField::UpdatedRange),
                            "updatedRows" | "updated_rows" => Ok(GeneratedField::UpdatedRows),
                            "updatedColumns" | "updated_columns" => Ok(GeneratedField::UpdatedColumns),
                            "updatedCells" | "updated_cells" => Ok(GeneratedField::UpdatedCells),
                            _ => Ok(GeneratedField::__SkipField__),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UpdateValuesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct kei.connectors.google_sheets.v1.UpdateValuesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UpdateValuesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut spreadsheet_id__ = None;
                let mut updated_range__ = None;
                let mut updated_rows__ = None;
                let mut updated_columns__ = None;
                let mut updated_cells__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SpreadsheetId => {
                            if spreadsheet_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("spreadsheetId"));
                            }
                            spreadsheet_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::UpdatedRange => {
                            if updated_range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedRange"));
                            }
                            updated_range__ = Some(map_.next_value()?);
                        }
                        GeneratedField::UpdatedRows => {
                            if updated_rows__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedRows"));
                            }
                            updated_rows__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::UpdatedColumns => {
                            if updated_columns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedColumns"));
                            }
                            updated_columns__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::UpdatedCells => {
                            if updated_cells__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedCells"));
                            }
                            updated_cells__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::__SkipField__ => {
                            let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                        }
                    }
                }
                Ok(UpdateValuesResponse {
                    spreadsheet_id: spreadsheet_id__.unwrap_or_default(),
                    updated_range: updated_range__.unwrap_or_default(),
                    updated_rows: updated_rows__,
                    updated_columns: updated_columns__,
                    updated_cells: updated_cells__,
                })
            }
        }
        deserializer.deserialize_struct("kei.connectors.google_sheets.v1.UpdateValuesResponse", FIELDS, GeneratedVisitor)
    }
}
