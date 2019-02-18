class ToggleSwitch extends HTMLElement {
    constructor() {
        super()
        const shadow = this.attachShadow({mode: 'open'});
        shadow.innerHTML = `
        <style>
            /* The switch - the box around the slider */
                .switch {
                position: relative;
                display: inline-block;
                --width: 30px;
                --height: 17px;
                width: var(--width, 60px);
                height: var(--height, 34px);
            }
            
            /* Hide default HTML checkbox */
            .switch input {
                opacity: 0;
                width: 0;
                height: 0;
            }
            
            /* The slider */
            .slider {
                position: absolute;
                cursor: pointer;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background-color: #ccc;
                -webkit-transition: .4s;
                transition: .4s;
            }
            
            .slider:before {
                position: absolute;
                content: "";
                --size: 1px;
                height: calc(var(--height, 26px) - var(--size));
                width: calc(var(--height, 26px) - var(--size));
                left: var(--size);
                bottom: var(--size);
                background-color: white;
                -webkit-transition: .4s;
                transition: .4s;
            }
            
            input:checked + .slider {
                background-color: var(--clr-primary, #2196F3);
            }
            
            input:focus + .slider {
                box-shadow: 0 0 1px var(--clr-primary, #2196F3);
            }
            
            input:checked + .slider:before {
                -webkit-transform: translateX(var(--height));
                -ms-transform: translateX(var(--height));
                transform: translateX(var(--height));
            }
            
            /* Rounded sliders */
            .slider.round {
                border-radius: var(--height);
            }
            
            .slider.round:before {
                border-radius: 50%;
            } 
        </style>
        <label class="switch">
            <input type="checkbox">
            <span class="slider round"></span>
        </label>
        `
        this.shadowRoot.querySelector('input')
            .addEventListener('change', event => {
                this.shadowRoot.CustomEvent('change', event)
            })
    }
}
customElements.define('toggle-switch', ToggleSwitch)