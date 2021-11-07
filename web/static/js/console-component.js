class TRCConsole extends HTMLElement {
    static get name() { return 'trc-console' }

    constructor() {
        self = super()
        this.attachShadow({ mode: 'open' })
        const style = document.createElement('style')
        this.shadowRoot.append(style)
        self.content = document.createElement('dialog')
        this.shadowRoot.append(self.content)
    }

    printMessage(msg, type_='log') {
        if (msg == null || msg == undefined) {
            console.error(`The message '${msg}' is not defined.`)
        }
        console.log('message:', msg, type_)
        let msgElem = self.createMessageElement(msg, type_);
        const maxMessages = 100

        if (self.content.children.length == maxMessages) {
            self.content.firstChild.remove()
        }

        self.content.appendChild(msgElem)
        self.content.lastChild.scrollIntoView()
    }

    createMessageElement(msg, type_='log') {
        const message = document.createElement('p')
        message.innerHTML = `<b>${type_}:</b> ${msg}`
        console.log('message:', message.style.color)

        switch (type_) {
            case 'error':
                message.style.color = 'red'
                break;

            case 'warning':
                message.style.color = 'orange'
                break;

            case 'info':
                message.style.color = 'darkgray'
                break;

            case 'log':
            default:
                console.warn(`The type '${type_}' is unknown to 'trc-console'.`)
                break;
        }

        return message
    }

    open() {
        self.content.open()
    }
}

window.customElements.define(TRCConsole.name, TRCConsole)
