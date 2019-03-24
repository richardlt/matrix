import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import ReduxMixin from './store/store.js';
import './blocs/matrix-gamepad-matrix.js';
import './components/matrix-gamepad-left-pad.js';
import './components/matrix-gamepad-right-pad.js';
import './blocs/matrix-gamepad-slots.js';
import './blocs/matrix-gamepad-select-slots.js';
import './components/matrix-gamepad-button.js';
import './socket.js';

class MatrixGamepadApp extends ReduxMixin(PolymerElement) {
  static get template() {
    return html`
      <style>
        :host {
          display: block;
          width: 100%;
          height: 100%;
        }

        .content {
          position: absolute;
          width: 100%;
          height: 100%;
          display: flex;
          flex-direction: row;
          align-items: center;
          background-color: #111111;
        }

        .side {
          flex: 0.25;
          display: flex;
          flex-direction: column;
          align-items: center;
          margin-left: 12px;
          margin-right: 12px;
        }

        .center {
          flex: 0.50;
          display: flex;
          flex-direction: column;
        }

        matrix-gamepad-matrix {
          border-radius: 5px;
          border: 5px solid #222222;
          width: calc(100% - 10px)!important;
          margin-bottom: 20px;
        }

        matrix-gamepad-left-pad,
        matrix-gamepad-right-pad {
          margin-top: 20px;
          margin-bottom: 20px;
        }

        matrix-gamepad-button.top {
          width: 80%;
          height: 30px;
          font-size: 25px;
          line-height: 34px;
        }

        matrix-gamepad-button.bottom {
          width: 40%;
          height: 19px;
          font-size: 10px;
          line-height: 22px;
        }
      </style>
      
      <div class="content">
        <div class="side">
          <matrix-gamepad-button id="l" class="top">l</matrix-gamepad-button>
          <matrix-gamepad-left-pad on-command="handleCommand"></matrix-gamepad-left-pad>
          <matrix-gamepad-button id="select" class="bottom">select</matrix-gamepad-button>
        </div>
        <div class="center">
          <matrix-gamepad-matrix></matrix-gamepad-matrix>
          <matrix-gamepad-slots on-click="handleClickSlots"></matrix-gamepad-slots>
        </div>
        <div class="side">
          <matrix-gamepad-button id="r" class="top">r</matrix-gamepad-button>
          <matrix-gamepad-right-pad on-command="handleCommand"></matrix-gamepad-right-pad>
          <matrix-gamepad-button id="start" class="bottom">start</matrix-gamepad-button>
        </div>
      </div>
      <matrix-gamepad-select-slots style\$="display: [[getSelectDisplay(selectVisible)]]" on-select-slot="handleSelectSlot"></matrix-gamepad-select-slots>
      <matrix-gamepad-socket id="socket"></matrix-gamepad-socket>
    `;
  }

  static get properties() {
    return {
      selectVisible: {
        type: Boolean,
        readOnly: true
      }
    };
  }

  static mapStateToProps(state) {
    return {
      selectVisible: state.player.selectVisible
    };
  }

  static mapDispatchToEvents(dispatch, element) {
    return {
      setSelectVisible: event => dispatch({
        type: 'SET_SELECT_VISIBLE',
        visible: event.detail
      })
    };
  }

  ready() {
    super.ready();

    this.$.l.addEventListener('mousedown', this._handleClick('l_down'));
    this.$.l.addEventListener('mouseup', this._handleClick('l_up'));
    this.$.select.addEventListener('pointerdown', this._handleClick('select_down'));
    this.$.select.addEventListener('mouseup', this._handleClick('select_up'));
    this.$.start.addEventListener('pointerdown', this._handleClick('start_down'));
    this.$.start.addEventListener('mouseup', this._handleClick('start_up'));
    this.$.r.addEventListener('mousedown', this._handleClick('r_down'));
    this.$.r.addEventListener('mouseup', this._handleClick('r_up'));

    document.addEventListener('keydown', this._handleKey('down'));
    document.addEventListener('keyup', this._handleKey('up'));
  }

  _convertKeyToCmd(key) {
    switch (key) {
      case 'z':
        return 'up';
      case 'd':
        return 'right';
      case 's':
        return 'down';
      case 'q':
        return 'left';
      case 'o':
        return 'x';
      case 'p':
        return 'a';
      case 'l':
        return 'b';
      case 'k':
        return 'y';
      case 'f':
        return 'l';
      case 'j':
        return 'r';
      case 'g':
        return 'select';
      case 'h':
        return 'start';
      default:
        return null;
    }
  }

  _handleKey(type) {
    return (e) => {
      const cmd = this._convertKeyToCmd(e.key);
      if (cmd) {
        this.handleCommand(Object.assign({}, e, { detail: cmd + '_' + type }));
      }
    };
  }

  _handleClick(cmd) {
    return (e) => {
      this.handleCommand(Object.assign({}, e, { detail: cmd }));
    };
  }

  handleCommand(cmd) { this.$.socket.command(cmd.detail); }

  handleSelectSlot(slot) { this.$.socket.selectSlot(slot.detail); }

  getSelectDisplay(visible) { return visible ? 'flex' : 'none'; }

  handleClickSlots(e) {
    this.dispatchEvent(new CustomEvent('set-select-visible', {
      detail: !this.selectVisible
    }));
  }
}

customElements.define('matrix-gamepad-app', MatrixGamepadApp);
