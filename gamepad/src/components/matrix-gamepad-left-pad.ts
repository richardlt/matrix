import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

@customElement('matrix-gamepad-left-pad')
export class MatrixGamepadLeftPad extends LitElement {
  static styles = css`
    :host {
      position: relative;
      background-color: #222222;
      border-radius: 50%;
      width: 100%;
      padding-top: 100%;
    }

    .button {
      position: absolute;
      box-sizing: border-box;
      width: 30%;
      height: 30%;
    }

    .button:active {
      opacity: 0.7;
    }

    .button.left {
      top: 35%;
      left: 5%;
    }

    .button.up {
      top: 5%;
      left: 35%;
    }

    .button.right {
      top: 35%;
      right: 5%;
    }

    .button.down {
      bottom: 5%;
      left: 35%;
    }

    .square {
      position: absolute;
      height: 100%;
      width: 100%;
      background-color: #444444;
    }

    .square:hover {
      cursor: pointer;
    }

    .square.left {
      border-top-left-radius: 25%;
      border-bottom-left-radius: 25%;
      box-shadow: 1px 2px 10px 3px rgba(0, 0, 0, 0.2);
    }

    .square.up {
      border-top-left-radius: 25%;
      border-top-right-radius: 25%;
      box-shadow: 1px -6px 12px 1px rgba(0, 0, 0, 0.1);
    }

    .square.right {
      border-top-right-radius: 25%;
      border-bottom-right-radius: 25%;
      box-shadow: 1px 2px 10px 3px rgba(0, 0, 0, 0.2);
    }

    .square.down {
      border-bottom-left-radius: 25%;
      border-bottom-right-radius: 25%;
      box-shadow: 1px 6px 12px 1px rgba(0, 0, 0, 0.2);
    }

    .arrow {
      position: absolute;
      transform: rotate(45deg);
      width: 72%;
      padding-top: 72%;
      background-color: #444444;
    }

    .arrow.left {
      left: 65%;
      top: 14.5%;
    }

    .arrow.up {
      left: 14.5%;
      top: 65%;
    }

    .arrow.right {
      right: 65%;
      top: 14.5%;
    }

    .arrow.down {
      left: 14.5%;
      bottom: 65%;
    }
  `;

  render() {
    return html`      
      <div class="button left" @mousedown=${this.handleClick('left_down')} @mouseup=${this.handleClick('left_up')}>
        <div class="square left"></div>
        <div class="arrow left"></div>
      </div>
      <div class="button right" @mousedown=${this.handleClick('right_down')} @mouseup=${this.handleClick('right_up')}>
        <div class="square right"></div>
        <div class="arrow right"></div>
      </div>
      <div class="button up" @mousedown=${this.handleClick('up_down')} @mouseup=${this.handleClick('up_up')}>
        <div class="square up"></div>
        <div class="arrow up"></div>
      </div>
      <div class="button down" @mousedown=${this.handleClick('down_down')} @mouseup=${this.handleClick('down_up')}>
        <div class="square down"></div>
        <div class="arrow down"></div>
      </div>
    `;
  }

  handleClick(cmd: string) {
    return (e: Event) => {
      this.dispatchEvent(new CustomEvent('command', { detail: cmd }));
    };
  }
}
