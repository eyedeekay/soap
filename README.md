# soap

Go implementation of an Unciv multiplayer server which operates on hidden services.
It uses the SAMv3 API provided by I2P to set up an efficient, secure Unciv server
for you and your friends.

STATUS: This project is maintained. It is feature-complete as far as I know.

Usage
-----

Simply compile and run. More options might be added in the future but for now it's
no-frills.

```sh
go build
./soap
```

You can obtain the URL from the terminal output, or by visiting [http://127.0.0.1:7669](http://127.0.0.1:7669)
to reach the information page which helps you share the server.

Client Configuration
--------------------

This server is intended for use with regular UnCiv. In order to launch UnCiv with
I2P as an HTTP Proxy, use the following comand:

```sh
java -Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=4444 -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=4444 -jar Unciv.jar
```
