# ShellyActionRouter
ShellyActionRouter is a tiny action proxy written in Golang to allow multiple Actions per Shelly..  

### Installation
no installation needed, just extract the release files to a directory and execute the related ShellyActionRouter binary for your system. 
For linux it is recommended to place all files from the archive to /opt/shellyactionrouter

#### Linux
make the binary executable via chmod
```
sudo chmod +x /opt/shellyactionrouter/ShellyActionRouter-linux-amd64
```

#### Linux on ARM (Rapsberry Pi)
make the binary executable via chmod
```
sudo chmod +x /opt/shellyactionrouter/ShellyActionRouter-linux-arm-5
```

#### make ShellyActionRouter a systemd service
create a file for ShellyActionRouter

```
sudo vi /lib/systemd/system/shellyactionrouter.service 
```

Content:
```
[Unit]
Description=ShellyActionRouter service.

[Service]
Type=simple
WorkingDirectory=/opt/shellyactionrouter
ExecStart=/opt/shellyactionrouter/ShellyActionRouter-linux-amd64

[Install]
WantedBy=multi-user.target
```

*Attention*: ExecStart needs to be the path to the binary, e.g. ExecStart=/opt/shellyactionrouter/ShellyActionRouter-linux-amd64 for 64 bit Linux or 
ShellyActionRouter-linux-arm5 for Linux on ARM (Raspberry Pi)

copy the file to /etc/systemd/system

```
sudo cp /lib/systemd/system/shellyactionrouter.service /etc/systemd/system/
```
Now you can start it via 
```
sudo systemctl start shellyactionrouter
```

automatic startup on boot can be enabled by:
```
sudo systemctl enable shellyactionrouter
```

### Configuration
Open a browser and go to the router url, the WebUI should be self explaining. all URLs and Actions will be stored in the ActionRouter.db (sqlite)
```
http://<shelly-proxy:8888>
```

#### Port configuration
To confire a custom port you can edit the config.ini. The default port is 8888
