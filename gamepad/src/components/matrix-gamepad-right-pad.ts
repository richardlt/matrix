import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

@customElement('matrix-gamepad-right-pad')
export class MatrixGamepadRightPad extends LitElement {
  static styles = css`
    :host {
      position: relative;
      background-color: #222222;
      border-radius: 50%;
      width: 100%;
      padding-top: 100%;
    }

    div {
      position: absolute;
      width: 34%;
      height: 34%;
      border-radius: 50%;
      box-shadow: 1px 2px 15px 3px rgba(0, 0, 0, 0.3);
    }

    div:hover {
      cursor: pointer;
    }

    div:active {
      opacity: 0.7;
    }

    div.y {
      top: 33%;
      left: 5%;
      background-color: #27ae60;
    }

    div.x {
      top: 5%;
      left: 33%;
      background-color: #2980b9;
    }

    div.a {
      top: 33%;
      right: 5%;
      background-color: #c0392b;
    }

    div.b {
      bottom: 5%;
      left: 33%;
      background-color: #f1c40f;
    }
  `;

  render() {
    return html`
      <div class="x" @mousedown=${this.handleClick('x_down')} @mouseup=${this.handleClick('x_up')}></div>
      <div class="y" @mousedown=${this.handleClick('y_down')} @mouseup=${this.handleClick('y_up')}></div>
      <div class="a" @mousedown=${this.handleClick('a_down')} @mouseup=${this.handleClick('a_up')}></div>
      <div class="b" @mousedown=${this.handleClick('b_down')} @mouseup=${this.handleClick('b_up')}></div>
    `;
  }

  handleClick(cmd: string) {
    return (e: Event) => {
      this.dispatchEvent(new CustomEvent('command', { detail: cmd }));
    };
  }
}
