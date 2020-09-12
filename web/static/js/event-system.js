class EventSystem {
    constructor() {
        this.listeners = {}
    }

    listen(eventName, callback) {
        if (this.listeners[eventName] === undefined) {
            this.listeners[eventName] = []
        }

        this.listeners[eventName].push(callback)
    }

    throw(eventName, data) {
        if (this.listeners[eventName] !== undefined) {
            this.listeners[eventName]
                .forEach(callback => callback(data))
        } else {
            console.warn(`There is no listener for the event '${eventName}' defined.`)
        }
    }
}