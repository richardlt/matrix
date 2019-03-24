import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import ReduxMixin from '../store/store.js';

class MatrixGamepadSelectSlots extends ReduxMixin(PolymerElement) {
  static get template() {
    return html`
      <style>
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
      </style>
      
      <div class="row">
        <div class="column">
          <div class="container">
            <template is="dom-repeat" items="[[slots]]" as="s">
              <div class="slot" value="[[index]]" style\$="background-color: [[_getSlotColor(s)]]" on-click="handleClickSlot"></div>
            </template>
          </div>
        </div>
      </div>
    `;
  }

  static get properties() {
    return {
      slots: {
        type: Array,
        readOnly: true
      }
    };
  }

  static mapStateToProps(state) {
    return {
      slots: state.player.slots
    };
  }

  static mapDispatchToEvents(dispatch, element) { return {}; }

  _getSlotColor(slot) { return slot.mine ? "green" : "#222222"; }

  handleClickSlot(e) {
    const slot = e.target.value;
    if (!this.slots[slot].mine) {
      this.dispatchEvent(new CustomEvent('select-slot', { detail: slot }));
    }
  }
}

customElements.define('matrix-gamepad-select-slots', MatrixGamepadSelectSlots);
