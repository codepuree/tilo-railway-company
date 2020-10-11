class SignalTower {
    /* Sensoren:
        00 l|h z -> high no Zug; low Zug
    */

    /* Blocks:
        a-k|l-p d f|b|s|x z  -> (f)orward; (b)ackward; (s)top; (x)nothalt
        a-k|l-p 00-99 z
    */

    /* Weichen:
        y a-h 0|1 z -> 0 gerade aus (großer radius); 1 abbiegen (kleiner radius)
    */

    /* Example codes:
        jf1z
        5dfz
        550z
        jc1z
    */

    sensors  = {}
    blocks   = {}
    switches = {}

    #sensorListeners = {}

    constructor() {

    }

    async waitForSensor(id, value) {
        if (this.#sensorListeners[id] == undefined) {
            this.#sensorListeners[id] = []
        }
        
        await new Promise((resolve, reject) => {
            let clb = state => {
                if (state == value) {
                    let i = this.#sensorListeners[id].findIndex(clb)
                    this.#sensorListeners.splice(i, 1)

                    resolve()
                }
            }
            this.#sensorListeners[id].push(clb)
        })
    }

    async setSwitch(id, direction) {
        webSocket.send(`y${id}${direction}z`)
    }

    async setBlock(id, direction, speed) {
        await this.setBlockDirection(id, direction)
        await this.setBlockSpeed(id, speed)
    }

    async setBlockDirection(id, direction) {
        webSocket.send(`${id}d${direction}z`)
    }

    async setBlockSpeed(id, speed) {
        webSocket.send(`${id}${speed}z`)
    }

    async delay(duration) {
        await new Promise((resolve, reject) => {
            setTimeout(_ => {
                resolve()
            }, duration)
        })
    }

    receiveMessage(msg) {
        if (Number.isInteger(parseInt(msg.charAt(0)))) {
            // Sensor
            let id = msg.substring(0, 2)
            let state = msg.charAt(2)

            this.sensors[id] = state
            console.log("seonsor state change")
            if (this.#sensorListeners[id] != undefined) {
                this.#sensorListeners[id].forEach(clb => clb(state))
            }
        } else if (msg.charAt(0) == "y") {
            // Weiche
            let id = msg.charAt(1)
            let state = msg.charAt(2)

            this.switches[id] = state
        } else {
            // Block
            let id = msg.charAt(0)

            if (this.blocks[id] == undefined) {
                this.blocks[id] = {}
            }

            if (msg.charAt(1) == "d") {
                let direction = msg.charAt(2)

                this.blocks[id].direction = direction
            } else {
                let speed = msg.substring(2, 4)

                this.blocks[id].speed = speed
            }
        }
    }
}