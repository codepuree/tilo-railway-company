// GLOBALE KONFIGURATION - wird von Flask extrahiert!
const config = {
    "Dir": { 
        relW: 2, relH: 25, 
        pos: {
            "fw": { relTop: 33, relLeft: 10.3 },
            "bw": { relTop: 33, relLeft: 96.1 } 
            }
    }, 
    "Track": {
        relW: 5, relH: 26, relTop: 33, relLeft: 14.75,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { chars: [
                " ", "1", "2", "3", "4", "5", "6", "7", "8", "9¾"
            ], 
            options: {"default": ["weight-300"]}
        }
    },
    "Train": {
        relW: 12.35, relH: 26, relTop: 33, relLeft: 21.1,
        offsetTop: 4, charOffset: 8, fontSizeFactor: 50, lineHeight: 1.1,
        flapCharset: {
    chars: [
        " ", 
        "のぞみ\nNOZOMI", "ひかり\nHIKARI", "こだま\nKODAMA", 
        "はやぶさ\nHAYABUSA", "こまち\nKOMACHI", "かがやき\nKAGAYAKI", 
        "さくら\nSAKURA", "つばめ\nTSUBAME", "みずほ\nMIZUHO", 
        "あさま\nASAMA", "とき\nTOKI", "たにがわ\nTANIGAWA", 
        "なすの\nNASUNO", "やまびこ\nYAMABIKO", "はやて\nHAYATE", 
        "つばさ\nTSUBASA", "かもめ\nKAMOME", 
        "成田エクスプレス\nNARITA EXP", "はるか\nHARUKA", 
        "あずさ\nAZUSA", "ひたち\nHITACHI", "ソニック\nSONIC", 
        "ホグワーツ特急\nHOGWARTS EXP"
    ],
    options: {
        "default": ["weight-400"],
        "のぞみ\nNOZOMI": ["flap-blue", "txt-white"], 
        "ひかり\nHIKARI": ["flap-blue", "txt-white"], 
        "こだま\nKODAMA": ["flap-blue", "txt-white"],
        "はやぶさ\nHAYABUSA": ["flap-turquoise", "txt-pink"], 
        "はやて\nHAYATE": ["flap-turquoise", "txt-white"],
        "こまち\nKOMACHI": ["flap-red", "txt-white"], 
        "かがやき\nKAGAYAKI": ["flap-yellow", "txt-black"], 
        "あさま\nASAMA": ["flap-yellow", "txt-black"],
        "さくら\nSAKURA": ["flap-orange", "txt-black"], 
        "みずほ\nMIZUHO": ["flap-orange", "txt-black"],
        "つばめ\nTSUBAME": ["flap-white", "txt-black"],
        "かもめ\nKAMOME": ["flap-white", "txt-black"],
        "とき\nTOKI": ["flap-red", "txt-white"],
        "つばさ\nTSUBASA": ["flap-orange", "txt-black"],
        "たにがわ\nTANIGAWA": ["flap-white", "txt-black"],
        "なすの\nNASUNO": ["flap-white", "txt-black"],
        "やまびこ\nYAMABIKO": ["flap-white", "txt-black"],
        "成田エクスプレス\nNARITA EXP": ["bg-darkred", "txt-white"], 
        "はるか\nHARUKA": ["flap-blue", "txt-white"],
        "ソニック\nSONIC": ["flap-blue", "txt-white"],
        "あずさ\nAZUSA": ["flap-white", "txt-black"],
        "ひたち\nHITACHI": ["flap-white", "txt-black"],
        "ホグワーツ特急\nHOGWARTS EXP": ["bg-darkred", "txt-white"]
    }
}
    },
    "TrainNo_1": {
        relW: 2.25, relH: 26, relTop: 33, relLeft: 34.8,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [" ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
            options: {"default": ["weight-300"]}
        }
    },
    "TrainNo_2": {
        relW: 2.25, relH: 26, relTop: 33, relLeft: 37.2,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [" ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
            options: {"default": ["weight-300"]}
        }
    },
    "TrainNo_3": {
        relW: 2.25, relH: 26, relTop: 33, relLeft: 39.6,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [" ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
            options: {"default": ["weight-300"]}
        }
    },
    "TrainNo_4": {
        relW: 2.25, relH: 26, relTop: 33, relLeft: 42,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [" ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"],
            options: {"default": ["weight-300"]}
        }
    },
    "Destination": {
        relW: 15.5, relH: 26, relTop: 33, relLeft: 45.5,
        offsetTop: 4, charOffset: 8, fontSizeFactor: 50, lineHeight: 1.1,
        flapCharset: {
            chars: [
                " ", "東京\nTOKYO", "新大阪\nSHIN-OSAKA", "京都\nKYOTO", "名古屋\nNAGOYA", 
                "横浜\nYOKOHAMA", "博多\nHAKATA", "札幌\nSAPPORO", "仙台\nSENDAI", 
                "広島\nHIROSHIMA", "神戸\nKOBE", "岡山\nOKAYAMA", "金沢\nKANAZAWA", 
                "新潟\nNIIGATA", "静岡\nSHIZUOKA", "浜松\nHAMAMATSU", "奈良\nNARA", 
                "長野\nNAGANO", "盛岡\nMORIOKA", "青森\nAOMORI", "熊本\nKUMAMOTO", 
                "鹿児島\nKAGOSHIMA", "長崎\nNAGASAKI", "大分\nOITA", "高山\nTAKAYAMA", 
                "箱根\nHAKONE", "日光\nNIKKO", "鎌倉\nKAMAKURA", "松本\nMATSUMOTO", 
                "姫路\nHIMEJI", "宮島\nMIYAJIMA", "伊勢\nISE", "別府\nBEPPU", 
                "熱海\nATAMI", "コモ\nCOMO", "ミラノ\nMILANO", "ローマ\nROMA", 
                "ベルリン\nBERLIN", "ミュンヘン\nMÜNCHEN","ホグワーツ\nHOGWARTS", "ロンドン\nLONDON", 
                "パリ\nPARIS", "ニューヨーク\nNEW YORK", "ソウル\nSEOUL", 
                "バンコク\nBANGKOK", "チューリッヒ\nZÜRICH", "ウィーン\nWIEN", 
                "マドリード\nMADRID", "アムステルダム\nAMSTERDAM", "オスロ\nOSLO", 
                "ストックホルム\nSTOCKHOLM"
            ],
            options: {"default": ["weight-400"]} 
        }
    },
    "Remark": {
        relW: 18.5, relH: 26, relTop: 33, relLeft: 62.3,
        offsetTop: 4, charOffset: 8, fontSizeFactor: 50, lineHeight: 1.1,
        flapCharset: {
            chars: [
                " ", "運休\nCANCELLED", "遅れ\nDELAYED", "編成逆転\nREVERSE FORMATION", "全車指定\nALL SEATS RESERVED", "自由席\nNON-RESERVED CAR", "当駅止まり\nTERMINATES HERE", "回送\nOUT OF SERVICE", "臨時\nSPECIAL TRAIN", "満席\nNO VACANCY", "空席あり\nSEATS AVAILABLE","魔法界より\nWIZARDING WORLD", "姿現し禁止\nNO APPARATING", "ふくろう便のみ\nBY OWL POST"
            ],
            options: {
                "default": ["weight-400"],
                "運休\nCANCELLED": ["flap-white", "txt-red"], "遅れ\nDELAYED": ["flap-yellow", "txt-black"],
                "編成逆転\nREVERSE FORMATION": ["flap-orange", "txt-black"], "全車指定\nALL SEATS RESERVED": ["flap-blue", "txt-white"],
                "自由席\nNON-RESERVED CAR": ["flap-green", "txt-black"], "当駅止まり\nTERMINATES HERE": ["txt-white"],
                "回送\nOUT OF SERVICE": ["flap-white", "txt-black"], "臨時\nSPECIAL TRAIN": ["flap-turquoise", "txt-white"],
                "満席\nNO VACANCY": ["bg-darkred", "txt-white"], "空席あり\nSEATS AVAILABLE": ["flap-green", "txt-black"]
            }
        }
    },
    "Hour": {
        relW: 5, relH: 26, relTop: 33, relLeft: 82,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [
                "  ", "00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", 
                "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", 
                "21", "22", "23" 
            ], 
            options: {"default": ["weight-300"]}
        }
    },
    "Minute": {
        relW: 5, relH: 26, relTop: 33, relLeft: 88.3,
        offsetTop: 6, charOffset: 15, fontSizeFactor: 110, 
        flapCharset: { 
            chars: [
                "  ", "00", "01", "02", "03", "04", "05", "06", "07", "08", "09",
                "10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
                "20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
                "30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
                "40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
                "50", "51", "52", "53", "54", "55", "56", "57", "58", "59"
            ], 
            options: {"default": ["weight-300"]}
        }
    }
};

