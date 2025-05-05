import {CollectionSignal, Signal} from "../signal.js";
import {ItemFormUploadFiles} from "./item-form-upload-files.js";
import {ItemFormImages} from "./item-form-images.js";
import {ItemFormInput} from "./item-form-input.js";

export function ItemNewForm(model, service, templater) {
    const frag = templater.cloneTemplate("item-new-template");
    const itemModel = new ItemNewFormModel();
    const title = frag.querySelector('.title');
    const description = frag.querySelector('.description');
    const tags = frag.querySelector('.tags');
    const previewList = frag.querySelector('.image-preview-list');

    new ItemFormInput(title, itemModel.title)

    new ItemFormInput(description, itemModel.description)

    new ItemFormInput(tags, itemModel.tags)

    new ItemFormImages(previewList, service, itemModel.images);

    new ItemFormUploadFiles(frag, service, itemModel.images);

    const discardUploadedImages = async () => {
        const ids = itemModel.images.get();

        const deletePromises = ids.map((id) =>
            service.deleteImage(id).then(r => {
                console.log("Deleting image:", id, r.ok() ? "OK" : "Failed");
                return r;
            })
        );

        await Promise.all(deletePromises);
    };

    const close = () => {
        model.closeModal();
        itemModel.destroy();
    }
    const saveButton = frag.querySelector('.save');
    saveButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        model.createItem(itemModel.data()).then(() => {close()})
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

function ItemNewFormModel() {
    const title = Signal("New");
    const description = Signal("");
    const tags = Signal("");
    const images = CollectionSignal([]);

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
        tags.unsubscribeAll();
        images.unsubscribeAll();
    };

    return { title, description, tags, images, data, destroy };
}