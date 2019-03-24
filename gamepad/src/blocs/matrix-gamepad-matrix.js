import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import ReduxMixin from '../store/store.js';

class MatrixGamepadMatrix extends ReduxMixin(PolymerElement) {
  static get template() {
    return html`
      <style>
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
      </style>
      <div class="matrix">
        <template is="dom-repeat" items="[[lines]]" as="line">
          <div class="line">
            <template is="dom-repeat" items="[[line]]" as="pixel">
              <div class="pixel" style\$="background-color: rgb([[pixel.r]],[[pixel.g]],[[pixel.b]])"></div>
            </template>
          </div>
        </template>
      </div>
    `;
  }

  static get properties() {
    return {
      pixels: {
        type: Array,
        observer: '_pixelsChanged',
        readOnly: true
      }
    };
  }

  static mapStateToProps(state) {
    return {
      pixels: state.display.frame
    };
  }

  static mapDispatchToEvents(dispatch, element) { return {}; }

  _pixelsChanged(newValue, oldValue) {
    let lines = [];
    for (let i = 0; i < 9; i++) {
      let line = [];
      for (let j = 0; j < 16; j++) {
        line.push(newValue[i * 16 + j])
      }
      lines.push(line);
    }
    this.lines = lines;
  }
}

customElements.define('matrix-gamepad-matrix', MatrixGamepadMatrix);
