# transmission-monitor

Tool to track Transmission torrents state. Intended to be run as a cron job each few minutes. Interacts with Transmission RPC and notifies (via external command) about unhealthy torrents once.

For each faulty torrent external command gets fed with JSON describing the problem and the torrent. If command succeeds (zero exit code), notification delivery is considered successful and will not repeated again. Otherwise notification delivery will be retried on next run. You may clear entire database or remove single key.

Remote RPC must be enabled in Transmission for this program to work.

## Installation

#### Binaries

Pre-built binaries are available [here](https://github.com/Snawoot/transmission-monitor/releases/latest).

#### Build from source

Alternatively, you may install transmission-monitor from source. Run the following within the source directory:

```
make install
```

## Configuration

Configuration example:

#### /home/user/.config/transmission-monitor.yaml

```yaml
rpc:
  user: transmissionuser
  password: transmissionpassword
notify:
  command:
    - /home/user/.config/transmission-notify.sh
```

#### /home/user/.config/transmission-notify.sh

```bash
#!/bin/bash

set -euo pipefail

jq -r '"There is a problem with following torrent:\n\nName: \"" + .torrent.name + "\"\nHash: " + .torrent.hashString + "\nComment: " + .torrent.comment + "\nCause: " + .reason' | \
mailx -v \
-r "sender@example.com" \
-s "Torrent requires attention" \
-S smtp="mx.example.com:587" \
-S smtp-use-starttls \
-S smtp-auth=login \
-S smtp-auth-user="sender@example.com" \
-S smtp-auth-password="mailpassword" \
recipient@example.com
```

Make sure to run `transmission-monitor` command every few minutes with scheduler of your choice.

## Synopsis

```
$ ./bin/transmission-monitor -h
Usage of transmission-monitor:
  -clear-db
    	clear database
  -clear-key string
    	delete specified hash from database
  -conf string
    	path to configuration file (default "/home/user/.config/transmission-monitor.yaml")
  -version
    	show program version and exit
```
