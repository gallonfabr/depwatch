# depwatch

A Python daemon that monitors dependency changelogs and sends digests to Slack or email.

## Installation

```bash
pip install depwatch
```

Or install from source:

```bash
git clone https://github.com/yourname/depwatch
cd depwatch
pip install .
```

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
| `notify.email.from` | Sender email address | — |
| `notify.email.smtp_host` | SMTP server hostname | — |
| `notify.email.smtp_port` | SMTP server port | `587` |

## License

MIT © yourname
