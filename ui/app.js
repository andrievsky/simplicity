import {BackendService} from './backend.js';
import {Model} from "./model.js";
import {HeaderComponent} from "./components/header.js";
import {ItemListComponent} from "./components/item-list.js";
import {FooterComponent} from "./components/footer.js";
import {ModalComponent} from "./components/modal.js";

const init = async () => {
    const header = document.getElementById('headerContainer');
    const content = document.getElementById('contentContainer');
    const footer = document.getElementById('footerContainer');
    const modal = document.getElementById("modalContainer");
    const itemTemplate = document.getElementById('item-template');
    const itemEditTemplate = document.getElementById('item-edit-template');

    let backend = new BackendService("");
    let model = new Model();
    let modalComponent = new ModalComponent(modal, model, itemEditTemplate);
    modalComponent.init();
    let headerComponent = new HeaderComponent(header, model);
    headerComponent.init();
    let itemListComponent = new ItemListComponent(content, model, itemTemplate);
    itemListComponent.init();
    let footerComponent = new FooterComponent(footer);
    footerComponent.init();
    await backend.listItems().then((result) => {
        if (result.ok()) {
            model.items.set(result.data);
        } else {
            console.error("Error loading items:", result.error);
        }
    });
};

window.addEventListener('load', init);