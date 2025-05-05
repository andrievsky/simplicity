import {ItemListCard} from "./item-list-card.js";
import {ItemEditForm} from "./item-edit-form.js";

export function ItemListComponent(container, model, service, templater) {
    const subscriptions = [];
    const showItems = (items) => {
        if (!items || items.length === 0) {
            container.innerHTML = '<p>No items available.</p>';
            return;
        }

        let nodes = [];

        items.forEach((item) => {
            const frag = templater.cloneTemplate("item-template");
            const card = new ItemListCard(item, frag, (item) => {
                const editForm = new ItemEditForm(item, model, service, templater);
                model.selectedItem.set(editForm);
            });
            nodes.push(card.frag);
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