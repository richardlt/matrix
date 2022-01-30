import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

import store from './store/store';
import { PlayerAction, PlayerState } from './store/player';
import { MatrixGamepadSocket } from './socket';

import './blocs/matrix-gamepad-matrix';
import './components/matrix-gamepad-left-pad';
import './components/matrix-gamepad-right-pad';
import './blocs/matrix-gamepad-slots';
import './blocs/matrix-gamepad-select-slots';
import './components/matrix-gamepad-button';
import './socket';

@customElement('matrix-gamepad-app')
export class MatrixGamepadApp extends LitElement {
  @property({ type: Boolean }) selectVisible = false;

  static styles = css`
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
  `;

  render() {
    return html`      
      <div class="content">
        <div class="side">
          <matrix-gamepad-button class="top" @mousedown=${this.handleClick('l_down')} @mouseup=${this.handleClick('l_up')}>l</matrix-gamepad-button>
          <matrix-gamepad-left-pad @command="${this.handleCommand}"></matrix-gamepad-left-pad>
          <matrix-gamepad-button class="bottom" @mousedown=${this.handleClick('select_down')} @mouseup=${this.handleClick('select_up')}>select</matrix-gamepad-button>
        </div>
        <div class="center">
          <matrix-gamepad-matrix></matrix-gamepad-matrix>
          <matrix-gamepad-slots @click="${this.handleClickSlots}"></matrix-gamepad-slots>
        </div>
        <div class="side">
          <matrix-gamepad-button class="top" @mousedown=${this.handleClick('r_down')} @mouseup=${this.handleClick('r_up')}>r</matrix-gamepad-button>
          <matrix-gamepad-right-pad @command="${this.handleCommand}"></matrix-gamepad-right-pad>
          <matrix-gamepad-button class="bottom" @mousedown=${this.handleClick('start_down')} @mouseup=${this.handleClick('start_up')}>start</matrix-gamepad-button>
        </div>
      </div>
      <matrix-gamepad-select-slots style="display: ${this.selectVisible ? 'flex' : 'none'}" @select-slot="${this.handleSelectSlot}"></matrix-gamepad-select-slots>
      <matrix-gamepad-socket id="socket"></matrix-gamepad-socket>
    `;
  }

  connectedCallback() {
    super.connectedCallback();

    store.subscribe(() => {
      this.selectVisible = (store.getState().player as PlayerState).selectVisible;
    });

    document.addEventListener('keydown', this.handleKey('down'));
    document.addEventListener('keyup', this.handleKey('up'));
  }

  convertKeyToCmd(key: string) {
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

  handleKey(type: string) {
    return (e: KeyboardEvent) => {
      const cmd = this.convertKeyToCmd(e.key);
      if (cmd) {
        this.handleCommand(Object.assign({}, e, { detail: `${cmd}_${type}` }));
      }
    };
  }

  handleClick(cmd: string) {
    return (e: Event) => {
      this.handleCommand(Object.assign({}, e, { detail: cmd }));
    };
  }

  handleCommand(cmd: any) {
    (this.renderRoot.querySelector('#socket') as MatrixGamepadSocket).command(cmd.detail);
  }

  handleSelectSlot(slot: any) {
    (this.renderRoot.querySelector('#socket') as MatrixGamepadSocket).selectSlot(slot.detail);
  }

  handleClickSlots(e: Event) {
    store.dispatch(<PlayerAction>{ type: 'SET_SELECT_VISIBLE', visible: !this.selectVisible });
  }
}
