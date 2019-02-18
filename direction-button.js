class DirectionButton extends HTMLElement {
    constructor() {
        super()
        const shadow = this.attachShadow({ mode: 'open' });
        shadow.innerHTML = `<style>
            :host {
                width: 2rem;
                height: 2rem;
                display: inline-block;
                background-image: var(--img-url, url('./resources/richtung_links_inaktiv.svg'));
                background-size: contain;
                cursor: pointer;
            }
        </style>`

        this.addEventListener('click', event => {
            if (this.state === 'inaktiv') {
                this.state = 'aktiv'
            } else {
                this.state = 'inaktiv'
            }

            // this.shadowRoot.CustomEvent('change', event)

            var event = new CustomEvent("change", {
                detail: {
                    direction: this.direction === 'links' ? 'left' : 'right',
                    state: this.state == 'aktiv' ? 'active' : 'inactive'
                }
            });
            this.dispatchEvent(event);
        })
    }

    get state() {
        return this.getAttribute('state') ? this.getAttribute('state') : 'inaktiv';
    }

    set state(newState) {
        this.setAttribute('state', newState);
    }

    get direction() {
        return this.getAttribute('direction') ? this.getAttribute('direction') : 'links';
    }

    set direction(newDirection) {
        this.setAttribute('direction', newDirection);
    }

    static get observedAttributes() {
        return ['state', 'direction'];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        this.style.setProperty('--img-url', `url(./resources/richtung_${this.direction}_${this.state}.svg)`)
    }
}
customElements.define('direction-button', DirectionButton)