import { LitElement, html, css } from 'lit';
import { customElement } from 'lit/decorators.js';

@customElement('matrix-gamepad-button')
export class MatrixGamepadButton extends LitElement {
  static styles = css`
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
  `;

  render() {
    return html`
      <slot></slot>
    `;
  }
}
