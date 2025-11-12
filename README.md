# OpenTunnel V4 Proxy Server

### OpenTunnel (V4) is an open source proxy server made to bypass network download restrictions.

## Client Usage

- Download release from:

`https://github.com/KoPlayz/OpenTunnel/releases/tag/client-v4.1.0`

for your device

- Update port, server ip/domain, and set API Key

- run client



## Server Usage

- Download latest server release from:

`https://github.com/KoPlayz/OpenTunnel/releases/tag/server-v4.0.0`

- Run server to generate initial config (replace {python} with the path to your python binary)

`$ {python} fileserver.py`

- Edit config.json to specify port, logs directory, and API key (use any text editor)

`$ nano config.json`

- Rerun server with new config

`$ {python} fileserver.py`