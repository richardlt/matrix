export default function (state, action) {
  if (typeof state === 'undefined') {
    return {
      frame: []
    }
  }

  switch (action.type) {
    case 'NEW_FRAME':
      return Object.assign({}, state, {
        frame: action.frame
      });
  }

  return state
};