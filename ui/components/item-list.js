export function ItemList(model, container) {
    this.items = [];
    this.onUpdateItems = new Signal();
    this.onUpdateItems.add(this.updateItems.bind(this));
}