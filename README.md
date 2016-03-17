udptest
=======

Discover your firewall UDP output policy with this app.

## Usage

On server : `sudo ./server.go`

On client : `./client`  
Client will ask for server IP and Port. Default port is 80.

You have to authorize all UDP traffic on your test server (be careful to add a DMZ if you're behind a NAT).

This app has been successfully tested on DigitalOcean servers.

## Compile

Sorry, does not follow golang best practises...

```
go build client.go porttest.go checkerror.go
go build server.go porttest.go checkerror.go
```
