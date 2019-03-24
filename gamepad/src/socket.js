import { PolymerElement } from '@polymer/polymer/polymer-element.js';
import ReduxMixin from './store/store.js';

class MatrixGamepadSocket extends ReduxMixin(PolymerElement) {
  constructor() {
    super();
    this.getSlot = this.getSlot.bind(this);
    this.newFrame = this.newFrame.bind(this);
  }

  static get properties() { return {}; }

  static mapStateToProps(state) { return {}; }

  static mapDispatchToEvents(dispatch, element) {
    return {
      newFrame: event => dispatch({ type: 'NEW_FRAME', frame: event.detail }),
      getSlot: event => dispatch({ type: 'GET_SLOT', slot: event.detail })
    };
  }

  command(cmd) { this._socket.emit('command', cmd); }

  selectSlot(slot) { this._socket.emit('select-slot', slot); }

  ready() {
    super.ready();
    this._socket = io(window.origin);
    this._socket.on('connect', () => { this.selectSlot(0); });
    this._socket.on('frame', this.newFrame);
    this._socket.on('slot', this.getSlot);
    this._socket.on('disconnect', () => { });
  }

  getSlot(slot) {
    this.dispatchEvent(new CustomEvent('get-slot', {
      detail: slot
    }));
  }

  newFrame(data) {
    if (data.number == 0) {
      this.dispatchEvent(new CustomEvent('new-frame', {
        detail: data.pixels
      }));
    }
  }
}

customElements.define('matrix-gamepad-socket', MatrixGamepadSocket);
