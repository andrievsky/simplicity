import {Signal} from './signal.js';

export function Model() {
    const items = Signal([]);
    const selectedItem = Signal(null);


    return {items, selectedItem};
}