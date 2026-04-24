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
            # Suche das config Objekt
            config_match = re.search(r'const config = (\{[\s\S]*?\n\s*\};)', content)
            if config_match:
                return config_match.group(1)
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
        // UI-Parameter (8): Was der Benutzer sieht und steuert
        window.baseIds = ["Dir", "Track", "Train", "TrainNo", "Destination", "Remark", "Hour", "Minute"];
        // Intern-Parameter (11): Was an das Frontend gesendet wird
        const internalBaseIds = ["Dir", "Track", "Train", "TrainNo_1", "TrainNo_2", "TrainNo_3", "TrainNo_4", "Destination", "Remark", "Hour", "Minute"];
        
        // Status der aktuellen Anzeige (11 Slots pro Reihe × 2 Reihen = 22 Slots)
        let currentData = Array(22).fill(" ");
        // Warteschlange für zufällige Züge
        let trainQueue = [];
        
        // Füge TrainNo Config für UI-Steuerung hinzu (von config kopieren, wenn nicht vorhanden)
        if (!config["TrainNo"]) {
            config["TrainNo"] = {
                flapCharset: {
                    chars: [" ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
                    options: {"default": ["weight-300"]}
                }
            };
        }
        
        // Mappe UI-Indizes (0-7 pro Reihe) zu Data-Indizes (0-10 pro Reihe)
        function uiToDataIndex(uiIdx) {
            const row = Math.floor(uiIdx / 8);
            const colUI = uiIdx % 8;
            const rowOffset = row * 11;
            
            // UI: Dir(0), Track(1), Train(2), TrainNo(3), Destination(4), Remark(5), Hour(6), Minute(7)
            // Data: Dir(0), Track(1), Train(2), TrainNo_1(3), TrainNo_2(4), TrainNo_3(5), TrainNo_4(6), Destination(7), Remark(8), Hour(9), Minute(10)
            const mapping = [0, 1, 2, 3, 7, 8, 9, 10]; // UI Col → Data Col (TrainNo bleibt bei 3)
            return rowOffset + mapping[colUI];
        }
        
        // Berechne den Bereich für TrainNo in currentData
        function getTrainNoRange(uiIdx) {
            const row = Math.floor(uiIdx / 8);
            const rowOffset = row * 11;
            return { start: rowOffset + 3, end: rowOffset + 6 }; // Indices 3-6 (bzw. 14-17 in Row 2)
        }
        
        // Reverse-Lookup: Finde den exakten Charset-Wert zu einem gespeicherten Wert
        function findCharsetValue(charset, storedValue) {
            // Falls exact match
            if (charset.indexOf(storedValue) !== -1) return storedValue;
            // Substring-reverse: wenn storedValue ein Substring eines Charset-Elements ist, nutze es
            // z.B. storedValue="9", charset enthält "9¾", return "9¾"
            for (let char of charset) {
                if (char.includes(storedValue) && storedValue.length < char.length) {
                    return char;
                }
            }
            return storedValue;
        }
        function applyTrainColors(trainValue, trainNoValue) {
            // Zerlege TrainNo in 4 Ziffern (rechts-aligned)
            // Verwende Leerzeichen statt Nullen für führende Positionen
            const trainNoStr = String(trainNoValue).padStart(4, ' ').slice(-4);
            const result = [];
            
            for (let i = 0; i < 4; i++) {
                const digit = trainNoStr[i];
                const displayValue = /^\d$/.test(digit) ? digit : " ";
                result.push(displayValue);
            }
            console.log("applyTrainColors() result (no colors):", result);
            return result;
        }

        function generateRandomTrain() {
            // baseIds hat 8 UI-Parameter, aber config hat 11 (mit TrainNo_1-4)
            // Daher erstelle die Mapping zwischen UI-Parametern und Config-Keys
            const configKeys = ["Dir", "Track", "Train", "TrainNo_1", "TrainNo_2", "TrainNo_3", "TrainNo_4", "Destination", "Remark", "Hour", "Minute"];
            const result = [];
            
            for (let i = 0; i < configKeys.length; i++) {
                const configKey = configKeys[i];
                
                // TrainNo-Ziffern werden separat behandelt - füge Placeholder für alle 4 ein
                if (configKey.startsWith("TrainNo_")) {
                    result.push(null); // Wird später mit Farben gefüllt
                    continue;
                }
                
                // Andere Parameter
                const chars = (config[configKey] && config[configKey].flapCharset) ? config[configKey].flapCharset.chars : ["fw", "bw"];
                // Wähle ein zufälliges Element ab Index 1 (um space zu vermeiden)
                const randIdx = 1 + Math.floor(Math.random() * (chars.length - 1));
                result.push(chars[randIdx] || chars[0] || " ");
            }
            
            // Generiere TrainNo zwischen 1 und 9999
            const trainNo = Math.floor(Math.random() * 9999) + 1;
            
            // Wende Farben an und zerlege in 4 Ziffern (Train ist bei Index 2)
            const trainValue = result[2];  // "のぞみ..." oder ähnlich
            console.log("generateRandomTrain(): trainValue=", trainValue);
            const trainNoZiffern = applyTrainColors(trainValue, trainNo);
            
            // Ersetze die Placeholders mit den farbigen Ziffern
            result[3] = trainNoZiffern[0]; // TrainNo_1
            result[4] = trainNoZiffern[1]; // TrainNo_2
            result[5] = trainNoZiffern[2]; // TrainNo_3
            result[6] = trainNoZiffern[3]; // TrainNo_4
            
            console.log("generateRandomTrain() result:", result);
            return result;
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
            const row1Empty = currentData.slice(0, 11).every(val => val === " ");
            const row2Empty = currentData.slice(11, 22).every(val => val === " ");

            console.log("nextTrain(): row1Empty=" + row1Empty + ", row2Empty=" + row2Empty + ", trainQueue.length=" + trainQueue.length);

            if (row1Empty) {
                // Reihe 1 füllen
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    console.log("nextTrain(): Filling Row 1 with:", newTrain);
                    console.log("  - TrainNo_1-4 (indices 3-6):", newTrain[3], newTrain[4], newTrain[5], newTrain[6]);
                    for(let i = 0; i < 11; i++) currentData[i] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            } else if (row2Empty) {
                // Reihe 2 füllen
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    console.log("nextTrain(): Filling Row 2 with:", newTrain);
                    console.log("  - TrainNo_1-4 (indices 3-6):", newTrain[3], newTrain[4], newTrain[5], newTrain[6]);
                    for(let i = 0; i < 11; i++) currentData[i + 11] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            } else {
                // Beide voll -> R2 rutscht auf R1, R2 bekommt neuen Zug
                console.log("nextTrain(): Both rows full, shifting and adding new train to Row 2");
                for(let i = 0; i < 11; i++) currentData[i] = currentData[i + 11];
                if (trainQueue.length > 0) {
                    const newTrain = trainQueue.shift();
                    console.log("nextTrain(): Filling Row 2 with:", newTrain);
                    console.log("  - TrainNo_1-4 (indices 3-6):", newTrain[3], newTrain[4], newTrain[5], newTrain[6]);
                    for(let i = 0; i < 11; i++) currentData[i + 11] = newTrain[i];
                    trainQueue.push(generateRandomTrain());
                }
            }
            console.log("nextTrain(): currentData after update (Row 1 TrainNo):", currentData[3], currentData[4], currentData[5], currentData[6]);
            sendUpdate();
        }

        // Manuelle Steuerung - nur für UI-Parameter
        window.move = function(uiIdx, dir) {
            try {
                console.log("move() called with uiIdx=" + uiIdx + ", dir=" + dir);
                const step = parseInt(document.getElementById('jump-step').value) || 1;
                const colUI = uiIdx % 8;  // 8 UI-Parameter pro Reihe
                const id = window.baseIds[colUI];
                console.log("move(): colUI=" + colUI + ", id=" + id);
                const charset = (config[id] && config[id].flapCharset) ? config[id].flapCharset.chars : [" ", "fw", "bw"];
                
                if (id === "TrainNo") {
                    // Spezialfall für TrainNo - behandle als 4-stellige Zahl
                    const range = getTrainNoRange(uiIdx);
                    
                    // Kombiniere die 4 Ziffern zu einer Zahl
                    let trainNo = 0;
                    for (let i = 0; i < 4; i++) {
                        const digit = currentData[range.start + i];
                        trainNo = trainNo * 10 + (digit === " " ? 0 : parseInt(digit));
                    }
                    
                    console.log("move() TrainNo: vorher=" + trainNo + ", step=" + step + ", dir=" + dir);
                    
                    // Addiere/Subtrahiere step zur gesamten Zahl
                    if (dir === 'next') {
                        trainNo = (trainNo + step) % 10000; // Max 9999
                    } else {
                        trainNo = (trainNo - step + 10000) % 10000; // Min 0
                    }
                    
                    console.log("move() TrainNo: nachher=" + trainNo);
                    
                    // Zerlege die Nummer wieder in 4 Ziffern
                    const trainIdx = Math.floor(uiIdx / 8) === 0 ? 2 : 13; // Train Index in currentData
                    const trainValue = currentData[trainIdx];
                    const trainNoZiffern = applyTrainColors(trainValue, trainNo);
                    
                    // Schreibe die 4 Ziffern zurück
                    for (let i = 0; i < 4; i++) {
                        currentData[range.start + i] = trainNoZiffern[i];
                    }
                } else {
                    // Normale Parameter
                    const dataIdx = uiToDataIndex(uiIdx);
                    const currentValue = currentData[dataIdx];
                    console.log("move(): Normale Parameter, dataIdx=" + dataIdx + ", currentValue=" + currentValue);
                    
                    // Extrahiere nur die Basis-Ziffer (falls "9¾", nimm "9")
                    const baseValue = currentValue ? currentValue.split(" ")[0] : " ";
                    
                    let currentIdx = charset.indexOf(baseValue);
                    if (currentIdx === -1) {
                        // Versuche Substring-Match NICHT für TrainNo!
                        for (let i = 0; i < charset.length; i++) {
                            if (charset[i].includes(baseValue)) {
                                currentIdx = i;
                                break;
                            }
                        }
                        if (currentIdx === -1) currentIdx = 0;
                    }
                    
                    if (dir === 'next') {
                        currentIdx = (currentIdx + step) % charset.length;
                    } else {
                        currentIdx = (currentIdx - step + charset.length) % charset.length;
                    }
                    
                    // Speichere die Basis-Zahl nur (z.B. "9" für "9¾")
                    const newExactValue = charset[currentIdx];
                    const newBaseValue = newExactValue.split(" ")[0];
                    console.log("move(): Setting index " + dataIdx + " to " + newBaseValue);
                    currentData[dataIdx] = newBaseValue;
                }
                
                sendUpdate();
            } catch(e) {
                console.error("Fehler in move():", e);
            }
        }

        function createButtons(containerId, offset) {
            const container = document.getElementById(containerId);
            window.baseIds.forEach((id, i) => {
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