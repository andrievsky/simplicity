import {cloneTemplate, updateTemplate} from "../template.js";
import {FormReader} from "../form-reader.js";

export function ItemEditForm(item, model, service, templates) {
    const frag = cloneTemplate(templates["item-edit"]);
    updateTemplate(frag, item);
    const close = () => {
        model.selectedItem.set(null);
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
    let uploadedImages = [];

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
            const placeholder = createPreview("loading");
            previewList.appendChild(placeholder);

            const response = await service.uploadImage(file);
            if (response.ok()) {
                const url = "/api/image/files/"+response.data.id+"?format=canonical";
                placeholder.replaceWith(createPreview("success", url));
                uploadedImages.push(url);
            } else {
                placeholder.replaceWith(createPreview("error"));
            }
        }
    }

    function createPreview(state, url = null) {
        const wrapper = document.createElement("div");
        wrapper.className = "preview";

        if (state === "loading") {
            wrapper.textContent = "Uploading...";
        } else if (state === "success") {
            console.log("Image uploaded successfully:", url)
            const img = document.createElement("img");
            img.src = url;
            img.className = "preview-image";

            const remove = document.createElement("button");
            remove.textContent = "✖";
            remove.className = "remove-image";
            remove.addEventListener("click", () => {
                uploadedImages = uploadedImages.filter(u => u !== url);
                wrapper.remove();
            });

            wrapper.appendChild(img);
            wrapper.appendChild(remove);
        } else if (state === "error") {
            wrapper.textContent = "Failed to upload";
            wrapper.classList.add("error");
        }

        return wrapper;
    }


    return frag;
}
