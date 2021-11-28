class TRCConsole extends HTMLElement {
    static get name() { return 'trc-console' }

    constructor() {
        self = super()
        this.attachShadow({ mode: 'open' })
        const style = document.createElement('style')
        style.innerHTML = `
            dialog {
                color: red;
                background-color: white;
                position: fixed;
                top: 50%;
                left: 50%;
                right: 50%;
                bottom: 50%;
                min-height: 80vh;
                max-height: 80vh;
                min-width: 30vw;
                max-width: 80vw;
                /*display: grid;*/
                grid-template-rows: 1fr auto 1fr;
            }

            dialog[open] {
                display: grid;
            }

            header {
                width: 100%;
                display: flex;
                justify-content: space-between;
            }

            main {
                overflow-y: auto;
                /* max-height: 60vh; */
            }

            footer { }

            #logLevels { }

            #logLevels>input[type=checkbox] {
                display: none;
            }

            #logLevels>label {
                border: 1px solid lightgray;
                border-radius: 2em;
                padding: 5px;
                text-align: center;
            }

            #logLevels>input[type=checkbox]:checked + label{
                border-color: red;
            }
        `
        this.shadowRoot.append(style)
        this.container = document.createElement('dialog')
        this.shadowRoot.append(this.container)

        // Header
        const header = document.createElement('header')
        this.container.append(header)

        // Log-Level Selection
        this.logLevels = document.createElement('form')
        this.logLevels.id = 'logLevels'
        header.append(this.logLevels)
        this.logLevels.append(this.#createLogLevel('ℹ️', 'info'))
        this.logLevels.append(this.#createLogLevel('⚠️', 'warning'))
        this.logLevels.append(this.#createLogLevel('🔥', 'error'))

        // Shutdown + Reboot
        const boots = document.createElement('form')
        header.append(boots)
        const btnShutdown = document.createElement('button')
        btnShutdown.title = 'Shutdown'
        btnShutdown.innerText = '🛑'
        btnShutdown.addEventListener('click', (event => {
            this.dispatchEvent(new CustomEvent('shutdown', {}))
        }).bind(this))
        boots.append(btnShutdown)

        const btnReboot = document.createElement('button')
        btnReboot.title = 'Reboot'
        btnReboot.innerText = '🔄'
        btnReboot.addEventListener('click', (event => {
            this.dispatchEvent(new CustomEvent('reboot', {}))
        }).bind(this))
        boots.append(btnReboot)

        // Main Content
        this.main = document.createElement('main')
        this.container.append(this.main)

        // Message Send
        this.footer = document.createElement('footer')
        const sendMessage = document.createElement('form')
        sendMessage.id = 'sendMessage'
        const inMessage = document.createElement('input')
        inMessage.name = 'message'
        const submitMessage = document.createElement('input')
        submitMessage.type = 'submit'
        sendMessage.addEventListener('submit', (event => {
            event.preventDefault()
            const message = inMessage.value
            if (message.trim().length == 0) {
                return
            }
            this.dispatchEvent(new CustomEvent('message', { detail: { message } }))
        }).bind(this))
        sendMessage.append(inMessage)
        sendMessage.append(submitMessage)
        this.footer.append(sendMessage)
        this.container.append(this.footer)
    }

    #createLogLevel(icon, name, isSelected=true) {
        const fragment = new DocumentFragment()
        
        const input = document.createElement('input')
        input.type = 'checkbox'
        input.name = 'logLevel'
        input.id = `logLevel_${name}`
        input.value = name
        input.checked = isSelected
        fragment.append(input)
        
        const label = document.createElement('label')
        label.setAttribute('for', `logLevel_${name}`)
        label.innerText = icon
        label.title = name
        fragment.append(label)

        return fragment
    }

    #getLogLevels() {
        const fd = new FormData(this.logLevels)
        return Array.from(fd.values())
    }

    printMessage(msg, type_='log') {
        if (msg == null || msg == undefined) {
            console.error(`Received an empty message with type ${type_}!`)
            return
        }

        if (!this.#getLogLevels().includes(type_.replace('log', 'info'))) {
            return
        }
        
        if (typeof msg == 'object') {
            msg = JSON.stringify(msg, null, 2)
        }

        let msgElem = self.createMessageElement(msg, type_);
        const maxMessages = 100  // TODO: allow `maxMessages` to be set from the attributes

        if (this.main.children.length == maxMessages) {
            this.main.firstChild.remove()
        }

        this.main.appendChild(msgElem)
        this.main.lastChild.scrollIntoView()
    }

    createMessageElement(msg, type_='log') {
        const message = document.createElement('p')
        let icon = 'ℹ️'
        let color = 'darkgray'

        switch (type_) {
            case 'error':
                icon = '🔥'
                color = 'red'
                break;

            case 'warning':
                icon = '⚠️'
                color = 'orange'
                break;

            case 'log':
            case 'info':
                icon = 'ℹ️'
                color = 'darkgray'
                break;

            default:
                console.warn(`The type '${type_}' is unknown to 'trc-console'.`)
                break;
        }

        message.innerHTML = `<b title="${type_}">${icon}</b> ${msg}`
        message.style.color = color

        return message
    }

    set open(isOpen) {
        if (isOpen) {
            this.container.setAttribute('open', isOpen)
        } else {
            this.container.removeAttribute('open')
        }
    }

    get open() {
        return false  // this.content.getAttribute('open') == 'true' ? true : false
    }
}

window.customElements.define(TRCConsole.name, TRCConsole)
