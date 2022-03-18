# ovpn-bot
Simple telegram bot for automatic OpenVPN configuration

## How to install

### Create configuration

Create `docker-compose.yml` with the following content:

```yml
version: '3.7'
services:
  bot:
    image: arthurwow/ovpnbot
    env_file:
      - .env
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped
```

Also create `.env` file with configuration:

```env
ADMIN_TELEGRAM_ID=
BOT_TOKEN=
```

ADMIN_TELEGRAM_ID is your Telegram ID, which will be used to manage users. Use https://t.me/userinfobot to get it. Bot will reply with a message with your ID:

```
@ivanonov
Id: 123456789 <==== USE THIS ID TO WRITE IN .env FILE
First: Ivan
Last: Ivanov
Lang: en
```

BOT_TOKEN is a Telegram token issued by https://t.me/BotFather. You should talk to this bot to create a new bot and get a token for it.

So, configured `.env` file should look like this:

```env
ADMIN_TELEGRAM_ID=123456789
BOT_TOKEN=1231231231:AAAAAAAAABBBBCCCCCCCCCCCCCC
```

### Start bot

When configuration is written to `.env` file, start the bot:

```bash
# Start the bot in the background
docker-compose up -d

# Now bot should be working, check if it is:
docker-compose ps -a           

# If output is like this, State is Up, then bot is running
# 
#      Name        Command   State   Ports
# ----------------------------------------
# ovpn-bot_bot_1   ./app     Up       

# Get public IP address of the instance, to use it later for 
# VPN configuration. Remember this IP address for bot init!
curl ifconfig.me
```

Now it's time to setup VPN via bot.

### Use bot

First of all, write something to your bot. It should reply with help message. If not, something was wrong on the previous step.

If bot is working, you can init your VPN server. To do that, write init command to the bot. You must replace 0.0.0.0 with the actual IP address of your VPN server, which can be obtained by running `curl ifconfig.me` on the same instance.
```
/init tcp://0.0.0.0:443
```

After some time to generate keys (it can take a while, sending dots in PM), you should see a message like this:

```
....
An updated CRL has been created.
CRL file: /etc/openvpn/pki/crl.pem


Executing command: `docker run -v ovpn_data:/etc/openvpn -d --restart=always --name ovpn_udp -p 34231:1194/udp --cap-add=NET_ADMIN kylemanna/openvpn ovpn_run --proto udp`fb79eba733f495b79878e3bc66710421a9ad0a931d217bc85c06e26fa098e659
Executing command: `docker run -v ovpn_data:/etc/openvpn -d --restart=always --name ovpn_tcp -p 34231:1194/tcp --cap-add=NET_ADMIN kylemanna/openvpn ovpn_run --proto tcp`dc32c846369adf2a08dfd629bcb021ce23e709b571fb794c64898edecf3708ac
All done, init completed!
```

This message means that VPN is now running!

The last step is to create user profiles. This is very simple, you just need to write `/generate profile_name` to create .ovpn profile with the name `profile_name`. You can create as many profiles as you want.


<details>
<summary>How to run bot without docker-compose</summary>

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

</details>