class DirectionIndicator {
    constructor(id, config) {
        this.id = id;
        this.cfg = config;
        this.$el = $("#" + id);
        this.currentDir = ""; 
    }

    resize(bgDim) {
        this.lastBgDim = bgDim; 
        this.updatePosition(this.currentDir);
    }

    updatePosition(val) {
        if (!this.lastBgDim) return;

        const pos = (this.cfg.pos && this.cfg.pos[val]) ? this.cfg.pos[val] : this.cfg;
        
        const pW = this.lastBgDim.w * (this.cfg.relW / 100);
        const pH = this.lastBgDim.h * (this.cfg.relH / 100);
        const pT = this.lastBgDim.h * (pos.relTop / 100);
        const pL = this.lastBgDim.w * (pos.relLeft / 100);

        this.$el.css({
            "width": pW + "px",
            "height": pH + "px",
            "top": pT + "px",
            "left": pL + "px"
        });
    }

    display(val) {
        if (this.currentDir === val) return;
        this.currentDir = val;

        this.updatePosition(val);

        if (val === "fw") {
            this.$el.html('<img src="dir_f.svg">');
        } else if (val === "bw") {
            this.$el.html('<img src="dir_b.svg">');
        } else {
            this.$el.empty(); 
        }
    }
}

class SplitFlap {
    constructor(id, config) {
        this.id = id;
        this.cfg = config;
        this.$el = $("#" + id);
        
        const savedIndex = localStorage.getItem('flap_index_' + this.id);
        this.currentIndex = savedIndex ? parseInt(savedIndex) : 0;
        
        if (this.currentIndex >= this.cfg.flapCharset.chars.length) {
            this.currentIndex = 0;
        }

        this.currentValue = this.cfg.flapCharset.chars[this.currentIndex];
        
        this.buildStructure();
    }

