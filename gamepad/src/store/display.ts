import { Pixel } from "../types";

export class DisplayState {
  frame: Array<Pixel> = [];
}

export class DisplayAction {
  type: string = '';
  frame: Array<Pixel> = [];
}

export default function (state: DisplayState, action: DisplayAction): DisplayState {
  if (typeof state === 'undefined') {
    return {
      frame: []
    };
  }

  switch (action.type) {
    case 'NEW_FRAME':
      return {
        ...state,
        frame: action.frame
      };
  }

  return state;
};