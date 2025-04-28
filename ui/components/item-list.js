import {ItemListCard} from "./item-list-card.js";

export function ItemListComponent(container, model, templates) {
    const subscriptions = [];
    const showItems = (items) => {
        if (!items || items.length === 0) {
            container.innerHTML = '<p>No items available.</p>';
            return;
        }

        let nodes = [];
        const template = templates["item"];
        items.forEach((item) => {
            const card = new ItemListCard(item, template, (item) => {
                model.selectedItem.set(item);
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