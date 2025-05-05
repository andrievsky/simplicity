import {ItemNewForm} from "./item-new-form.js";

export function HeaderComponent(container, model, service, templater) {
    const init = () => {
        const headerContainer = document.getElementById("headerContainer");

        // Wire the 'New Item' button
        const newItemButton = document.getElementById("newItemButton");
        newItemButton.addEventListener("click", () => {
            const editForm = new ItemNewForm(model, service, templater);
            model.openModal(editForm);
        });

        // Center the search bar (already handled in styles)
        const searchBar = headerContainer.querySelector(".search-bar");
        if (!searchBar) {
            console.error("Search bar not found in header");
        }
    };

    return { init };
}