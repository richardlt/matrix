const HID = require('node-hid'),
  EventEmitter = require('events');

const configs = require('./configs.json');

module.exports = class Gamepad extends EventEmitter {
  constructor() {
    super();

    this._devices = [];
    this._mapConfigs = {};
    configs.forEach(config => {
      this._mapConfigs[config.vendorID + '_' + config.productID] = config;
    });
  }

  loadDevices() {
    this._devices = HID.devices().map(device => {
      return {
        device: device,
        states: {},
        config: this._mapConfigs[device.vendorId + '_' + device.productId]
      };
    }).filter(device => {
      return device.config;
    }).map(device => {
      device.hid = new HID.HID(device.device.path);
      return device;
    });
  }

  openDevices() {
    this._devices.forEach(device => {
      device.hid.on('data', this._handleData(device));
      device.hid.on('error', error => { console.log('error', error) });
    });
  }

  _handleData(device) {
    return data => {
      device.config.buttons.forEach(button => {
        const isPressed = (data[button.pin] & 0xff) == button.value;

        const currentState = device.states[button.name];

        if (isPressed && !currentState) {
          this.emit('event', button.name + ':press');
        } else if (!isPressed && currentState) {
          this.emit('event', button.name + ':release');
        }

        device.states[button.name] = isPressed;
      });
    };
  }
};