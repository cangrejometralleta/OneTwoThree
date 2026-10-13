# Configuration Examples

Optional format examples for the rules in [Constants](../rules/constants.md). Choose one format per configuration set.

## JSON

`constants/enrollment.json`:

```json
{
  "enrollmentStates": ["pending", "accepted", "rejected"]
}
```

`config/enrollment.production.json`:

```json
{
  "serviceUrl": "https://enrollment.example.com",
  "requestTimeoutSeconds": 15
}
```

## TOML

`constants/enrollment.toml`:

```toml
enrollmentStates = ["pending", "accepted", "rejected"]
```

`config/enrollment.production.toml`:

```toml
serviceUrl = "https://enrollment.example.com"
requestTimeoutSeconds = 15
```

Defaults and development files use the same keys in the selected format. The YAML examples remain in the rule.
