const { EventEmitter } = require('events');
const { WebSocket } = require('ws');
const { merge, cloneDeep, debounce } = require('lodash');

class ChatClient extends EventEmitter {
    constructor(url) {
        super();
        this.url = url;
        this.ws = null;
        this.messageQueue = [];
        this.connected = false;
    }

    connect() {
        this.ws = new WebSocket(this.url);

        this.ws.on('open', () => {
            this.connected = true;
            this.emit('connected');
            this.flushQueue();
        });

        this.ws.on('message', (data) => {
            const message = JSON.parse(data);
            this.emit('message', message);
        });

        this.ws.on('close', () => {
            this.connected = false;
            this.emit('disconnected');
        });

        this.ws.on('error', (error) => {
            this.emit('error', error);
        });
    }

    sendMessage(content) {
        const message = {
            content: content,
            timestamp: Date.now(),
            type: 'text'
        };

        if (this.connected) {
            this.ws.send(JSON.stringify(message));
        } else {
            this.messageQueue.push(message);
        }
    }

    flushQueue() {
        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            this.ws.send(JSON.stringify(message));
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

class MessageStore {
    constructor() {
        this.messages = [];
        this.users = new Map();
    }

    addMessage(message) {
        this.messages.push(cloneDeep(message));
        this.emit('message:added', message);
    }

    getRecent(limit = 50) {
        return this.messages.slice(-limit);
    }

    clear() {
        this.messages = [];
    }
}

function formatMessage(message) {
    const formatted = merge({}, message, {
        formatted_time: new Date(message.timestamp).toLocaleString()
    });
    return formatted;
}

const debouncedSend = debounce((client, message) => {
    client.sendMessage(message);
}, 300);

module.exports = {
    ChatClient,
    MessageStore,
    formatMessage,
    debouncedSend
};
