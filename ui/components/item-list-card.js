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

    return {frag};
}