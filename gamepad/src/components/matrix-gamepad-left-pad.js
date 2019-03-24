import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

class MatrixGamepadLeftPad extends PolymerElement {
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

        .button {
          position: absolute;
          box-sizing: border-box;
          width: 30%;
          height: 30%;
        }

        .button:active {
          opacity: 0.7;
        }

        .button.left {
          top: 35%;
          left: 5%;
        }

        .button.up {
          top: 5%;
          left: 35%;
        }

        .button.right {
          top: 35%;
          right: 5%;
        }

        .button.down {
          bottom: 5%;
          left: 35%;
        }

        .square {
          position: absolute;
          height: 100%;
          width: 100%;
          background-color: #444444;
        }

        .square:hover {
          cursor: pointer;
        }

        .square.left {
          border-top-left-radius: 25%;
          border-bottom-left-radius: 25%;
          box-shadow: 1px 2px 10px 3px rgba(0, 0, 0, 0.2);
        }

        .square.up {
          border-top-left-radius: 25%;
          border-top-right-radius: 25%;
          box-shadow: 1px -6px 12px 1px rgba(0, 0, 0, 0.1);
        }

        .square.right {
          border-top-right-radius: 25%;
          border-bottom-right-radius: 25%;
          box-shadow: 1px 2px 10px 3px rgba(0, 0, 0, 0.2);
        }

        .square.down {
          border-bottom-left-radius: 25%;
          border-bottom-right-radius: 25%;
          box-shadow: 1px 6px 12px 1px rgba(0, 0, 0, 0.2);
        }

        .arrow {
          position: absolute;
          transform: rotate(45deg);
          width: 72%;
          padding-top: 72%;
          background-color: #444444;
        }

        .arrow.left {
          left: 65%;
          top: 14.5%;
        }

        .arrow.up {
          left: 14.5%;
          top: 65%;
        }

        .arrow.right {
          right: 65%;
          top: 14.5%;
        }

        .arrow.down {
          left: 14.5%;
          bottom: 65%;
        }
      </style>
      
      <div id="left" class="button left">
        <div class="square left"></div>
        <div class="arrow left"></div>
      </div>
      <div id="right" class="button right">
        <div class="square right"></div>
        <div class="arrow right"></div>
      </div>
      <div id="up" class="button up">
        <div class="square up"></div>
        <div class="arrow up"></div>
      </div>
      <div id="down" class="button down">
        <div class="square down"></div>
        <div class="arrow down"></div>
      </div>
    `;
  }

  static get properties() { return {}; }

  ready() {
    super.ready();
    this.$.left.addEventListener('mousedown', this._handleClick('left_down'));
    this.$.left.addEventListener('mouseup', this._handleClick('left_up'));
    this.$.up.addEventListener('mousedown', this._handleClick('up_down'));
    this.$.up.addEventListener('mouseup', this._handleClick('up_up'));
    this.$.right.addEventListener('mousedown', this._handleClick('right_down'));
    this.$.right.addEventListener('mouseup', this._handleClick('right_up'));
    this.$.down.addEventListener('mousedown', this._handleClick('down_down'));
    this.$.down.addEventListener('mouseup', this._handleClick('down_up'));
  }

  _handleClick(cmd) {
    return event => {
      this.dispatchEvent(new CustomEvent('command', { detail: cmd }));
    };
  }
}

customElements.define('matrix-gamepad-left-pad', MatrixGamepadLeftPad);
