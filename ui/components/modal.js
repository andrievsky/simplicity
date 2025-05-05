export function ModalComponent(container, model, service, templater) {
    const subscriptions = [];
    const show = function (frag) {
        container.replaceChildren(frag);
        container.style.display = "flex";
    };

    const hide = function () {
        container.style.display = "none";
        container.replaceChildren();
    };

    const init = () => {
        subscriptions.push(model.selectedItem.subscribe((view) => {
            if (view) {
                show(view);
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