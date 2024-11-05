#!/bin/bash
# delete all "retained" messages
echo "cleaning " $1 " :: usage: cleanmqtt <host>"
host="$1"
# get your MQTT credentials from HA UI: Settings > Devices & Services > MQTT > Configure > Reconfigure MQTT
user=homeassistant
paswd="yourPassword"

mosquitto_sub -h "$host" -u "$user" -P "$paswd" -t "#" -v --retained-only | while read topic value
do
  echo "cleaning topic $topic"
  mosquitto_pub -h "$host" -u "$user" -P "$paswd" -t "$topic" -r -n
done
