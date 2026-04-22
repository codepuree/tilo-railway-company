import os
import re
import random
from flask import Flask, render_template_string, send_from_directory, jsonify

base_dir = os.path.abspath(os.path.dirname(__file__))
app = Flask(__name__)

def get_config_raw():
    path = os.path.join(base_dir, 'splitflap_trc.js')
    try:
        with open(path, 'r', encoding='utf-8') as f:
            content = f.read()
            config_match = re.search(r'const config = (\{.*?\});', content, re.DOTALL)
            if config_match: return config_match.group(1)
    except: pass
    return "{}"

@app.route('/')
def index():
    return render_template_string('''
<!DOCTYPE html>
<html>
<head>
    <title>Splitflap Train Manager</title>
    <style>
        body, html { margin: 0; padding: 0; height: 100%; font-family: 'Segoe UI', sans-serif; background: #1a1a1a; color: white; overflow: hidden; }
        #display-frame { width: 100%; height: 60vh; border: none; background: #000; }
        #control-panel { width: 100%; height: 40vh; background: #2d2d2d; border-top: 3px solid #444; padding: 15px; box-sizing: border-box; overflow-y: auto; }
        
        .main-nav { display: flex; gap: 20px; align-items: center; margin-bottom: 20px; padding: 10px; background: #383838; border-radius: 8px; }
        .btn-main { padding: 12px 24px; background: #72ea7a; color: black; border: none; border-radius: 5px; font-weight: bold; cursor: pointer; font-size: 16px; }
        .btn-main:hover { background: #5cd664; }
        
        .row-container { display: flex; gap: 15px; }
        .row-section { flex: 1; background: #222; padding: 10px; border-radius: 8px; border: 1px solid #444; }
        .row-title { color: #72ea7a; font-size: 12px; font-weight: bold; margin-bottom: 10px; text-transform: uppercase; }
        .controls { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 5px; }
        .card { background: #3d3d3d; padding: 5px; border-radius: 4px; text-align: center; font-size: 10px; }
        .btn-group { display: flex; justify-content: center; gap: 3px; margin-top: 3px; }
        button { cursor: pointer; background: #555; color: white; border: none; border-radius: 2px; padding: 2px 6px; }
        button:hover { background: #72ea7a; color: black; }
        
        .queue-info { font-size: 12px; color: #888; }
    </style>
</head>
<body>
    <iframe id="display-frame" src="splitflap_trc.html"></iframe>
    <div id="control-panel">
        
        <div class="main-nav">
            <button class="btn-main" onclick="nextTrain()">NEXT TRAIN ❯❯</button>
            <div class="queue-info">Züge in Warteschlange: <span id="q-count">0</span></div>
            <div style="margin-left: auto;">Jump: <input type="number" id="jump-step" value="1" style="width:40px; background:#111; color:#72ea7a; border:1px solid #444;"></div>
        </div>

        <div class="row-container">
            <div class="row-section">
                <div class="row-title">Obere Reihe (Reihe 1)</div>
                <div id="controls-r1" class="controls"></div>
            </div>
            <div class="row-section">
                <div class="row-title">Untere Reihe (Reihe 2)</div>
                <div id="controls-r2" class="controls"></div>
            </div>
        </div>
    </div>

    <script>
        const config = {{ charsets | safe }};
        const baseIds = ["Dir", "Track", "Train", "TrainNo", "Destination", "Remark", "Hour", "Minute"];
        
        // Status der aktuellen Anzeige (16 Slots)
        let currentData = Array(16).fill(" ");
        // Warteschlange für zufällige Züge
        let trainQueue = [];

        function generateRandomTrain() {
            const train = [];
            baseIds.forEach(id => {
                const chars = (config[id] && config[id].flapCharset) ? config[id].flapCharset.chars : ["fw", "bw"];
                // Zufälligen Index wählen (ohne das erste Leerzeichen, wenn möglich)
                const randIdx = Math.floor(Math.random() * (chars.length - 1)) + 1;
                train.push(chars[randIdx] || " ");
            });
            return train;
        }

        // Initial 10 Züge generieren
        for(let i=0; i<10; i++) trainQueue.push(generateRandomTrain());

        function sendUpdate() {
            const iframe = document.getElementById('display-frame');
            if (iframe.contentWindow) {
                iframe.contentWindow.postMessage(currentData, "*");
            }
            document.getElementById('q-count').innerText = trainQueue.length;
        }

        function nextTrain() {
            const row1Empty = currentData.slice(0,8).every(val => val === " ");
            const row2Empty = currentData.slice(8,16).every(val => val === " ");

            if (row1Empty) {
                // Reihe 1 füllen
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    for(let i=0; i<8; i++) currentData[i] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            } else if (row2Empty) {
                // Reihe 2 füllen
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    for(let i=0; i<8; i++) currentData[i+8] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            } else {
                // Beide voll -> R2 rutscht auf R1, R2 bekommt neuen Zug
                // 1. R2 Daten nach R1 kopieren
                for(let i=0; i<8; i++) currentData[i] = currentData[i+8];
                // 2. Neuen Zug für R2 holen
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    for(let i=0; i<8; i++) currentData[i+8] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            }
            sendUpdate();
        }

        // Manuelle Steuerung
        window.move = function(idx, dir) {
            const step = parseInt(document.getElementById('jump-step').value) || 1;
            const id = baseIds[idx % 8];
            const charset = (config[id] && config[id].flapCharset) ? config[id].flapCharset.chars : [" ", "fw", "bw"];
            
            let currentIdx = charset.indexOf(currentData[idx]);
            if (currentIdx === -1) currentIdx = 0;

            if (dir === 'next') {
                currentIdx = (currentIdx + step) % charset.length;
            } else {
                currentIdx = (currentIdx - step + charset.length) % charset.length;
            }
            currentData[idx] = charset[currentIdx];
            sendUpdate();
        }

        function createButtons(containerId, offset) {
            const container = document.getElementById(containerId);
            baseIds.forEach((id, i) => {
                const card = document.createElement('div');
                card.className = 'card';
                card.innerHTML = `<h3>${id}</h3><div class="btn-group">
                    <button onclick="move(${i + offset}, 'prev')">◀</button>
                    <button onclick="move(${i + offset}, 'next')">▶</button>
                </div>`;
                container.appendChild(card);
            });
        }

        createButtons('controls-r1', 0);
        createButtons('controls-r2', 8);
        
        // Kurzer Delay beim Start damit Iframe bereit ist
        setTimeout(sendUpdate, 1000);
    </script>
</body>
</html>
    ''', charsets=get_config_raw())

@app.route('/<path:filename>')
def serve_file(filename):
    return send_from_directory(base_dir, filename)

if __name__ == '__main__':
    app.run(port=8080, debug=True)