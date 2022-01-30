import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { PlayerState } from '../store/player';

import store from '../store/store';
import { Slot } from '../types';

@customElement('matrix-gamepad-select-slots')
export class MatrixGamepadSelectSlots extends LitElement {
  @property({ type: Array }) slots: Array<Slot> = [];

  static styles = css`
    :host {
      position: absolute;
      height: 100%;
      width: 100%;
      display: flex;
      justify-content: center;
    }

    .row {
      flex: 0.5;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }

    .column {
      background-color: #111111;
      border: #333333 10px solid;
      text-align: center;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }

    .container {
      display: flex;
      flex-direction: row;
      justify-content: center;
    }

    .slot {
      position: relative;
      width: 22%;
      padding-top: 22%;
      margin-left: 2%;
      margin-top: 2%;
      margin-bottom: 2%;
      height: 0;
      cursor: pointer;
    }

    .container> :first-child {
      margin-left: 0;
    }
  `;

  render() {
    return html`      
      <div class="row">
        <div class="column">
          <div class="container">
            ${this.slots.map((slot, index) => html`
              <div class="slot" value="${index}" style="background-color: ${slot.mine ? 'green' : '#222222'}" @click="${this.handleClickSlot(index)}"></div>
            `)}
          </div>
        </div>
      </div>
    `;
  }

  connectedCallback() {
    super.connectedCallback();
    store.subscribe(() => {
      this.slots = (store.getState().player as PlayerState).slots;
    });
  }

  handleClickSlot(index: number) {
    return (e: Event) => {
      if (!this.slots[index].mine) {
        this.dispatchEvent(new CustomEvent('select-slot', { detail: index }));
      }
    };
  }
}
