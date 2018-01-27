const SerialPort = require('serialport');

const ARDUINO_NANO_VENDOR_ID = '1a86',
    ARDUINO_NANO_PRODUCT_ID = '7523',
    ARDUINO_UNO_VENDOR_ID = '2a03',
    ARDUINO_UNO_PRODUCT_ID = '0043';

module.exports = class Matrix {
    constructor(size) {
        this._size = size;
        this._brightness = 255;
        this._frame = [];
        for (let i = 0; i < this._size; i++) { this._frame.push(0, 0, 0); }
        this._ports = [];
        this._buffer = [];
    }

    async loadPorts() {
        (await SerialPort.list()).filter(port => {
            return (port.vendorId == ARDUINO_NANO_VENDOR_ID &&
                port.productId == ARDUINO_NANO_PRODUCT_ID) ||
                (port.vendorId == ARDUINO_UNO_VENDOR_ID &&
                    port.productId == ARDUINO_UNO_PRODUCT_ID);
        }).forEach(meta => {
            this._ports.push(new SerialPort(meta.comName, { autoOpen: false, baudRate: 115200 }));
        });
    }

    openPorts() {
        this._ports.forEach(port => {
            let interval = null;
            let lastBuffer = [];

            port.on('open', () => { console.log('open', port.path); });

            port.on('data', () => {
                interval = setInterval(_ => {
                    const buffer = this._buffer;
                    if (!this._equalArray(buffer, lastBuffer)) {
                        lastBuffer = buffer;
                        port.write(buffer);
                    }
                }, 45);
            });

            port.on('close', err => {
                console.log('close', port.path);
                clearInterval(interval);
            });

            port.on('error', err => { console.log('error', err); });

            port.open(err => { if (err) { console.log('opening', err); } });
        });
        return this;
    }

    setFrame(frame) { this._frame = frame; this._updateBuffer(frame); }

    setBrightness(value) { this._brightness = value; this._updateBuffer(this._frame); }

    _updateBuffer(frame) {
        const buffer = [this._brightness];
        for (let i = 0; i < this._size; i++) {
            const color = frame[i];
            if (color.a > 0) {
                buffer.push(color.r, color.g, color.b);
            } else {
                buffer.push(0, 0, 0);
            }
        }
        this._buffer = buffer;
    }

    _equalArray(array1, array2) {
        return (array1.length == array2.length && array1.every((u, i) => {
            return u == array2[i];
        }));
    }
};