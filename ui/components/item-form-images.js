import {ItemFormImage} from "./item-form-image.js";

export function ItemFormImages(container, service, signal) {
    const uploadedImageViews = {};
    const subscriptions = [];
    const originalImages = signal.get()

    const setPreviewHandler = (id) => {
        console.log("setPreviewHandler", id)
    }

    const addImageView = (id, removeHandler, setPreviewHandler) => {
        const imageView = new ItemFormImage(id, removeHandler, setPreviewHandler);
        uploadedImageViews[id] = container.appendChild(imageView.wrapper);
    }

    const addImage = (id) => {
        console.log("addImage", id);
        const canBeRemoved = !originalImages.includes(id);
        const removeHandler = canBeRemoved ? (id) => {signal.remove(id);} : null;
        addImageView(id, removeHandler, setPreviewHandler);
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
    signal.get().forEach(addImage);
    subscriptions.push(signal.subscribeAdd(addImage));
    subscriptions.push(signal.subscribeRemove(removeImage));

    const destroy = () => {
        subscriptions.forEach(unsub => unsub());
    }

    return { destroy };
}