import {ItemEditFormImage} from "./item-edit-form-image.js";

export function ItemEditFormUploadFiles(frag, service, itemModel) {
    const dropZone = frag.querySelector('.drop-zone');
    const input = frag.querySelector('.image-upload');
    const previewList = frag.querySelector('.image-preview-list');

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

    const uploadedImageViews = {};
    itemModel.uploadedImages.subscribeDelta((delta => {
        if (delta.type === "add") {
            const imageView = new ItemEditFormImage(delta.item, (id) => {
                itemModel.uploadedImages.remove(id);
            });
            uploadedImageViews[delta.item] = previewList.appendChild(imageView.wrapper);
        }
        if (delta.type === "remove") {
            service.deleteImage(delta.item).then((response) => {
                if (response.ok()) {
                    removeImage(delta.item);
                } else {
                    console.error("Error deleting image:", response.error);
                }
            });
        }
    }));

    const removeImage = (id) => {
        const wrapper = uploadedImageViews[id];
        if (wrapper) {
            wrapper.remove();
            delete uploadedImageViews[id];
        }
    }

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
}