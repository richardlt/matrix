import React from 'react';
import { render } from 'react-dom';
import io from 'socket.io-client';
import styled from 'styled-components';

import Matrix from './components/matrix/matrix';

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

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      frame0: [],
      frame1: [],
      frame2: [],
      frame3: [],
      frame4: []
    };
  }

  componentDidMount() {
    const socket = io(window.origin);
    socket.on('connect', () => { });
    socket.on('frame', (data) => {
      let state = {};
      state['frame' + data.number] = data.pixels;
      this.setState(state);
    });
    socket.on('disconnect', () => { });
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