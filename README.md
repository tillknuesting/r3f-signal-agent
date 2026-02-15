# r3f-signal-agent

Trend collection tool for gathering trending topics from tech sources (HackerNews, etc.) and preparing them for blog content generation.

## Features

- **Multi-source collection**: HackerNews (front page + newest), extensible to other sources
- **Hexagonal architecture**: Swappable components (storage, collectors, LLM agents)
- **REST API**: Control collection, browse trends, manage sources
- **Fancy TUI**: Terminal UI built with Bubble Tea
- **LLM integration**: GLM-5 (z.ai) for summarization and blog topic suggestions
- **Markdown storage**: Human-readable trend files
- **Profile system**: Switch between different content domains (tech, finance, etc.)

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    DRIVING ADAPTERS                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  HTTP API   │  │    TUI      │  │       CLI           │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      PORTS (Interfaces)                     │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Services: Collector, Trend, Source, Profile, Agent │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                      DOMAIN (Entities)                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Trend, Source, Profile, Collection                 │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                    DRIVEN ADAPTERS                          │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌────────────┐  │
│  │ Markdown  │ │   HTTP    │ │  Chrome   │ │   GLM-5    │  │
│  │ Storage   │ │ Collector │ │ Collector │ │   Agent    │  │
│  └───────────┘ └───────────┘ └───────────┘ └────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Clone
git clone https://github.com/tillknuesting/r3f-signal-agent.git
cd r3f-signal-agent

# Set API key
export ZAI_API_KEY=your_api_key

# Run server
make run-server

# Collect trends
curl -X POST http://localhost:8080/api/v1/collect

# View trends
curl http://localhost:8080/api/v1/trends

# Run TUI
make run-tui
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/trends` | List all trends |
| GET | `/api/v1/trends?date=2026-02-15` | Trends by date |
| GET | `/api/v1/trends/search?q=query` | Search trends |
| POST | `/api/v1/trends/:id/star` | Star trend |
| POST | `/api/v1/collect` | Trigger collection |
| GET | `/api/v1/sources` | List sources |
| POST | `/api/v1/agent/summarize` | Summarize with LLM |

## Configuration

```yaml
# config/config.yaml
active_profile: "tech"

server:
  host: "0.0.0.0"
  port: 8080

scheduler:
  enabled: true
  interval: 24h

llm:
  provider: "zai"
  model: "glm-5"
  api_key: "${ZAI_API_KEY}"

storage:
  type: "markdown"
  base_path: "./data/profiles"
```

## Adding Sources

Sources are defined in YAML files under `config/sources/`:

```yaml
sources:
  - id: "hackernews-newest"
    name: "Hacker News (Newest)"
    type: "http"
    enabled: true
    config:
      url: "https://hacker-news.firebaseio.com/v0/newstories.json"
      item_url: "https://hacker-news.firebaseio.com/v0/item/{id}.json"
      limit: 30
    display:
      icon: "🟠"
      color: "#FF6600"
```

## Data Format

Trends are stored as Markdown files with embedded JSON:

```markdown
---
date: 2026-02-15
count: 47
---

```json
[
  {
    "id": "hackernews-newest-12345",
    "title": "New Rust 2.0 Features",
    "url": "https://...",
    "score": 342,
    "source": "Hacker News (Newest)"
  }
]
```
```

## License

MIT