    buildStructure() {
        this.$el.html(`
            <div class="trc-container">
                <div class="trc-spacer"></div>
                <div class="trc-char-wrapper">
                    <span class="trc-char"> </span>
                </div>
                <div class="trc-divider-base">
                    <div class="trc-divider-line"></div>
                    <div class="trc-link trc-left"></div>
                    <div class="trc-link trc-right"></div>
                </div>
            </div>
        `);
    }

    resize(bgDim) {
        const pW = bgDim.w * (this.cfg.relW / 100);
        const pH = bgDim.h * (this.cfg.relH / 100);
        const pT = bgDim.h * (this.cfg.relTop / 100);
        const pL = bgDim.w * (this.cfg.relLeft / 100);
        
        const scale = pH / 120;
        
        const fs = Math.floor((this.cfg.fontSizeFactor || 140) * scale);
        const lh = this.cfg.lineHeight ? this.cfg.lineHeight : 0.8;

        this.$el.css({
            "width": pW + "px",
            "height": pH + "px",
            "top": pT + "px",
            "left": pL + "px"
        });

        this.$el.find(".trc-char").css({
            "font-size": fs + "px",
            "line-height": lh,
            "margin-top": (this.cfg.charOffset || 0) * scale + "px"
        });

        const divH = Math.max(0.2, 2 * scale);
        const linkSize = Math.max(2, 12 * scale);

        this.$el.find(".trc-spacer").css("height", (this.cfg.offsetTop || 0) * scale + "px");
        this.$el.find(".trc-divider-line").css("height", divH + "px");
        this.$el.find(".trc-link").css({ "width": linkSize + "px", "height": linkSize + "px" });
        
        this.updateView();
    }

