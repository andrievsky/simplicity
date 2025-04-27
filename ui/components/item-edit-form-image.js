export function ItemEditFormImage(id, removeHandler) {
    const wrapper = document.createElement("div");
    wrapper.className = "preview";
    const url = "/api/image/files/"+id+"?format=web-thumb-sq";
    console.log("ItemEditFormImage:", url)
    const img = document.createElement("img");
    img.src = url;
    img.className = "preview-image";
    wrapper.appendChild(img);

    if (removeHandler) {
        const remove = document.createElement("button");
        remove.textContent = "✖";
        remove.className = "remove-image";

        remove.addEventListener("click", () => {
            removeHandler(id);
        });
        wrapper.appendChild(remove);
    }

    return {wrapper};
}