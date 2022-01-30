import React from 'react';
import styled from 'styled-components';

import { Pixel } from '../types';

const Container = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
`;

const Line = styled.div`
  flex: 1;
  display: flex;
  flex-direction: row;
`;

const Square = styled.div`
  flex: 1;
  box-shadow: -1px 2px 10px 3px rgba(0, 0, 0, 0.5) inset;
`;

type Props = {
  height: number,
  width: number,
  frame: Array<Pixel>
};

export default class Matrix extends React.Component<Props> {
  render() {
    let lines = [];
    for (let i = 0; i < this.props.height; i++) {
      let line = [];
      for (let j = 0; j < this.props.width; j++) {
        const index = i * this.props.width + j;
        const color = this.props.frame.length > index ? this.props.frame[index] : { r: 0, g: 0, b: 0 };
        line.push(<Square key={j} style={{
          backgroundColor: `rgb(${color.r},${color.g},${color.b})`
        }}>&nbsp;</Square>);
      }
      lines.push(<Line key={i}>{line}</Line>);
    }
    return <Container>{lines}</Container>;
  }
}