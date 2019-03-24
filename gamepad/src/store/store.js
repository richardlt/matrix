import { createMixin } from 'polymer-redux';
import { createStore, combineReducers } from 'redux';
import display from './display';
import player from './player';

const store = createStore(combineReducers({
  display,
  player
}));

export default createMixin(store);
