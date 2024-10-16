![alt text](resources/inhalt/schrift_merkmale/Logo.png)

# tilo-railway-company

## Setup

To install the required packages, navigate into the server folder and run the following command.

`> pip install -r .\requirements.txt`

## Start the server

To start the server, navigate into the server folder and run the following command.

`> python .\server.py`

You can now open a browser an navigate to this url http://localhost:8080/

## Commands

### Compile and send to raspberry pi

```bash
GOOS=linux GOARCH=arm GOARM=7 go build -a -tags netgo -ldflags '-w' -o ./bin/trc ./cmd/ && scp bin/trc root@train:/root/Schreibtisch/trc/trc
```

### Extract trclib

```bash
yaegi extract -name github.com/codepuree/tilo-railway-company/pkg/trclib github.com/codepuree/tilo-railway-company/pkg/traincontrol
```

### install yaegi
```bash
go get -u github.com/traefik/yaegi/cmd/yaegi
```

### Clear chromium cache

```bash
rm -r /home/pi/.cache/chromium/Default/Cache/*
```

### Layout Information

Direction: 
clockwise: F
counterclockwise: B

Blocks: A, B, C. D, F, G
A: Track 1, closest to viewer/edge, station
B: Track 2
C: Track 4
D: Track 4, most far away from viewer/edge, station

F: Turnouts (east, west combined)
G: open rail (countryside)
