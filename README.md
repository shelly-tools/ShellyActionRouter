# ShellyActionRouter
ShellyActionRouter is a tiny action proxy written in Golang to allow multiple Actions per Shelly..  

### Installation
no installation needed, just extract the release files to a directory and execute the ShellyActionRouter binary.

### Configuration
add your actions to the actions.ini file.. see actions.ini for examples

## configure the Shelly to trigger the action proxy
```
http://<shelly-proxy:8888>/api/action/action1

```
where ```action1``` is the section name in the actions.ini file
