class TRCMap extends HTMLElement {
    static get name() { return 'trc-map' }

    constructor() {
        self = super()
        this.attachShadow({ mode: 'open' })
        this.mapContainer = document.createElement('object')
        this.map = null
        this.mapContainer.addEventListener('load', this.#mapLoaded.bind(this))
        this.shadowRoot.append(this.mapContainer)
    }

    static set map(path) { this.#setMap(path) }

    #setMap(mapPath) {
        this.mapContainer.setAttribute('data', mapPath)
    }

    #mapLoaded() {
        this.map = this.mapContainer.contentDocument.querySelector('svg')
        const additionalStyle = document.createElement('style')
        additionalStyle.innerHTML = `
        .inactive {
            filter: invert(100%);
        }

        .interactive {
            cursor: pointer;
        }

        #speed {
            user-select: none;
        }
        `
        this.map.appendChild(additionalStyle)

        // Set up line focus button
        this.btnFocusLine = this.map.querySelector('#line')
        this.#svgAddClass(this.btnFocusLine, 'interactive')
        this.#svgAddClass(this.btnFocusLine, 'inactive')
        this.btnFocusLine.addEventListener('click', (event => {
            this.focus = !this.focus // TODO: Maybe it's better to send an event
            if (this.focus) {
                this.dispatchEvent(new CustomEvent('focus', { detail: { value: this.focus } }))
            } else {
                this.dispatchEvent(new CustomEvent('blur', { detail: { value: this.focus } }))
            }
        }).bind(this))


        // Set up direction west
        this.directionWest = this.map.querySelector('#dirw')
        if (!this.directionWest) {
            console.error(`Did not find direction west button for '${this.mapContainer.data}'`)
            return
        }
        this.#svgAddClass(this.directionWest, 'interactive')
        this.#svgAddClass(this.directionWest, 'inactive')
        this.directionWest.addEventListener('click', (event => {
            this.direction = 'w'
            this.dispatchEvent(new CustomEvent('change', { detail: { value: 'w' } }))
        }).bind(this))

        // Set up direction east
        this.directionEast = this.map.querySelector('#diro')
        if (!this.directionEast) {
            console.error(`Did not find direction ost button for '${this.mapContainer.data}'`)
            return
        }
        this.#svgAddClass(this.directionEast, 'interactive')
        this.#svgAddClass(this.directionEast, 'inactive')
        this.directionEast.addEventListener('click', (event => {
            this.direction = 'e'
            this.dispatchEvent(new CustomEvent('change', { detail: { value: 'e' } }))
        }).bind(this))
    }

    #svgAddClass(elem, className) {
        const classes = elem.getAttribute('class')?.split(' ') || []
        classes.push(className)
        elem.setAttribute('class', classes.join(' '))
    }

    #svgRemoveClass(elem, className) {
        const classes = elem.getAttribute('class').split(' ')
        elem.setAttribute('class', classes.filter(c => c != className).join(' '))
    }

    #svgToggleClass(elem, className) {
        const classes = elem.getAttribute('class').split(' ')
        if (classes.includes(className)) {
            this.#svgRemoveClass(elem, className)
        } else {
            this.#svgAddClass(elem, className)
        }
    }

    set speed(speed) {
        console.info('Set map speed:', this.map, self.map)
        this.map.querySelector('#speed').innerHTML = speed
    }

    set direction(direction) {
        if (direction[0] == 'w' || direction == 'west') {
            this.#svgRemoveClass(this.directionWest, 'inactive')
            this.#svgAddClass(this.directionEast, 'inactive')
        } else if (direction[0] == 'o' || direction == 'ost' || direction[0] == 'e' || direction == 'east') {
            this.#svgRemoveClass(this.directionEast, 'inactive')
            this.#svgAddClass(this.directionWest, 'inactive')
        } else {
            console.warn(`The direction '${direction}' is unknown.`)
        }
    }

    set focus(focus) {
        if (focus) {
            this.#svgRemoveClass(this.btnFocusLine, 'inactive')
        } else {
            this.#svgAddClass(this.btnFocusLine, 'inactive')
        }
    }

    static get observedAttributes() {
        return ['src', 'speed']
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue != newValue) {
            switch (name) {
                case 'src':
                    self.#setMap(newValue)
                    break

                case 'speed':
                    this.speed = newValue
                    break

                case 'focus':
                    // self.#handleFocus()
                    break

                default:
                    console.warn(`The value '${name}' that changes from '${oldValue}' to '${newValue}' is unknown to the system!`)
            }
        }
    }
}

window.customElements.define(TRCMap.name, TRCMap)
