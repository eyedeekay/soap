# soap

Go implementation of an Unciv multiplayer server which operates on hidden services.
It uses the SAMv3 API provided by I2P to set up an efficient, secure Unciv server
for you and your friends.

Usage
-----

Simply compile and run. More options might be added in the future but for now it's
no-frills.

```sh
go build
./soap
```


Client Configuration
--------------------

This server is intended for use with regular UnCiv. In order to launch UnCiv with
I2P as an HTTP Proxy, use the following comand:

```sh
java -Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=4444 -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=4444 -jar Unciv.jar
```
