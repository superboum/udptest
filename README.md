udptest
=======

Discover your firewall UDP output policy with this app.

## Known limitations

 * This is absolutely not secure
 * Only one client at a time can use the server
 * If a client does not close a port, the port will not be closed
 * Use it on a server without any service listening on UDP. Or cross fingers...

## Usage

On server : `sudo ./server.go`

On client : `./client`  
Client will ask for server IP and Port. Default port is 80.

A full log will be available in `port-checker.txt`

You have to authorize all UDP traffic on your test server (be careful to add a DMZ if you're behind a NAT).
This app has been successfully tested on DigitalOcean servers.

## Compile

Sorry, does not follow golang best practises...

```
go build client.go porttest.go checkerror.go
go build server.go porttest.go checkerror.go
```

## How does it work ?

The server binary will listen on port 80.

There is 2 available HTTP Requests :

```
POST /port {"port": "1005"} - Will open port 1005
DELETE /port {"port": "1005"} - Will close port 1005
```

The client will iterate on the following algorithm for every ports :

 1. Send a request to open a port
 2. Try to connect on UDP on this port
 3. Wait 1 sec for the answer `OK`
 4. will log the port as open if it receives this answer or as closed if it receives nothing or something else.
 5. After that, the client will send a request to close the port.
