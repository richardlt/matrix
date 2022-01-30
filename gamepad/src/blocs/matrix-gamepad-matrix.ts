import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { DisplayState } from '../store/display';
import store from '../store/store';

import { Pixel } from '../types';

@customElement('matrix-gamepad-matrix')
export class MatrixGamepadMatrix extends LitElement {
  @property({ type: Array }) pixels: Array<Pixel> = [];

  static styles = css`
    :host {
      position: relative;
      width: 100%;
      padding-top: 57%;
    }

    .matrix {
      position: absolute;
      top: 0;
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;
    }

    .line {
      flex: 1;
      display: flex;
      flex-direction: row;
    }

    .pixel {
      flex: 1;
      box-shadow: -1px 2px 10px 3px rgba(0, 0, 0, 0.3) inset;
    }
  `;

  render() {
    let lines: Array<Array<Pixel>> = [];
    for (let i = 0; i < 9; i++) {
      let line: Array<Pixel> = [];
      for (let j = 0; j < 16; j++) {
        line.push(this.pixels[i * 16 + j] ?? new Pixel(0, 0, 0));
      }
      lines.push(line);
    }

    return html`
      <div class="matrix">
        ${lines.map(line => html`
          <div class="line">
            ${line.map(pixel => html`
              <div class="pixel" style="background-color: rgb(${pixel.r},${pixel.g},${pixel.b})"></div>
            `)}
          </div>        
        `)}
      </div>
    `;
  }

  connectedCallback() {
    super.connectedCallback();
    store.subscribe(() => {
      this.pixels = (store.getState().display as DisplayState).frame;
    });
  }
}
