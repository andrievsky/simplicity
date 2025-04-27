import {cloneTemplate, updateTemplate} from "../template.js";
import {FormReader} from "../form-reader.js";
import {CollectionSignal, Signal} from "../signal.js";
import {ItemEditFormImage} from "./item-edit-form-image.js";

export function ItemEditForm(item, model, service, templates) {
    const frag = cloneTemplate(templates["item-edit"]);
    const itemModel = new ItemEditFormModel(item);

    updateTemplate(frag, item);
    const close = () => {
        model.selectedItem.set(null);
        itemModel.destroy();
    }
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

    const dropZone = frag.querySelector('.drop-zone');
    const input = frag.querySelector('.image-upload');
    const previewList = frag.querySelector('.image-preview-list');

    const uploadedImageViews = {};

    itemModel.uploadedImages.subscribeDelta((delta => {
        if (delta.type === "add") {
            const imageView = new ItemEditFormImage(delta.item, (id) => {
                itemModel.uploadedImages.remove(id);
            });
            uploadedImageViews[delta.item] = previewList.appendChild(imageView.wrapper);
        }
        if (delta.type === "remove") {
            const wrapper = uploadedImageViews[delta.item];
            if (wrapper) {
                wrapper.remove();
                delete uploadedImageViews[delta.item];
            }
        }
    }));

    dropZone.addEventListener("click", () => input.click());

    dropZone.addEventListener("dragover", (e) => {
        e.preventDefault();
        dropZone.classList.add("highlight");
    });
    dropZone.addEventListener("dragleave", () => {
        dropZone.classList.remove("highlight");
    });
    dropZone.addEventListener("drop", async (e) => {
        e.preventDefault();
        dropZone.classList.remove("highlight");
        await handleFiles(e.dataTransfer.files);
    });

    input.addEventListener("change", async (e) => {
        await handleFiles(e.target.files);
    });

    async function handleFiles(files) {
        for (const file of files) {
            const response = await service.uploadImage(file);
            if (response.ok()) {
                itemModel.uploadedImages.add(response.data.id);
            } else {
                console.error("Error uploading image:", response.error);
            }
        }
    }

    return frag;
}

function ItemEditFormModel(item) {
    const title = Signal(item.title || "Untitled");
    const description = Signal(item.description || "");
    const images = CollectionSignal(item.images || []);
    const tags = CollectionSignal(item.tags || []);
    const uploadedImages = CollectionSignal([]);

    const destroy = () => {
        title.unsubscribeAll();
        description.unsubscribeAll();
        images.unsubscribeAll();
        tags.unsubscribeAll();
        uploadedImages.unsubscribeAll();
    };

    return { title, description, images, tags, uploadedImages, destroy };
}