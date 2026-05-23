# depwatch

A Go daemon that monitors dependency changelogs and sends digests to Slack or email.

## Installation

```bash
go install github.com/yourname/depwatch@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourname/depwatch/releases).

## Usage

Create a configuration file `depwatch.yaml`:

```yaml
dependencies:
  - name: requests
    registry: pypi
  - name: express
    registry: npm

notify:
  slack:
    webhook_url: "https://hooks.slack.com/services/your/webhook/url"
  email:
    to: "team@example.com"

interval: "24h"
```

Then start the daemon:

```bash
depwatch start --config depwatch.yaml
```

Depwatch will poll each dependency's changelog on the configured interval and send a digest summarizing new releases to your configured channels.

### One-off check

```bash
depwatch check --config depwatch.yaml
```

## Configuration Options

| Key | Description | Default |
|-----|-------------|---------|
| `interval` | How often to check for updates | `24h` |
| `notify.slack.webhook_url` | Slack incoming webhook URL | — |
| `notify.email.to` | Recipient email address | — |

## License

MIT © yourname