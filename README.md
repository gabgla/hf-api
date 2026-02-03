# HFAPI - Hellfall Card Search API

HF API is the API counterpart to [Hellfall](https://skeleton.club/hellfall/)

If you're unfamiliar, Hellfall in turn is an equivalent of Scryfall for [Hellscube](https://www.reddit.com/r/HellsCube): (taken from the sub)

> **HellsCube**
>
> A subreddit to create a fun and functional draft cube with cards that push the limits of Magic. We look for balanced and playable cards with strong comedy.

This is still very much WIP so expect many breaking changes. I'm looking to implement Scryfall's query language and response format next.

The rest of this README is AI generated until I stop being lazy and actually review the text:

---

A fast, serverless REST API for searching Hellscube/Hellfall cards. Built with Go and deployed on AWS Lambda with full-text search powered by Bleve.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Initial Setup](#initial-setup)
  - [Running Locally](#running-locally)
- [API Reference](#api-reference)
- [Development](#development)
  - [Running Tests](#running-tests)
  - [Code Generation](#code-generation)
- [Deployment](#deployment)
  - [Infrastructure Setup](#infrastructure-setup)
  - [Deploying Code Updates](#deploying-code-updates)
- [Configuration](#configuration)

---

## Overview

HFAPI provides a search endpoint for querying Hellscube cards using a syntax similar to [Scryfall](https://scryfall.com/docs/syntax). The API supports filtering by various card attributes like name, colors, mana value, type, creator, and more.

**Key Features:**
- Full-text search with Bleve
- Scryfall-like query syntax
- Serverless architecture (AWS Lambda + API Gateway)
- Infrastructure as Code (Terraform)

---

## Prerequisites

Before working with this project, ensure you have the following installed:

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.25+ | Primary language |
| [AWS CLI](https://aws.amazon.com/cli/) | v2 | AWS interactions |
| [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) | Latest | Local Lambda testing |
| [Terraform](https://developer.hashicorp.com/terraform/downloads) | 1.0+ | Infrastructure deployment |
| [Docker](https://docs.docker.com/get-docker/) | Latest | Required by SAM for local Lambda emulation |
| [curl](https://curl.se/) | Any | Downloading database |
| [zip](https://infozip.sourceforge.net/) | Any | Creating deployment packages |

**AWS Requirements (for deployment):**
- An AWS account with appropriate permissions
- A Route 53 hosted zone for your domain
- AWS credentials configured (`aws configure` or environment variables)

---

## Project Structure

```
hf-api/
├── src/
│   ├── cmd/
│   │   ├── lambda/          # Lambda entry point
│   │   ├── httpserver/      # Local HTTP server for development
│   │   ├── gendb/           # Database generator
│   │   └── codegens/        # Code generation utilities
│   ├── internal/
│   │   ├── app/api/         # API handlers and routing
│   │   └── data/            # Database and search index
│   ├── pkg/
│   │   ├── cards/           # Card data structures
│   │   └── hellfall/        # Hellfall database parser
│   └── utils/               # Utility functions
├── deployment/
│   ├── config/              # Terraform variable files and environment config
│   ├── scripts/             # Deployment scripts
│   └── terraform/           # Infrastructure as Code
│       ├── modules/
│       │   ├── dns/         # ACM certificate and DNS
│       │   └── vpc/         # VPC (optional)
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── index.bleve/             # Pre-built search index (generated)
├── template.yaml            # SAM template for local development
└── Makefile                 # Build and development commands
```

---

## Getting Started

### Initial Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd hf-api
   ```

2. **Install Go dependencies:**
   ```bash
   go mod download
   ```

3. **Set up the database and search index:**
   ```bash
   make setup
   ```
   This downloads the Hellscube card database and generates the Bleve search index.

### Running Locally

There are two ways to run the API locally:

#### Option 1: Simple HTTP Server (Recommended for development)
```bash
make run
# or
make run-http
```
The API will be available at `http://localhost:8080`

#### Option 2: SAM Local (Lambda emulation)
```bash
make run-lambda
```
This uses AWS SAM to emulate the Lambda environment locally. The API will be available at `http://localhost:3000`

**Test the API:**
```bash
# Health check
curl http://localhost:8080/v1/health

# Search for cards
curl "http://localhost:8080/v1/cards/search?q=creator:Leslie"
```

---

## API Reference

### Base URL
- Local: `http://localhost:8080/v1`
- Production: `https://hfapi.saguinus.net/v1`

### Endpoints

#### Health Check
```
GET /v1/health
```
Returns the API health status.

**Response:**
```json
{
  "status": "ok"
}
```

#### Card Search
```
GET /v1/cards/search?q=<query>
```
Search for cards using Scryfall-like syntax.

**Query Parameters:**
| Parameter | Required | Description |
|-----------|----------|-------------|
| `q` | Yes | Search query string |

**Supported Search Tokens:**

| Token | Aliases | Example |
|-------|---------|---------|
| `name` | `n` | `name:Fireball` |
| `colors` | `c`, `color` | `c:red` |
| `mv` | `cmc` | `mv:>3` |
| `mana` | `m` | `m:2RR` |
| `identity` | `id` | `id:gruul` |
| `type_line` | `type`, `t` | `t:creature` |
| `oracle` | `o` | `o:draw` |
| `flavor_text` | `ft`, `flavor` | `ft:ancient` |
| `power` | `pow` | `pow:>4` |
| `toughness` | `tou`, `tough` | `tou:<=2` |
| `set` | `s`, `edition`, `e` | `set:HFC` |
| `creator` | `author` | `creator:Leslie` |
| `tags` | `tag` | `tags:combo` |

**Operators:**
- `:` - Contains/matches
- `=` - Equals (numeric)
- `>`, `<`, `>=`, `<=` - Numeric comparisons

**Example Queries:**
```bash
# Cards by creator with MV > 4 that are artifacts
curl "https://hfapi.saguinus.net/v1/cards/search?q=creator:Leslie+mv:>4+t:Artifact"

# Red creatures with power > 3
curl "https://hfapi.saguinus.net/v1/cards/search?q=c:red+t:creature+pow:>3"
```

---

## Development

### Running Tests
```bash
make test
```

### Code Generation

Regenerate token aliases (if you modify `KnownTokens` in handlers.go):
```bash
make generate-token-aliases
```

Regenerate the database and search index:
```bash
make setup
```

### Build Commands

| Command | Description |
|---------|-------------|
| `make setup` | Download database and generate search index |
| `make run` | Run local HTTP server |
| `make run-lambda` | Run with SAM local |
| `make test` | Run tests |
| `make build-for-lambda` | Build Lambda deployment package |
| `make clean` | Remove build artifacts |

---

## Deployment

### Infrastructure Setup

The infrastructure is managed with Terraform and includes:
- S3 bucket for Lambda artifacts
- Lambda function
- API Gateway (HTTP API)
- ACM certificate
- Route 53 DNS records
- CloudWatch log groups

#### First-time Setup

1. **Create your configuration files:**
   ```bash
   # Copy and edit the example configs
   cp deployment/config/dev.tfvars.example deployment/config/dev.tfvars
   cp deployment/config/dev.backend.hcl.example deployment/config/dev.backend.hcl
   ```

2. **Edit the configuration** with your values:
   - `hosted_zone_name`: Your Route 53 hosted zone (e.g., `example.com`)
   - `api_subdomain`: Subdomain for the API (e.g., `hfapi`)
   - Backend config: S3 bucket for Terraform state

3. **Create the bootstrap Lambda package:**
   ```bash
   ./deployment/terraform/scripts/generate_bootstrap.sh
   ```

4. **Initialize and apply Terraform:**
   ```bash
   cd deployment/terraform

   # Initialize with backend config
   terraform init -backend-config=../config/live.backend.hcl

   # Review the plan
   terraform plan -var-file=../config/live.tfvars

   # Apply
   terraform apply -var-file=../config/live.tfvars
   ```

### Deploying Code Updates

After making code changes, deploy to Lambda:

```bash
./deployment/scripts/build_and_deploy.sh
```

This script:
1. Builds the Lambda binary
2. Creates a deployment package with the search index
3. Uploads to S3
4. Updates the Lambda function

---

## Configuration

### Terraform Variables

Key variables in `deployment/config/*.tfvars`:

| Variable | Description | Default |
|----------|-------------|---------|
| `aws_region` | AWS region | `us-east-1` |
| `namespace` | Resource name prefix | `hfapi` |
| `environment` | Environment name | - |
| `hosted_zone_name` | Route 53 zone | - |
| `api_subdomain` | API subdomain | `api` |
| `lambda_memory_size` | Lambda memory (MB) | `128` |
| `lambda_timeout` | Lambda timeout (seconds) | `10` |
| `cors_allow_origins` | CORS origins | `["*"]` |
| `cors_allow_methods` | CORS methods | `["GET", "OPTIONS"]` |
| `log_retention_days` | CloudWatch log retention | `7` |

### Environment Variables

Lambda environment variables (set in tfvars):

| Variable | Description |
|----------|-------------|
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) |

---

## License

MIT
