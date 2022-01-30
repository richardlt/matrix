import { createStore, combineReducers } from 'redux';
import display from './display';
import player from './player';

export default createStore(combineReducers({
  display,
  player
}));