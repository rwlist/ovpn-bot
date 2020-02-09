# ovpn-bot
Simple telegram bot for automatic OpenVPN configuration

## Installation

TODO: one-liner bash script or single docker run.

```bash
docker run -d \
    --name ovpn-tg-bot \
    --volume /var/run/docker.sock:/var/run/docker.sock:ro \
    --env ADMIN_TELEGRAM_ID=123456789 \
    --env BOT_TOKEN=1231231231:AAAAAAAAABBBBCCCCCCCCCCCCCC \
    ovpnbot

docker logs -f ovpn-tg-bot
```

## Keep it as simple as that

Supported commands:
- [X] Initialize containers
- [X] Show status
- [X] Generate config
- [X] Remove everything
- [X] Show help
