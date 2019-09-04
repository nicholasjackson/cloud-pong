# cloud-pong
Terminal based Cloud Pong

## Running the game

The game requires two servers which are connected together and two clients connect to each of the servers.  To start a basic setup you can use the Makefile provided in the repo.

## Controls

### Player 1
* W - Bat up
* S - Bat down
* E - Serve
* Ctrl+R - Reset game

### Player 2
* O - Bat up
* L - Bat down
* P - Serve

### Run player 1 API

```
$make start-server-1
(cd api && PLAYER=1 BIND_PORT=6000 UPSTREAM_ADDRESS=localhost:6001 go run main.go)
2019-08-23T19:48:55.993+0100 [INFO]  Dialing connection: server=true
2019-08-23T19:48:55.994+0100 [INFO]  Listening on port: port=6000 player=1
```

### Run player 2 API

```
$ make start-server-2
(cd api && PLAYER=2 BIND_PORT=6001 UPSTREAM_ADDRESS=localhost:6000 go run main.go)
2019-08-23T19:49:33.088+0100 [INFO]  Dialing connection: server=true
2019-08-23T19:49:33.089+0100 [INFO]  Listening on port: port=6001 player=2
2019-08-23T19:49:33.091+0100 [INFO]  Client succesfully connected to the server
```

### Start the player one application which connects to player 1 API

```
$ make player-1
```

![](images/player1.png)


### Start the player two application which connects to player 2 API

```
$ make player-2
```

![](images/player2.png)


## WARNING
This code was a fun weekend hack, the code is largely held together by luck and bandages.
