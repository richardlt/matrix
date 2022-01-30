import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { PlayerState } from '../store/player';

import store from '../store/store';
import { Slot } from '../types';

@customElement('matrix-gamepad-slots')
export class MatrixGamepadSlots extends LitElement {
  @property({ type: Array }) slots: Array<Slot> = [];

  static styles = css`
    :host {
      height: 15px;
      display: block;
      text-align: center;
      cursor: pointer;
    }

    .slot {
      display: inline-block;
      height: 15px;
      width: 15px;
    }
  `;

  render() {
    return html`
      <div>
        ${this.slots.map(slot => html`
          <div class="slot" style="background-color: ${slot.mine ? 'green' : '#222222'}"></div>
        `)}
      </div>
    `;
  }

  connectedCallback() {
    super.connectedCallback();
    store.subscribe(() => {
      this.slots = (store.getState().player as PlayerState).slots;
    });
  }
}
