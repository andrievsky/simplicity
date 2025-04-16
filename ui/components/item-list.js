import {cloneTemplate, updateTemplate} from "../template.js";

export function ItemListComponent(container, model, itemTemplate) {
    const subscriptions = [];
    const showItems = (items) => {
        if (!items || items.length === 0) {
            container.innerHTML = '<p>No items available.</p>';
            return;
        }

        let nodes = [];
        items.forEach((item) => {
            const frag = cloneTemplate(itemTemplate);
            updateTemplate(frag, item);
            const editButton = frag.querySelector('.edit-button');
            editButton.addEventListener("click", (e) => {
                e.preventDefault();
                e.stopPropagation();
                model.selectedItem.set(item);
            });
            nodes.push(frag);
        });
        container.replaceChildren(...nodes);
    }
    const init = () => {
        subscriptions.push(model.items.subscribe(showItems));
    }
    const destroy = () => {
        subscriptions.forEach(unsub => unsub());
        if (container) {
            container.replaceChildren();
        }
    }

    return {init, destroy};
}