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
addScriptButton("Stop Fast", async function simple() {
    await webSocket.send("adszbdszcdszddszedszfdsz")
})
addScriptButton("Gleis 1 Weiche", async function simple() {
    await signalTower.setSwitch("e", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("f", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("a", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("b", "0")
})
addScriptButton("Gleis 2 Weiche", async function simple() {
    await signalTower.setSwitch("e", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("f", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("d", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("a", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("b", "1")
})
addScriptButton("Gleis 3 Weiche", async function simple() {
    await signalTower.setSwitch("d", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("e", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("f", "0")
    await signalTower.delay(50)
    await signalTower.setSwitch("a", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("c", "0")
})
addScriptButton("Gleis 4 Weiche", async function simple() {
    await signalTower.setSwitch("f", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("a", "1")
    await signalTower.delay(50)
    await signalTower.setSwitch("c", "1")
})
addScriptButton("Strecke", async function simple() {
    await signalTower.setBlock("g", "f", "30")
    await signalTower.delay(50)
    await signalTower.setBlock("f", "f", "30")
})
addScriptButton("Gleis 1", async function simple() {
    await signalTower.setBlock("a", "f", "30")
})
addScriptButton("Gleis 2", async function simple() {
    await signalTower.setBlock("b", "f", "30")
})
addScriptButton("Gleis 3", async function simple() {
    await signalTower.setBlock("c", "f", "30")
})
addScriptButton("Gleis 4", async function simple() {
    await signalTower.setBlock("d", "f", "30")
})
addScriptButton("Stop", async function simple() {
    signalTower.setBlockDirection("b", "s")
    signalTower.setBlockDirection("a", "s")
    signalTower.setBlockDirection("c", "s")
    signalTower.setBlockDirection("d", "s")
    signalTower.setBlockDirection("e", "s")
    signalTower.setBlockDirection("f", "s")
})