    updateView() {
        const formattedText = this.currentValue.replace(/\n/g, "<br>");
        this.$el.find(".trc-char").html(formattedText);
        this.applyStyles(this.currentValue);
    }

    display(target) {
    const charset = this.cfg.flapCharset.chars;
    let targetIndex = charset.indexOf(target);
    
    if (targetIndex === -1) targetIndex = 0;
    if (this.currentIndex === targetIndex) return;

    let path = [];
    let i = this.currentIndex;

    while (i !== targetIndex) {
        i++;
        if (i >= charset.length) {
            i = 0; 
        }
        path.push(charset[i]);
    }

    let step = 0;
    
    // --- Brake Configuration ---
    const decelerateSteps = 3;  
    const decelerateFactor = 3; 

    const anim = () => {
        if (step < path.length) {
            this.currentValue = path[step];
            this.currentIndex = charset.indexOf(this.currentValue);
            localStorage.setItem('flap_index_' + this.id, this.currentIndex);
            this.updateView();
            step++;

            let remaining = path.length - step;
            let currentDelay = 80; // Default Speed

            // Slowdown for the last few steps
            if (remaining < decelerateSteps) {
                currentDelay = 100 * decelerateFactor;
            }

            setTimeout(anim, currentDelay);
        }
    };
    anim();
}

    applyStyles(val) {
        const opt = this.cfg.flapCharset.options;
        const $span = this.$el.find(".trc-char");
        const $cont = this.$el.find(".trc-container");

        if (!opt) return;

        $cont.removeClass(function (index, className) {
            return (className.match(/(^|\s)(flap-|bg-)\S+/g) || []).join(' ');
        });
        $span.removeClass(function (index, className) {
            return (className.match(/(^|\s)txt-\S+/g) || []).join(' ');
        });

        if (opt[val]) {
            opt[val].forEach(cls => {
                if (cls.startsWith("flap-") || cls.startsWith("bg-")) $cont.addClass(cls);
                else $span.addClass(cls);
            });
        }
    }
    
    applyDynamicColors(colorClasses) {
        const $span = this.$el.find(".trc-char");
        const $cont = this.$el.find(".trc-container");
        
        // Entferne alte dynamische Farben
        $cont.removeClass(function (index, className) {
            return (className.match(/(^|\s)(flap-|bg-)\S+/g) || []).join(' ');
        });
        $span.removeClass(function (index, className) {
            return (className.match(/(^|\s)txt-\S+/g) || []).join(' ');
        });
        
        // Wende neue Farben an
        if (Array.isArray(colorClasses)) {
            colorClasses.forEach(cls => {
                if (cls.startsWith("flap-") || cls.startsWith("bg-")) $cont.addClass(cls);
                else $span.addClass(cls);
            });
        }
    }
}

let instances = [];

// Finde einen Character im Charset durch Substring-Matching
function findCharInCharset(chars, searchTerm) {
    if (!searchTerm) return null;
    
    // Exakte Übereinstimmung
    let idx = chars.indexOf(searchTerm);
    if (idx !== -1) return chars[idx];
    
    // Substring-Matching (case-insensitive)
    searchTerm = searchTerm.toLowerCase();
    for (let char of chars) {
        if (char.toLowerCase().includes(searchTerm)) {
            return char;
        }
    }
    return null;
}

function updateLayout() {
    const $img = $("#bg-image");
    if (!$img.length || $img.width() === 0) return;
    const bgDim = { w: $img.width(), h: $img.height() };
    instances.forEach(ins => ins.resize(bgDim));
}

