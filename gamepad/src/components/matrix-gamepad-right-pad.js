import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

class MatrixGamepadRightPad extends PolymerElement {
  static get template() {
    return html`
      <style>
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
      </style>

      <div id="x" class="x"></div>
      <div id="y" class="y"></div>
      <div id="a" class="a"></div>
      <div id="b" class="b"></div>
    `;
  }

  static get properties() { return {}; }

  ready() {
    super.ready();
    this.$.x.addEventListener('mousedown', this._handleClick('x_down'));
    this.$.x.addEventListener('mouseup', this._handleClick('x_up'));
    this.$.y.addEventListener('mousedown', this._handleClick('y_down'));
    this.$.y.addEventListener('mouseup', this._handleClick('y_up'));
    this.$.a.addEventListener('mousedown', this._handleClick('a_down'));
    this.$.a.addEventListener('mouseup', this._handleClick('a_up'));
    this.$.b.addEventListener('mousedown', this._handleClick('b_down'));
    this.$.b.addEventListener('mouseup', this._handleClick('b_up'));
  }

  _handleClick(cmd) {
    return event => {
      this.dispatchEvent(new CustomEvent('command', { detail: cmd }));
    };
  }
}

customElements.define('matrix-gamepad-right-pad', MatrixGamepadRightPad);
