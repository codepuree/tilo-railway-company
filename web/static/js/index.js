const eventSystem = new EventSystem()

// Setup direction buttons
const btnDirUpLeft = document.querySelector('#btnDirUpLeft')
const btnDirUpRight = document.querySelector('#btnDirUpRight')
const btnDirDownLeft = document.querySelector('#btnDirDownLeft')
const btnDirDownRight = document.querySelector('#btnDirDownRight')

eventSystem.listen('floor-direction', event => {
    if (event.floor === 'up') {
        if (event.direction === 'left') {
            btnDirUpLeft.state = 'aktiv'
            btnDirUpRight.state = 'inaktiv'
        } else if (event.direction === 'right') {
            btnDirUpLeft.state = 'inaktiv'
            btnDirUpRight.state = 'aktiv'
        } else {
            console.warn(`Unknown direction: '${event.direction}'`)
        }
    } else if (event.floor === 'down') {
        if (event.direction === 'left') {
            btnDirDownLeft.state = 'aktiv'
            btnDirDownRight.state = 'inaktiv'
        } else if (event.direction === 'right') {
            btnDirDownLeft.state = 'inaktiv'
            btnDirDownRight.state = 'aktiv'
        } else {
            console.warn(`Unknown direction: '${event.direction}'`)
        }
    } else {
        console.warn(`Unknown floor: '${event.floor}'`)
    }
})

btnDirUpLeft.addEventListener('change', event => {
    if (event.detail.state === 'active') {
        eventSystem.throw('floor-direction', { floor: 'up', direction: 'left' })
    } else {
        eventSystem.throw('floor-direction', { floor: 'up', direction: 'right' })
    }
})

btnDirUpRight.addEventListener('change', event => {
    if (event.detail.state === 'active') {
        eventSystem.throw('floor-direction', { floor: 'up', direction: 'right' })
    } else {
        eventSystem.throw('floor-direction', { floor: 'up', direction: 'left' })
    }
})

btnDirDownLeft.addEventListener('change', event => {
    if (event.detail.state === 'active') {
        eventSystem.throw('floor-direction', { floor: 'down', direction: 'left' })
    } else {
        eventSystem.throw('floor-direction', { floor: 'down', direction: 'right' })
    }
})

btnDirDownRight.addEventListener('change', event => {
    if (event.detail.state === 'active') {
        eventSystem.throw('floor-direction', { floor: 'down', direction: 'right' })
    } else {
        eventSystem.throw('floor-direction', { floor: 'down', direction: 'left' })
    }
})

// Modify SVGs
document.addEventListener('DOMContentLoaded', event => {
    const planObject = document.querySelector('#plan object')

    planObject.addEventListener('load', event => {
        const planObjectSVG = planObject.getSVGDocument()
    
        Array.from(planObjectSVG.querySelectorAll('circle'))
            .forEach(elem => {
                elem.style.cursor = 'pointer'
                // elem.setAttribute('r', '1.5rem')
    
                elem.addEventListener('click', event => {
                    const style = getComputedStyle(document.body);
                    const green = style.getPropertyValue('--clr-green')
                    const yellow = style.getPropertyValue('--clr-yellow')
                    const red = style.getPropertyValue('--clr-red')
    
                    if (elem.getAttribute('fill') === green) {
                        elem.setAttribute('fill', yellow)
                        webSocket.send(`${elem.id}: warn`)
                    } else if (elem.getAttribute('fill') === yellow) {
                        elem.setAttribute('fill', red)
                        webSocket.send(`${elem.id}: stop`)
                    } else {
                        elem.setAttribute('fill', green)
                        webSocket.send(`${elem.id}: go`)
                    }
                })
            })

        // Dragable arrows
        Array.from(planObjectSVG.querySelectorAll('.arrow'))
            .forEach(arrow => {
                const bbox = arrow.getBoundingClientRect()
                let isDown = false;

                eventSystem.listen('floor-direction', event => {
                    console.log('hide', arrow.classList, event)
                    if (arrow.classList.contains(event.direction)) {
                        arrow.style.display = 'none'
                    } else {
                        arrow.style.display = 'initial'
                    }
                })

                arrow.style.cursor = 'pointer'

                console.log(bbox)

                let down = event => {
                    isDown = true
                }
                arrow.addEventListener('mousedown', down)
                arrow.addEventListener('touchstart', down)

                let up = event => {
                    isDown = false
                }
                arrow.addEventListener('mouseup', up)
                arrow.addEventListener('touchend', up)
                arrow.addEventListener('touchcancel', up)

                let move = event => {
                    let clientX = event.touches ? event.touches[0].clientX : event.clientX
                    if (isDown) {
                        let offsetX = clientX - bbox.x

                        if (offsetX < 0 && arrow.classList.contains('right')) {
                            offsetX = 0
                        } else if (offsetX < -100 && arrow.classList.contains('left')) {
                            offsetX = 0
                        }

                        arrow.setAttribute('transform', `translate(${offsetX})`)
                    }
                }
                arrow.addEventListener('mousemove', move)
                arrow.addEventListener('touchmove', move)
            })
    })
})