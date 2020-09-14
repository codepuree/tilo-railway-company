addScriptButton("Simple", async function simple() {
    await signalTower.setBlock("a", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("e", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("f", "f", "50")

    await signalTower.waitForSensor("01", "l")
        .then(async (_) => {
            console.log("state switch")
            await signalTower.setBlockDirection("a", "s")
        })
})

addScriptButton("Gleis 1", async function simple() {
    await signalTower.setSwitch("a", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("b", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("e", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("f", "0")
    await signalTower.delay(200)
    await signalTower.setBlock("a", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("e", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("f", "f", "50")
})

addScriptButton("Gleis 2", async function simple() {
    await signalTower.setSwitch("a", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("b", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("d", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("e", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("f", "0")
    await signalTower.delay(200)
    await signalTower.setBlock("b", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("e", "f", "50")
    await signalTower.delay(50)
    await signalTower.setBlock("f", "f", "50")
})

addScriptButton("Stop", async function simple() {
    signalTower.setBlockDirection("b", "s")
    signalTower.setBlockDirection("a", "s")
    signalTower.setBlockDirection("c", "s")
    signalTower.setBlockDirection("d", "s")
    signalTower.setBlockDirection("e", "s")
    signalTower.setBlockDirection("f", "s")
})

addScriptButton("Stop Fast", async function simple() {
    await webSocket.send("adszbdszcdszddszedszfdsz")
})

