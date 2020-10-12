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

### Clear chromium cache

```bash
rm -r /home/pi/.cache/chromium/Default/Cache/*
```