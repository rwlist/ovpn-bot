# ovpn-bot
Simple telegram bot for automatic OpenVPN configuration

## Installation

```bash
docker run -d \
    --name ovpn-tg-bot \
    --volume /var/run/docker.sock:/var/run/docker.sock:ro \
    --env ADMIN_TELEGRAM_ID=123456789 \
    --env BOT_TOKEN=1231231231:AAAAAAAAABBBBCCCCCCCCCCCCCC \
    arthurwow/ovpnbot
```
or
```bash
docker run -d --name ovpn-tg-bot --volume /var/run/docker.sock:/var/run/docker.sock:ro --env ADMIN_TELEGRAM_ID=123456789 --env BOT_TOKEN=1231231231:AAAAAAAAABBBBCCCCCCCCCCCCCC arthurwow/ovpnbot
```

**ADMIN_TELEGRAM_ID** is a comma-separated list of the bot admins' Telegram IDs

### Proxy

If you running this bot in Russia, you might want to use proxy to run this bot. You can do this easily with HTTP_PROXY envvar, example:

```
docker run -d \
    --name ovpn-tg-bot \
    --volume /var/run/docker.sock:/var/run/docker.sock:ro \
    --env HTTP_PROXY=socks5://user:pass@proxy.example.com:1080 \
    --env ADMIN_TELEGRAM_ID=123456789 \
    --env BOT_TOKEN=1231231231:AAAAAAAAABBBBCCCCCCCCCCCCCC \
    arthurwow/ovpnbot
```

### docker-compose

TODO: docker-compose

## Commands

Supported commands:
- [X] Initialize containers
- [X] Show status
- [X] Generate config
- [X] Remove everything
- [X] Show help
