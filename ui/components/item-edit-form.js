import {cloneTemplate, updateTemplate} from "../template.js";
import {FormReader} from "../form-reader.js";

export function ItemEditForm(item, model, service, templates, close) {
    const frag = cloneTemplate(templates["item-edit"]);
    updateTemplate(frag, item);
    const element = frag.firstElementChild;
    const saveButton = frag.querySelector('.save');
    saveButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        updateItem(item.id, () => readItemData(element));
    });
    const cancelButton = frag.querySelector('.cancel');
    cancelButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        close();
    });

    const updateItem = (id, data) => {
        console.log("Updating item:", id, data());
        const result = service.updateItem(id, data());
        result.then((result) => {
            if (result.ok()) {
                console.log("Item updated successfully");
                refreshItems();
                close();
            } else {
                console.error("Error updating item:", result.error);
            }
        });
    }

    const refreshItems = () => {
        service.listItems().then((result) => {
            if (result.ok()) {
                model.items.set(result.data);
            } else {
                console.error("Error loading items:", result.error);
            }
        });
    }

    const readItemData = (element) => {
        const read = FormReader(element);
        return {
            title: read.string( 'title', 'Untitled'),
            description: read.string('description', ''),
            images: read.array('images', []),
            tags: read.array('tags', []),
        };
    };

    return frag;
}