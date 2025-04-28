import {cloneTemplate, updateTemplate} from "../template.js";
import {CollectionSignal, Signal} from "../signal.js";
import {ItemEditFormUploadFiles} from "./item-edit-form-upload-files.js";
import {ItemEditFormImages} from "./item-edit-form-images.js";

export function ItemEditForm(item, model, service, templates) {
    const frag = cloneTemplate(templates["item-edit"]);
    const itemModel = new ItemEditFormModel(item);
    const previewList = frag.querySelector('.image-preview-list');

    updateTemplate(frag, item);
    const close = () => {
        itemModel.destroy();
        model.selectedItem.set(null);
    }
    const saveButton = frag.querySelector('.save');
    saveButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        updateItem();
    });
    const cancelButton = frag.querySelector('.cancel');
    cancelButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        discardUploadedImages().then(() => {
            close();
        });
    });

    const updateItem = () => {
        const id = itemModel.id;
        const item = itemModel.data();
        console.log("Updating item:", id, item);
        const result = service.updateItem(id, item);
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

    new ItemEditFormImages(previewList, service, itemModel);

    new ItemEditFormUploadFiles(frag, service, itemModel);

    const discardUploadedImages = async () => {
        const ids = itemModel.uploadedImages.get();

        const deletePromises = ids.map((id) =>
            service.deleteImage(id).then(r => {
                console.log("Deleting image:", id, r.ok() ? "OK" : "Failed");
                return r;
            })
        );

        await Promise.all(deletePromises);
    };

    return frag;
}

function ItemEditFormModel(item) {
    const id = item.id;
    const title = Signal(item.title || "Untitled");
    const description = Signal(item.description || "");
    const images = CollectionSignal(item.images || []);
    const tags = CollectionSignal(item.tags || []);
    const uploadedImages = CollectionSignal([]);

    const data = () => {
        const mergedImages = images.get().concat(uploadedImages.get());
        return {
            title: title.get(),
            description: description.get(),
            images: mergedImages,
            tags: tags.get(),
        };
    }

    const destroy = () => {
        title.unsubscribeAll();
        description.unsubscribeAll();
        images.unsubscribeAll();
        tags.unsubscribeAll();
        uploadedImages.unsubscribeAll();
    };

    return { id, title, description, images, tags, uploadedImages, data, destroy };
}