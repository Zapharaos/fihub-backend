# Translation Package Demo

This demo folder contains practical examples of how to use the translation package, including initialization,
translation methods, and sample configuration files for all supported formats.

## Files in this demo:

- `main.go` - Complete example showing initialization and usage
- `config/` - Sample translation files in all supported formats

## Running the demo:

```bash
# From the example directory
go run main.go
```

## Translation files:

The `config/` directory contains sample translation files in all supported formats:

- `messages.en.toml` - English translations (TOML format)
- `messages.fr.json` - French translations (JSON format)
- `messages.es.yaml` - Spanish translations (YAML format)
- `messages.de.yml` - German translations (YML format)

Each file demonstrates different features like pluralization, template variables, and nested structures.
