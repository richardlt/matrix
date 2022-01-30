import { Slot } from '../types';

export class PlayerState {
  slots: Array<Slot> = [];
  selectVisible: boolean = false;
}

export class PlayerAction {
  type: string = '';
  slot: number = 0;
  visible: boolean = false;
}

export default function (state: PlayerState, action: PlayerAction): PlayerState {
  if (typeof state === 'undefined') {
    return {
      slots: [
        new Slot(false),
        new Slot(false),
        new Slot(false),
        new Slot(false)
      ],
      selectVisible: false
    };
  }

  switch (action.type) {
    case 'GET_SLOT':
      return {
        ...state,
        slots: state.slots.map((_, i) => new Slot(action.slot == i)),
        selectVisible: false
      };
    case 'SET_SELECT_VISIBLE':
      return {
        ...state,
        selectVisible: action.visible
      };
  }

  return state;
};