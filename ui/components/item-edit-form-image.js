export function ItemEditFormImage(id, removeHandler, setPreviewHandler) {
    const wrapper = document.createElement("div");
    wrapper.className = "preview";
    const url = "/api/image/files/" + id + "?format=web-thumb-sq";
    console.log("ItemEditFormImage:", url);
    const img = document.createElement("img");
    img.src = url;
    img.className = "preview-image";
    wrapper.appendChild(img);

    const buttonGroup = document.createElement("div");
    buttonGroup.className = "button-group";

    if (removeHandler) {
        const remove = document.createElement("button");
        remove.textContent = "✖";
        remove.className = "action-button";

        remove.addEventListener("click", () => {
            removeHandler(id);
        });
        buttonGroup.appendChild(remove);
    }

    if (setPreviewHandler) {
        const setPreview = document.createElement("button");
        setPreview.textContent = "P";
        setPreview.className = "action-button";

        setPreview.addEventListener("click", () => {
            setPreviewHandler(id);
        });
        buttonGroup.appendChild(setPreview);
    }

    wrapper.appendChild(buttonGroup);

    return { wrapper };
}