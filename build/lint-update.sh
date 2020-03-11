#!/bin/sh

echo "-- Downloading latest revive-config"
curl -L -H 'Cache-Control: no-cache' https://bitbucket.org/!api/2.0/snippets/tktip/yeba9j/files/revive.toml > revive.toml
echo "-- Revive-config updated\n"
