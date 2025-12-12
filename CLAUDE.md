# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**transaction-mapper** is a financial transaction processing tool that converts bank transaction export files into accounting software import templates. It implements a clean plugin-based architecture supporting both CLI and web interfaces.

## Architecture

The application follows a **plugin architecture** with clear separation of concerns:

- **Bank Plugins** (`pkg/bank/`): Parse different bank export formats (ICBC, CMB, SPDB)
- **Consumer Plugins** (`pkg/consumer/`): Export to different accounting apps (Bluecoins, Qianji)
- **Transaction Registry** (`pkg/transaction/`): Standardized intermediate transaction format
- **Configuration Engine** (`pkg/config/`): Rule-based categorization and account mapping
- **CLI Interface** (`cmd/`): Command-line processing using Cobra
- **Web Interface** (`web/`): Vue.js frontend with Gin backend

## Development Commands

### Backend (Go)

```bash
# Run CLI tool with bank file conversion
go run main.go -i input.csv -b icbc -a bluecoins -t "工商银行" -z "储蓄卡"

# Start web server
go run main.go serve

# Build Docker image locally (ARM64)
make build-local-image
```

### Frontend (Vue.js)

```bash
cd web/
npm install          # Install dependencies
npm run dev         # Development server
npm run build       # Production build
npm run preview     # Preview production build
```

## Key Configuration Files

- **`config.yaml`**: Transaction categorization rules, account transfer mappings, and keyword-based classification
- **`go.mod`**: Go 1.21.5 with dependencies for Gin, Cobra, gocsv, YAML handling, and samber/lo utilities
- **`web/package.json`**: Vue 3 + Vite frontend with Axios for API calls

## Plugin Architecture Patterns

### Bank Plugin Implementation
Implement the `bank.Plugin` interface:
```go
type Plugin interface {
    Name() string
    PreProcess(data []byte) (string, error)
    Parse(data string) ([]transaction.Transaction, error)
}
```

### Consumer Plugin Implementation
Implement the `consumer.Plugin` interface:
```go
type Plugin interface {
    Name() string
    Transform(transactions []transaction.Transaction, info transaction.AccountInfo) (interface{}, error)
}
```

### Registry Pattern
Both bank and consumer plugins use centralized registration:
- Register in `init()` functions using `Registry.register()`
- Type-safe plugin discovery and management
- Dependency injection support (consumers receive config)

## Configuration System

The `config.yaml` uses sophisticated keyword-based categorization:

- **Dual-Level Categories**: Support for hierarchical categories using `|` separator (e.g., "食|一日三餐")
- **Transfer Detection**: Automatic identification of account transfers with target account mapping
- **Keyword Mapping**: Efficient pre-built keyword-to-category maps for fast lookup
- **Skip Rules**: Define transaction patterns to ignore during processing
- **Default Categories**: Fallback categories for income/expense when no rules match

### Configuration Processing
- Pre-built lookup tables for efficient categorization
- Support for both single and dual-level categorization modes
- Transfer detection with configurable account-to-account mapping

## Web API Endpoints

- `GET /api/` - Health check
- `GET /api/version` - Version information
- `GET /api/banks` - List available bank parsers
- `GET /api/apps` - List available accounting app exporters
- `POST /api/transform` - Transform uploaded file with specified bank and app

## Transaction Processing Flow

1. **Input**: Bank export files (CSV/JSON from different banks)
2. **Pre-processing**: Bank-specific data cleaning and normalization
3. **Parsing**: Convert to standardized Transaction struct (`pkg/transaction/type.go`)
4. **Transformation**: Apply categorization rules from `config.yaml`
5. **Export**: Generate accounting app import format
6. **Output**: CSV file ready for import into target accounting app

## Key Design Patterns

### Provider Pattern
- Transaction processing uses provider interface for flexibility
- Banks implement `transaction.Provider` interface for data extraction
- Supports different data sources through consistent interface

### Configuration-Driven Architecture
- All categorization business logic externalized in YAML
- No hardcoded business rules in plugins
- Easy to add new rules without code changes

### Error Handling & Graceful Degradation
- Comprehensive error handling at each processing stage
- Fallback categories for unmatched transactions
- Detailed logging for debugging

## Adding New Features

### New Bank Support
1. Implement `bank.Plugin` interface in `pkg/bank/`
2. Handle bank-specific CSV/JSON format in `PreProcess` and `Parse`
3. Register plugin in `init()` function using `Registry.register()`

### New Accounting App Support
1. Implement `consumer.Plugin` interface in `pkg/consumer/`
2. Create mapping logic for target app format
3. Register plugin with config dependency injection
4. Update `config.yaml` with new categorization rules

### Code Quality Conventions
- **Interface-Based Development**: Dependencies through interfaces, not concrete types
- **Functional Programming**: Extensive use of samber/lo utilities for data transformation
- **Single Responsibility**: Each package has clear, focused purpose
- **Clean Separation**: Bank parsing, categorization, and export logic are completely separate