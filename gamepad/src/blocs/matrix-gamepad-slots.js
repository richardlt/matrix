import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import ReduxMixin from '../store/store.js';

class MatrixGamepadSlots extends ReduxMixin(PolymerElement) {
  static get template() {
    return html`
      <style>
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
      </style>
      
      <template is="dom-repeat" items="[[slots]]" as="s">
        <div class="slot" style\$="background-color: [[_getSlotColor(s)]]"></div>
      </template>
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
}

customElements.define('matrix-gamepad-slots', MatrixGamepadSlots);
