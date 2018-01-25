import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';

const Square = styled.div`
  flex: 1;
  box-shadow: -1px 2px 10px 3px rgba(0, 0, 0, 0.5) inset;
`;

class Pixel extends React.Component {
  render() {
    return <Square style={{
      backgroundColor: 'rgb(' + this.props.color.r + ',' + this.props.color.g + ',' + this.props.color.b + ')'
    }}>&nbsp;</Square>;
  }
}

Pixel.propTypes = {
  color: PropTypes.objectOf(PropTypes.number)
}

export default Pixel;