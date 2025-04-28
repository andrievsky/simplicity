import {ItemEditFormImage} from "./item-edit-form-image.js";

export function ItemEditFormImages(container, service, itemModel) {
    const uploadedImageViews = {};
    const subscriptions = [];

    const setPreviewHandler = (id) => {
        console.log("setPreviewHandler", id)
    }

    const addImageView = (id, removeHandler, setPreviewHandler) => {
        const imageView = new ItemEditFormImage(id, removeHandler, setPreviewHandler);
        uploadedImageViews[id] = container.appendChild(imageView.wrapper);
    }

    const addExistingImage = (id) => {
        console.log("addExistingImage", id);
        addImageView(id, (id) => {itemModel.images.remove(id);}, setPreviewHandler);
    }

    const addUploadedImage = (id) => {
        console.log("addUploadedImage", id);
        addImageView(id, (id) => {itemModel.uploadedImages.remove(id);}, setPreviewHandler);
    }

    const removeImage = (id) => {
        console.log("removeImage", id);
        service.deleteImage(id).then((response) => {
            if (response.ok()) {
                const wrapper = uploadedImageViews[id];
                if (wrapper) {
                    wrapper.remove();
                    delete uploadedImageViews[id];
                }
            } else {
                console.error("Error deleting image:", response.error);
            }
        });
    };

    itemModel.images.get().forEach(addExistingImage);
    subscriptions.push(itemModel.images.subscribeRemove(removeImage));
    subscriptions.push(itemModel.uploadedImages.subscribeAdd(addUploadedImage));
    subscriptions.push(itemModel.uploadedImages.subscribeRemove(removeImage));

    const destroy = () => {
        subscriptions.forEach(unsub => unsub());
    }

    return { destroy };
}