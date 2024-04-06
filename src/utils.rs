/// Merges two JSON objects recursively.
///
/// # Arguments
///
/// * `a` - The first JSON object.
/// * `b` - The second JSON object.
///
/// # Returns
///
/// The merged JSON object.
///
/// # Examples
///
/// ```
/// use serde_json::json;
/// use inferix::utils::merge_objects;
///
/// let a = json!({"person": {"name": "John", "age": 25}});
/// let b = json!({"person": {"age": 30, "city": "New York"}});
/// let expected = json!({"person": {"name": "John", "age": 30, "city": "New York"}});
/// assert_eq!(merge_objects(a, b), expected);
/// ```
pub fn merge_objects(a: serde_json::Value, b: serde_json::Value) -> serde_json::Value {
    let mut a = a.as_object().unwrap().clone();
    let b = b.as_object().unwrap();

    // Perform a deep merge
    for (key, value) in b.iter() {
        if a.contains_key(key) {
            if a[key].is_object() && value.is_object() {
                a.insert(
                    key.to_string(),
                    merge_objects(a[key].clone(), value.clone()),
                );
            } else {
                a.insert(key.to_string(), value.clone());
            }
        } else {
            a.insert(key.to_string(), value.clone());
        }
    }

    return serde_json::Value::Object(a);
}

/// Converts the provided datetime string to a Unix timestamp.
///
/// # Arguments
///
/// * `datetime` - The datetime string in RFC 3339 format.
///
/// # Returns
///
/// The Unix timestamp as an `i64` value.
///
/// # Examples
///
/// ```
/// use inferix::utils::convert_to_unix_timestamp;
///
/// // Test case 1: Valid datetime string
/// let datetime = "2022-01-01T00:00:00Z";
/// let expected = 1640995200;
/// assert_eq!(convert_to_unix_timestamp(datetime), expected);
///
/// // Test case 2: Invalid datetime string
/// let datetime = "2022-13-01T00:00:00Z";
/// // Assuming the function returns 0 for invalid datetime strings
/// let expected = 0;
/// assert_eq!(convert_to_unix_timestamp(datetime), expected);
/// ```
pub fn convert_to_unix_timestamp(datetime: &str) -> i64 {
    let dt = chrono::DateTime::parse_from_rfc3339(datetime);
    if let Ok(dt) = dt {
        return dt.timestamp();
    } else {
        // Return 0 if the datetime string is invalid
        return 0;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_merge_objects() {
        // Test case 1: Both objects are empty
        let a = serde_json::json!({});
        let b = serde_json::json!({});
        let expected = serde_json::json!({});
        assert_eq!(merge_objects(a, b), expected);

        // Test case 2: Object b is empty
        let a = serde_json::json!({"name": "John"});
        let b = serde_json::json!({});
        let expected = serde_json::json!({"name": "John"});
        assert_eq!(merge_objects(a, b), expected);

        // Test case 3: Object a is empty
        let a = serde_json::json!({});
        let b = serde_json::json!({"age": 25});
        let expected = serde_json::json!({"age": 25});
        assert_eq!(merge_objects(a, b), expected);

        // Test case 4: Both objects have overlapping keys
        let a = serde_json::json!({"name": "John", "age": 25});
        let b = serde_json::json!({"age": 30, "city": "New York"});
        let expected = serde_json::json!({"name": "John", "age": 30, "city": "New York"});
        assert_eq!(merge_objects(a, b), expected);

        // Test case 5: Nested objects
        let a = serde_json::json!({"person": {"name": "John", "age": 25}});
        let b = serde_json::json!({"person": {"age": 30, "city": "New York"}});
        let expected =
            serde_json::json!({"person": {"name": "John", "age": 30, "city": "New York"}});
        assert_eq!(merge_objects(a, b), expected);
    }

    #[test]
    fn test_convert_to_unix_timestamp() {
        // Test case 1: Valid datetime string
        let datetime = "2022-01-01T00:00:00Z";
        let expected = 1640995200;
        assert_eq!(convert_to_unix_timestamp(datetime), expected);

        // Test case 2: Another valid datetime string
        let datetime = "2022-12-31T23:59:59Z";
        let expected = 1672531199;
        assert_eq!(convert_to_unix_timestamp(datetime), expected);

        // Test case 3: Invalid datetime string
        let datetime = "2022-13-01T00:00:00Z";
        // Assuming the function returns 0 for invalid datetime strings
        let expected = 0;
        assert_eq!(convert_to_unix_timestamp(datetime), expected);
    }
}
