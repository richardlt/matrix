import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';

import Pixel from '../pixel/pixel';

const Line = styled.div`
  flex: 1;
  display: flex;
  flex-direction: row;
`;

const Container = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
`;

class Matrix extends React.Component {
  render() {
    let lines = [];
    for (let i = 0; i < this.props.height; i++) {
      let line = [];
      for (let j = 0; j < this.props.width; j++) {
        const index = i * this.props.width + j;
        const color = this.props.frame.length > index ? this.props.frame[index] : { r: 0, g: 0, b: 0 };
        line.push(<Pixel key={j} color={color} />);
      }
      lines.push(<Line key={i}>{line}</Line>);
    }
    return <Container>{lines}</Container>;
  }
}

Matrix.propTypes = {
  height: PropTypes.number,
  width: PropTypes.number,
  frame: PropTypes.arrayOf(PropTypes.objectOf(PropTypes.number))
}

export default Matrix;