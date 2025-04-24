import {BackendService} from './service.js';
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
    const templates = {
        "item": document.getElementById('item-template'),
        "item-edit": document.getElementById('item-edit-template')
    };

    let service = new BackendService("");
    let model = new Model();
    let modalComponent = new ModalComponent(modal, model, service, templates);
    modalComponent.init();
    let headerComponent = new HeaderComponent(header, model);
    headerComponent.init();
    let itemListComponent = new ItemListComponent(content, model, templates);
    itemListComponent.init();
    let footerComponent = new FooterComponent(footer, service);
    footerComponent.init();
    await service.listItems().then((result) => {
        if (result.ok()) {
            model.items.set(result.data);
        } else {
            console.error("Error loading items:", result.error);
        }
    });
};

window.addEventListener('load', init);