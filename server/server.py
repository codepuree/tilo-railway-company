""" This module hosts the server needed for "Tilos Train Terminal" """

from bottle import request, Bottle, abort, static_file
app = Bottle()
web_sockets = []

@app.route('/')
@app.route('/index.html')
def index():
    return static_file('index.html', root='../')

@app.route('/<filename:path>')
def return_static_files(filename):
    return static_file(filename, root='../')

@app.route('/websocket')
def handle_websocket():
    web_socket = request.environ.get('wsgi.websocket')
    web_sockets.append(web_socket)

    if not web_socket:
        abort(400, 'Expected WebSocket request.')

    while True:
        try:
            message = web_socket.receive()

            for other_web_socket in web_sockets:
                if other_web_socket != web_socket:
                    other_web_socket.send(message)
        except WebSocketError:
            print('WebSocketError:', WebSocketError)
            break

from gevent.pywsgi import WSGIServer
from geventwebsocket import WebSocketError
from geventwebsocket.handler import WebSocketHandler
server = WSGIServer(("0.0.0.0", 8080), app,
                    handler_class=WebSocketHandler)
server.serve_forever()

def main():
    print("Hello, world!")

if __name__ == '__main__':
    main()