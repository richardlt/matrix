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
    const mapPins = {};
    device.config.pins.forEach(pin => { mapPins[pin.number] = pin; });

    return data => {
      // sort buttons by values and substract value from pin value to detect 
      // button state for multitouch pins or compare value for monotouch pins
      device.config.buttons.sort((b1, b2) => {
        return b2.value - b1.value;
      }).forEach(button => {
        const isPressed = mapPins[button.pin].multi ?
          (data[button.pin] - button.value) >= 0 : data[button.pin] == button.value;

        if (isPressed && mapPins[button.pin].multi) {
          data[button.pin] -= button.value;
        }

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