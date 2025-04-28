import {cloneTemplate, updateTemplate} from "../template.js";

export function ItemListCard(item, template, editHandler) {
    const frag = cloneTemplate(template);
    updateTemplate(frag, item);
    const editButton = frag.querySelector('.edit-button');
    editButton.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        editHandler(item);
    });
    const preview = frag.querySelector('.preview');
    if (item.images && item.images.length > 0) {
        preview.src = "/api/image/files/" + item.images[0] + "?format=web-preview-280";
    }
    return {frag};
}