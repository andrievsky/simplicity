import Signal from './signal.js';

// Use Promise???
const AppModel = function() {
    this.onUpdateItems = new Promise();

    this.updateItems = function(data) {
        this.onUpdateItems.resolve(data);
    }





}

export default AppModel;