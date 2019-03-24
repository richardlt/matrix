export default function (state, action) {
  if (typeof state === 'undefined') {
    return {
      slots: [{}, {}, {}, {}],
      selectVisible: false
    }
  }

  switch (action.type) {
    case 'GET_SLOT':
      return Object.assign({}, state, {
        slots: state.slots.map((s, i) => { return { mine: action.slot == i }; }),
        selectVisible: false
      });
    case 'SET_SELECT_VISIBLE':
      return Object.assign({}, state, { selectVisible: action.visible });
  }

  return state
};