import {ItemEditForm} from "./item-edit-form.js";

export function ModalComponent(container, model, service, templates) {
    const subscriptions = [];
    const showEditItem = function (item) {
        if (!item) {
            console.error("No item to edit");
            return;
        }
        const close = function () {
            model.selectedItem.set(null);
        }
        const itemEditForm = new ItemEditForm(item, model, service, templates, close);
        show(itemEditForm);
    };
    const show = function (frag) {
        container.replaceChildren(frag);
        container.style.display = "flex";
    };

    const hide = function () {
        container.style.display = "none";
        container.replaceChildren();
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