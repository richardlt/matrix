import React from 'react';
import { render } from 'react-dom';
import styled from 'styled-components';

import Matrix from './components/matrix';
import { Pixel } from './types';

const Container = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
`;

const Row = styled.div`
  height: 30%;
  width: 100%;
  display: flex;
  flex-direction: row;
`;

type State = {
  frame0: Array<Pixel>,
  frame1: Array<Pixel>,
  frame2: Array<Pixel>,
  frame3: Array<Pixel>,
  frame4: Array<Pixel>
};

class App extends React.Component<{}, State> {
  interval: any;
  ws: any;

  constructor(props: {}) {
    super(props);
    this.connect = this.connect.bind(this);
    this.state = {
      frame0: [],
      frame1: [],
      frame2: [],
      frame3: [],
      frame4: []
    };
  }

  componentDidMount() {
    this.connect();
    this.interval = setInterval(this.connect, 1000);
  }

  componentWillUnmount() {
    if (this.interval) { clearInterval(this.interval); }
    if (this.ws) { this.ws.close(); }
  }

  connect() {
    if (this.ws) { return; }
    this.ws = new WebSocket(`ws://${window.location.host}/websocket`);
    this.ws.onopen = (event: any) => { };
    this.ws.onclose = (event: any) => { this.ws = null; };
    this.ws.onerror = (error: any) => { this.ws.close(); };
    this.ws.onmessage = (event: any) => {
      const message = JSON.parse(event.data);
      if (message.type === "frame") {
        let state: any = {};
        state[`frame${message.data.number}`] = message.data.pixels;
        this.setState(state);
      }
    };
  }

  render() {
    return (
      <Container>
        <Matrix width={16} height={9} frame={this.state.frame0} />
        <Row>
          <Matrix width={16} height={9} frame={this.state.frame1} />
          <Matrix width={16} height={9} frame={this.state.frame2} />
          <Matrix width={16} height={9} frame={this.state.frame3} />
          <Matrix width={16} height={9} frame={this.state.frame4} />
        </Row>
      </Container>
    );
  }
}

render(<App />, document.getElementById('root'));