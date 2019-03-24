import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

class MatrixGamepadButton extends PolymerElement {
  static get template() {
    return html`
      <style>
        :host {
          border-radius: 15px;
          background-color: #222222;
          text-align: center;
          color: #555555;
          text-transform: uppercase;
          font-family: Arial, Helvetica, sans-serif;
          border: 5px solid #555555;
          font-weight: bold;
          box-shadow: 1px 2px 10px 3px rgba(0, 0, 0, 0.3);
        }

        :host(:active) {
          opacity: 0.7;
        }
      </style>

      <slot></slot>
    `;
  }

  static get properties() { return {}; }
}

customElements.define('matrix-gamepad-button', MatrixGamepadButton);
