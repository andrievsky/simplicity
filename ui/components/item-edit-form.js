import {CollectionSignal, Signal} from "../signal.js";
import {ItemFormUploadFiles} from "./item-form-upload-files.js";
import {ItemFormImages} from "./item-form-images.js";
import {ItemFormInput} from "./item-form-input.js";

export function ItemEditForm(item, model, service, templater) {
    const frag = templater.cloneTemplate("item-edit-template");
    const itemModel = new ItemEditFormModel(item);
    const title = frag.querySelector('.title');
    const description = frag.querySelector('.description');
    const tags = frag.querySelector('.tags');
    const previewList = frag.querySelector('.image-preview-list');

    new ItemFormInput(title, itemModel.title)

    new ItemFormInput(description, itemModel.description);

    new ItemFormInput(tags, itemModel.tags);

    new ItemFormImages(previewList, service, itemModel.images);

    new ItemFormUploadFiles(frag, service, itemModel.images);

    const discardUploadedImages = async () => {
        const ids = itemModel.images.get().filter(id => !item.images.includes(id));

        const deletePromises = ids.map((id) =>
            service.deleteImage(id).then(r => {
                console.log("Deleting image:", id, r.ok() ? "OK" : "Failed");
                return r;
            })
        );

        await Promise.all(deletePromises);
    };

    const updateItem = () => {
        const id = itemModel.id;
        const item = itemModel.data();
        model.updateItem(id, item).then(close);
    }
    const close = () => {
        model.closeModal();
        itemModel.destroy();
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

    return frag;
}

function ItemEditFormModel(item) {
    const id = item.id;
    const title = Signal(item.title || "Untitled");
    const description = Signal(item.description || "");
    const images = CollectionSignal(item.images || []);
    const tags = Signal(item.tags.toString() || "");

    const data = () => {
        return {
            title: title.get(),
            description: description.get(),
            images: images.get(),
            tags: tags.get().split(",").map(tag => tag.trim()).filter(tag => tag !== ""),
        };
    }

    const destroy = () => {
        title.unsubscribeAll();
        description.unsubscribeAll();
        images.unsubscribeAll();
        tags.unsubscribeAll();
    };

    return { id, title, description, images, tags, data, destroy };
}