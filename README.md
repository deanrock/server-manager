server-manager
==============



How to generate SSL key?
========================

	$ cd . #to server-manager root (e.g. /home/manager/server-manager/)
	$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ssl.key -out ssl.crt

