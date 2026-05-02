# SmartHomeHub
A self-hosted smart home control system written in Go, integrating real devices (Philips Hue) and simulated devices behind a unified control plane.


Current routes:


## GET

`/health`               - Health check
`/devices`              - Lists all devices
`/devices/{id}/state`   - Gives state of specific device id


## POST
`/devices/{id}/command` - Sends a command to device id. Command structure depends on device type.