import {cloneTemplate, updateTemplate} from "../template.js";

export function ModalComponent(container, model, itemEditTemplate) {
    const subscriptions = [];
    const showEditItem = function (item) {
        const frag = cloneTemplate(itemEditTemplate);
        updateTemplate(frag, item);
        const saveButton = frag.querySelector('.save-button');
        saveButton.addEventListener("click", (e) => {
            e.preventDefault();
            e.stopPropagation();
            hide();
        });
        const cancelButton = frag.querySelector('.cancel-button');
        cancelButton.addEventListener("click", (e) => {
            e.preventDefault();
            e.stopPropagation();
            hide();
        });
        container.replaceChildren(frag);
        show();
    };
    const show = function () {
        container.style.display = "flex";
    };

    const hide = function () {
        container.style.display = "none";
    };

    const init = () => {
        subscriptions.push(model.selectedItem.subscribe((item) => {
            if (item) {
                showEditItem(item);
            } else {
                hide();
            }
        }));
    };

    const destroy = () => {
        subscriptions.forEach(unsub => unsub());
        if (container) {
            container.replaceChildren();
        }
    }

    return {init, destroy}
}