function runStart() {
    const params = window.location.search.substring(1).split('&');

    // config wird global definiert, nicht lokal!
    
    const baseIds = ["Dir", "Track", "Train", "TrainNo_1", "TrainNo_2", "TrainNo_3", "TrainNo_4", "Destination", "Remark", "Hour", "Minute"];
    
    // Row 1
    baseIds.forEach((id, idx) => {
        if ($("#" + id).length && config[id]) {
            let ins;
            if (id === "Dir") {
                ins = new DirectionIndicator(id, config[id]);
            } else {
                ins = new SplitFlap(id, config[id]);
            }
            
            instances.push(ins);
        }
    });

    // Row 2
    baseIds.forEach((id, idx) => {
        const id2 = id + "_2"; 
        if ($("#" + id2).length && config[id]) {
            const row2Config = JSON.parse(JSON.stringify(config[id])); 
            row2Config.relTop += 36; // Offset Row 2
            if (id === "Dir" && row2Config.pos) {
                for (let key in row2Config.pos) {
                    if (row2Config.pos[key].relTop !== undefined) {
                        row2Config.pos[key].relTop += 36;
                    }
                }
            }

            let ins;
            if (id === "Dir") {
                ins = new DirectionIndicator(id2, row2Config);
            } else {
                ins = new SplitFlap(id2, row2Config);
            }
            
            instances.push(ins);
        }
    });

    updateLayout();
    $(window).on('resize', updateLayout);

    // Initialisiere alle Klappen mit Leerzeichen
    instances.forEach(ins => {
        ins.display(" ");
    });

    // Setze die Parameter NACH updateLayout() und nach initialer Display
    baseIds.forEach((id, idx) => {
        const instance = instances.find(ins => ins.id === id);
        if (instance && params[idx]) {
            const param = decodeURIComponent(params[idx]);
            
            // Substring-Matching für alle außer TrainNo (für TrainNo wird die Farb-Logik in der Flask-App gehandhabt)
            if (!id.startsWith("TrainNo")) {
                const charset = config[id] && config[id].flapCharset ? config[id].flapCharset.chars : [];
                const found = findCharInCharset(charset, param);
                instance.display(found || param);
            } else {
                instance.display(param);
            }
        }
    });
    
    // Row 2 Parameter
    baseIds.forEach((id, idx) => {
        const id2 = id + "_2";
        const instance = instances.find(ins => ins.id === id2);
        const paramIdx = idx + baseIds.length;
        if (instance && params[paramIdx]) {
            const param = decodeURIComponent(params[paramIdx]);
            
            if (!id.startsWith("TrainNo")) {
                const charset = config[id] && config[id].flapCharset ? config[id].flapCharset.chars : [];
                const found = findCharInCharset(charset, param);
                instance.display(found || param);
            } else {
                instance.display(param);
            }
        }
    });
}

$(window).on('load', runStart);

window.addEventListener("message", (event) => {
    const newParams = event.data;
    if (Array.isArray(newParams)) {
        const baseIds = ["Dir", "Track", "Train", "TrainNo_1", "TrainNo_2", "TrainNo_3", "TrainNo_4", "Destination", "Remark", "Hour", "Minute"];
        
        console.log("Message Handler received:", newParams);
        console.log("  TrainNo indices 3-6:", newParams[3], newParams[4], newParams[5], newParams[6]);
        
        // Row 1
        baseIds.forEach((id, idx) => {
            if (newParams[idx] !== undefined) {
                const instance = instances.find(ins => ins.id === id);
                if (instance) {
                    const param = decodeURIComponent(newParams[idx]);
                    console.log("  Calling display for " + id + ": " + param);
                    if (!id.startsWith("TrainNo")) {
                        const charset = config[id] && config[id].flapCharset ? config[id].flapCharset.chars : [];
                        const found = findCharInCharset(charset, param);
                        instance.display(found || param);
                    } else {
                        instance.display(param);
                    }
                }
            }
        });
        
        // Row 2
        baseIds.forEach((id, idx) => {
            const id2 = id + "_2";
            const paramIdx = idx + baseIds.length;
            if (newParams[paramIdx] !== undefined) {
                const instance = instances.find(ins => ins.id === id2);
                if (instance) {
                    const param = decodeURIComponent(newParams[paramIdx]);
                    if (!id.startsWith("TrainNo")) {
                        const charset = config[id] && config[id].flapCharset ? config[id].flapCharset.chars : [];
                        const found = findCharInCharset(charset, param);
                        instance.display(found || param);
                    } else {
                        instance.display(param);
                    }
                }
            }
        });
    }
}, false);