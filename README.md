# MobWiz

> A batteries-included CLI that scaffolds production-ready mobile modules for Android, Flutter, and iOS using consistent architectural patterns.

MobWiz streamlines the tedious parts of module creation: it asks a few high-level questions (or accepts flags), picks the right template, and writes all the boilerplate for data, domain, and presentation layers. The result is a predictable module tree that mirrors the architecture standards your team already uses.

## Highlights

- **Multi-platform scaffolding**: Android (Kotlin), Flutter (Dart), and iOS (Swift) out of the box.
- **Architecture-aware**: Clean Architecture, MVVM, and BLoC templates with matching directory layouts.
- **Interactive or scripted**: Run `mobwiz create` for prompts or pass flags for CI/CD usage.
- **Smart templating**: Uses Go templates plus utility funcs (`pascalCase`, `snakeCase`, etc.) to keep files idiomatic.
- **Extensible**: Customize `templates/templates.yaml` or add your own `.tmpl` files to expand coverage.

## Requirements

- Go 1.24 or later (`go env GOVERSION` to confirm)
- Git (to clone the repo)
- macOS/Linux/Windows shell

## Installation

### Install globally (recommended)

```bash
go install github.com/chingiz/mobwiz@latest
```

This will place the `mobwiz` binary in your `GOBIN` (usually `$(go env GOPATH)/bin` or `$GOBIN` if set), so make sure that directory is on your `PATH`.

### Install from source

```bash
git clone https://github.com/chingiz/mobwiz.git
cd mobwiz
go build -o mobwiz ./...
```

After building, add the binary to your `PATH` or run it from the repo root via `./mobwiz`.

## Usage

### Interactive flow

```bash
./mobwiz create
```

You’ll be guided through module name, platform, architecture, package identifiers, and optional extras (tests, DI, networking, local DB).

### Flag-driven flow

```bash
./mobwiz create \
  --name Profile \
  --platform "Flutter" \
  --architecture "bloc" \
  --package com.example.profile
```

Flags make MobWiz deterministic—perfect for CI pipelines or scaffolding scripts. Omit the `--package` flag for non-Android modules.

## Template Catalog

| Platform         | Architectures | Layers generated                                                                                                     |
| ---------------- | ------------- | -------------------------------------------------------------------------------------------------------------------- |
| Android (Kotlin) | `mvvm`        | Data (models, DAO, API, repository impl), Domain (models, repository, use cases), Presentation (ViewModel, Fragment) |
| Flutter (Dart)   | `bloc`        | Domain entities/repositories/use cases, data repositories, presentation BLoC trio + page                             |
| iOS (Swift)      | `mvvm`        | Domain entities/repositories/use cases, Data DTO/service/storage/repo impl, Presentation view/viewmodel/coordinator  |

Templates live in `templates/<platform>/<architecture>/...` and are referenced through `templates/templates.yaml`. Add new template files and update the YAML to expand MobWiz’ reach.

## Generated Project Structure

Here is an example of what MobWiz generates for a **Flutter** module using **Clean Architecture + BLoC**:

```text
lib/
├── data/
│   └── repositories/
│       └── profile_repository_impl.dart
├── domain/
│   ├── entities/
│   │   └── profile.dart
│   ├── repositories/
│   │   └── profile_repository.dart
│   └── usecases/
│       └── get_profile_usecase.dart
└── presentation/
    ├── bloc/
    │   ├── profile_bloc.dart
    │   ├── profile_event.dart
    │   └── profile_state.dart
    └── pages/
        └── profile_page.dart
```

This predictable structure ensures that every team member knows exactly where to find logic, UI, and data handling code.

## Configuration

MobWiz uses `templates/templates.yaml` to map your choices (Platform + Architecture) to specific template files.

```yaml
- platform: "Flutter"
  architecture: "Clean Architecture + BLoC"
  templates:
    - template: "flutter/bloc/domain/entity.dart.tmpl"
      output: "lib/domain/entities/{{snakeCase .Name}}.dart"
    # ... other templates
```

- **`template`**: Path to the source `.tmpl` file relative to the `templates/` directory.
- **`output`**: Destination path for the generated file. You can use Go template syntax and helper functions like `{{snakeCase .Name}}` or `{{pascalCase .Name}}` to dynamically name files.

## Troubleshooting

### Common Issues

- **`command not found: mobwiz`**: Ensure your Go binary directory is in your `PATH`.
  - Run `go env GOPATH` to find your workspace.
  - Add `$GOPATH/bin` to your shell config (e.g., `~/.zshrc` or `~/.bashrc`).
- **Templates not found**: If you installed from source, make sure you run `mobwiz` from the root of the repository or that the `templates/` directory is available relative to the binary if you moved it. (Future versions will embed templates into the binary).

## Development

```bash
# Run tests (add more as features expand)
go test ./...

# Verify lint/tooling (example)
golangci-lint run ./...
```

When editing templates, keep placeholders in Go template syntax (`{{ }}`) and leverage helper funcs defined in the templating engine (see `internal/generator/engine.go`).

## Customization Tips

- Duplicate an existing platform folder in `templates/` to add new architectural patterns.
- Update `templates/templates.yaml` with new `path` + `template` entries. Paths support helper funcs like `{{pascalCase .Name}}`.
- Extend the `Config` struct (in `internal/config`) and prompt flow to add richer options such as networking stacks or dependency injection frameworks.

## Roadmap Ideas

- Template packs for Jetpack Compose, SwiftUI, and Flutter Riverpod.
- Additional flags for opting in/out of tests, DI, or persistence.
- Exportable JSON/YAML of module definitions for documentation automation.
- Embed templates into the binary for easier distribution.

## Support

Issues and feature requests are welcome. Open a ticket or start a discussion to share ideas that could make MobWiz even more helpful.
