import { LitElement } from 'lit';
import { customElement } from 'lit/decorators.js';
import { PlayerAction } from './store/player';

import store from './store/store';

@customElement('matrix-gamepad-socket')
export class MatrixGamepadSocket extends LitElement {
  ws: any;
  interval: any;

  constructor() {
    super();
    this.getSlot = this.getSlot.bind(this);
    this.newFrame = this.newFrame.bind(this);
    this.connect = this.connect.bind(this);
  }

  command(cmd: string) {
    if (!this.ws) { return; }
    this.ws.send(JSON.stringify({ type: 'command', data: cmd }));
  }

  selectSlot(slot: number) {
    if (!this.ws) { return; }
    this.ws.send(JSON.stringify({ type: 'select-slot', data: slot }));
  }

  connectedCallback() {
    super.connectedCallback();
    this.connect();
    this.interval = setInterval(this.connect, 1000);
  }

  connect() {
    if (this.ws) { return; }
    this.ws = new WebSocket(`ws://${window.location.host}/websocket`);
    this.ws.onopen = (event: any) => { this.selectSlot(0); };
    this.ws.onclose = (event: any) => { this.ws = null; };
    this.ws.onerror = (error: any) => { this.ws.close(); };
    this.ws.onmessage = (event: any) => {
      const message = JSON.parse(event.data);
      switch (message.type) {
        case 'frame':
          this.newFrame(message.data);
          break;
        case 'slot':
          this.getSlot(message.data);
          break;
      }
    };
  }

  getSlot(slot: number) {
    store.dispatch(<PlayerAction>{ type: 'GET_SLOT', slot: slot });
  }

  newFrame(data: any) {
    if (data.number == 0) {
      store.dispatch({ type: 'NEW_FRAME', frame: data.pixels });
    }
  }
}
