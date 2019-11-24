const messageLog = document.querySelector('#messageLog')

let webSocket = new WebSocket(`ws://${location.host}/websocket`)

webSocket.onopen = open => {
    console.log('Open:', open)
}

function calculateChecksum(messageBody) {
    let length = messageBody.length;
    let sum = Array.from(messageBody)
        .reduce((aggregation, character) => {
            aggregation += character.charCodeAt(0)
            return aggregation
        }, 0)
    let checksum = sum % length + 10

    console.log(`calculateChecksum:\n\tlength: ${length}\n\t   sum: ${sum}\n\t      = ${checksum}`)

    return checksum
}

webSocket.onmessage = message => {
    console.log('Message:', message)
    const messageElement = document.createElement('p')
    messageElement.innerText = `> ${message.data}`;
    messageLog.appendChild(messageElement)

    // if (message.data == 'l123: on') {
    //     document.querySelector('#l123').setAttribute('fill', 'green')
    // } else if (message.data + '' == 'l123: off') {
    //     document.querySelector('#l123').setAttribute('fill', 'red')
    // }

    const planObject = document.querySelector('#plan object')
    const planObjectSVG = planObject.getSVGDocument()
    let [id, state] = message.data.split(':')
    let obj = planObjectSVG.getElementById(id)

    if (obj != null && obj != undefined) {
        const style = getComputedStyle(document.body);
        const green = style.getPropertyValue('--clr-green')
        const yellow = style.getPropertyValue('--clr-yellow')
        const red = style.getPropertyValue('--clr-red')

        let col = red;
        switch (state.trim()) {
            case 'go':
                col = green
                break;

            case 'warn':
                col = yellow
                break
            
            case 'stop':
                col = red
                break
        
            default:
                col = red
                break;
        }

        obj.setAttribute('fill', col)
    }
}

webSocket.onerror = error => {
    console.log('Error:', error)
}

webSocket.onclose = close => {
    console.log('Close:', close)
}