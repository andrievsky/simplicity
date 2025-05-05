import {Signal} from './signal.js';

export function Model(service) {
    const items = Signal([]);
    const selectedItem = Signal(null);
    const refreshItems = () => {
        service.listItems().then((result) => {
            if (result.ok()) {
                items.set(result.data);
            } else {
                console.error("Error loading items:", result.error);
            }
        });
    }
    const openModal = (view) => {
        selectedItem.set(view);
    }

    const closeModal = () => {
        selectedItem.set(null);
    }

    const updateItem = (id, item) => {
        const result = service.updateItem(id, item);
        result.then((result) => {
            if (result.ok()) {
                console.log("Item updated successfully");
                refreshItems();
            } else {
                console.error("Error updating item:", result.error);
            }
        });
        return result;
    }

    const createItem = (item) => {
        const result = service.createItem(item);
        result.then((result) => {
            if (result.ok()) {
                console.log("Item created successfully");
                refreshItems();
            } else {
                console.error("Error creating item:", result.error);
            }
        });
        return result;
    }

    return { items, selectedItem , refreshItems, openModal, closeModal, updateItem, createItem };
}