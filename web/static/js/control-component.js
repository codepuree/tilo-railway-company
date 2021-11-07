class TRCControl extends HTMLElement {
    static get name() { return 'trc-control' }

	constructor() {
		console.log('Initializing TRCControl')
		self = super()
        console.log('TRCControl value:', this.getAttribute('value'))
		this.attachShadow({ mode: 'open' })

		self.controlSlider = document.createElement('object')
		self.controlSlider.setAttribute('data', '/static/svg/slider.svg')
		self.controlSlider.setAttribute('type', 'image/svg+xml')
		self.controlSlider.style.height = '100%'
		self.controlSlider.style.width  = '100%'
		self.controlSlider.style.overflow = 'hidden'

		this.shadowRoot.append(self.controlSlider)

        self.controlSlider.addEventListener('load', _ => {
            // Constants
            this.container      = this.controlSlider.contentDocument.querySelector('svg')
			this.slider         = this.container.querySelector('#Slider')
			const btnSlider     = this.slider.querySelector('#SliderButton')
			this.txtSlider      = this.slider.querySelector('#SliderSpeed')
			const sliderGlider  = this.container.querySelector('#SliderGlider')
            this.barActualSpeed = this.container.querySelector('#barActualSpeed')
			const btnEmergency  = this.container.querySelector('#EmergencyButton')
            this.arduino        = this.container.querySelector('#Arduino')
            this.raspberry      = this.container.querySelector('#Raspberry')

            // Variables
            this.isDragging  = false
            this.intStartPos = this.container.createSVGPoint()
            this.offsetY = 0
            this.y = 0
            this.speed = 0
            this.fullPosY = -sliderGlider.getBBox().y
            this.zeroPosY = btnSlider.getBBox().y + btnSlider.getBBox().height / 2
            this.offPosY  = 86  // TODO: find way to calculate value

            // Add styling
            const style = document.createElement('style')
            style.innerText = `
                #SliderButton {
                    cursor: pointer;
                    user-select: none;
                }

                #EmergencyButton {
                    cursor: pointer;
                }

                #Arduino > .st0 {
                    fill: #e74215;
                }

                #Arduino[active] > .st0 {
                    fill: #49ad3c;
                }

                #Raspberry > .st0 {
                    fill: #e74215;
                }

                #Raspberry[active] > .st0 {
                    fill: #49ad3c;
                }
            `;
            this.container.appendChild(style)

            // Configure the slider
            this.slider.setAttribute('draggable', true)

            // Mouse
            btnSlider.addEventListener('mousedown',  this.#startDrag.bind(this))
            btnSlider.addEventListener('mousemove',  this.#drag.bind(this))
            btnSlider.addEventListener('mouseup',    this.#endDrag.bind(this))
            btnSlider.addEventListener('mouseleave', this.#endDrag.bind(this))

            // Touch
            btnSlider.addEventListener('touchstart',  this.#startDrag.bind(this))
            btnSlider.addEventListener('touchmove',   this.#drag.bind(this))
            btnSlider.addEventListener('touchend',    this.#endDrag.bind(this))
            btnSlider.addEventListener('touchleave',  this.#endDrag.bind(this))
            btnSlider.addEventListener('touchcancel', this.#endDrag.bind(this))

            // Emergency Button
            btnEmergency.addEventListener('click', _ => {
                this.dispatchEvent(new CustomEvent('emergency', { }))
            })

            // Get & set default values from attributes
            this.#setSpeed(parseInt(this.getAttribute('value')))
            this.#setBarActualSpeed(parseInt(this.getAttribute('actual-speed')))
            this.#setArduinoConnection(this.getAttribute('is-arduino-connected') == 'true')
            this.#setRaspberryConnection(this.getAttribute('is-raspberry-connected') == 'true')
        })
	}

    get value() {
        return this.speed
    }

    set value(value) {
        this.#setSpeed(value)
    }

    set actualSpeed(speed) {
        if (speed > 100) {
            speed = 100
            console.warn(`The given actual speed '${speed}' is greater than 100.`)
        } else if (speed < 0) {
            speed = 0
            console.warn(`The given actual speed '${speed}' is lower than 0.`)
        }

        this.#setBarActualSpeed(speed)
    }

    set isArduinoConnected(isConnected) {
        self.#setArduinoConnection(isConnected)
    }

    set isRaspberryConnected(isConnected) {
        self.#setRaspberryConnection(isConnected)
    }

    #startDrag(event) {
        if (!this.isDragging) {
            if (event.isCancellable) event.preventDefault()
            event.stopPropagation()
            this.isDragging = true
            this.intStartPos = this.#getPosition(event)
            this.offsetY = this.intStartPos.y + this.y - this.intStartPos.y
        }
    }

    #drag(event) {
        if (this.isDragging) {
            const intCurrPos = this.#getPosition(event)

            const tempY = intCurrPos.y - this.intStartPos.y + this.offsetY
            if (tempY >= this.offPosY) {
                this.y = this.offPosY
            } else if (tempY <= this.fullPosY) {
                this.y = this.fullPosY
            } else {
                this.y = tempY
            }
            const tempSpeed = this.#getSpeed(this.y)
            if (tempSpeed != this.speed) {
                this.speed = tempSpeed
                if (this.speed >= 0) {
                    this.txtSlider.innerHTML = this.speed
                }
                this.dispatchEvent(new CustomEvent('input', { detail: { value: this.speed } }))
            }
            this.slider.setAttributeNS(null, 'transform', `translate(0, ${this.y})`)
        }
    }

    #endDrag(event) {
        if (this.isDragging) {
            this.isDragging = false
            
            if (Math.abs(this.y) <= 5) {
                this.y = 0
                this.slider.setAttributeNS(null, 'transform', `translate(0, ${this.y})`)
            }

            this.dispatchEvent(new CustomEvent('change', { detail: { value: this.speed } }))
        }
    }

    #getPosition(event) {
        const point = this.container.createSVGPoint()

        point.x = event.touches ? event.touches[0].clientX : event.clientX
        point.y = event.touches ? event.touches[0].clientY : event.clientY

        return point.matrixTransform(this.container.getScreenCTM().inverse())
    }

    #getSpeed(y) {
        if (y < 0) {
            return Math.round(y / this.fullPosY * 100)
        } else {
            return - Math.round(y / this.offPosY * 100)
        }
    }
     
    #setSpeed(value) {
        this.y = value / 100 * this.fullPosY
        if (this.slider && this.txtSlider) {
            this.slider.setAttributeNS(null, 'transform', `translate(0, ${this.y})`)
            this.txtSlider.innerHTML = value
            // Set bar position
            const ctm = this.slider.getCTM()

            if (value >= 0) {
                this.txtSlider.innerHTML = value
            }
        }
    }

    #setArduinoConnection(isConnected) {
        if (isConnected) {
            this.arduino.setAttribute('active', true)
        } else {
            this.arduino.removeAttribute('active')
        }
    }

    #setRaspberryConnection(isConnected) {
        if (isConnected) {
            this.raspberry.setAttribute('active', true)
        } else {
            this.raspberry.removeAttribute('active')
        }
    }

    #setBarActualSpeed(speed) {
        // max height 155 from 100 speed mark
        // TODO: find a way to retrieve the 155 e.g. from sliderGlider
        const fullHeight = this.barActualSpeed.points[1].y - 155
        const height = fullHeight * speed / 100

        this.barActualSpeed.points[0].y = this.barActualSpeed.points[1].y - height
        this.barActualSpeed.points[3].y = this.barActualSpeed.points[2].y - height
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue != newValue) {
            switch (name) {
                case 'value':
                    this.#setSpeed(newValue)
                    break;

                case 'actual-speed':
                    this.#setBarActualSpeed(newValue)
                    break;

                case 'is-arduino-connected':
                    this.#setArduinoConnection(newValue)
                    break;

                case 'is-raspberry-connected':
                    this.#setRaspberryConnection(newValue)
                    break;

                default:
                    console.log(`The value for '${name}' changed from '${oldValue}' to '${newValue}'`);
            }
        }
    }

	static get observedAttributes() { return ['value', 'is-arduino-connected', 'is-raspberry-connected']; }
}
window.customElements.define(TRCControl.name, TRCControl)
