# Sywer
HTTP server

## Installation
- Linux:
	- To install Sywer: `sudo sh setup.sh`
	- If you want to run the server at the boot: `sudo systemctl start sywer`
- Termux:

	- To install Sywer:
	```
	cd ressources
	sh termux_setup.sh
	```

## Settings
All the settings must be in the file `settings.swy`.
The differents settings are:
```
port	[PORT]
path	[SRC_PATH]
index	[FILE]
logs	[FILE]
icon	[FILE]
error_[error]	[FILE]
```

## Path
- `settings.swy`:
	- Linux: `/var/lib/sywer/settings.swy`
	- Termux: `/data/data/com.termux/files/usr/var/lib/sywer/settings.swy`
	- Windows: `./settings`
- `logs.swy`: 
	- Linux: `/var/log/sywer/logs.swy`
	- Termux: `/data/data/com.termux/files/usr/var/lib/sywer/settings.swy`
	- Windows: *None*
