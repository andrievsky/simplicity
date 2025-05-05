import {BackendService} from './service.js';
import {Model} from "./model.js";
import {HeaderComponent} from "./components/header.js";
import {ItemListComponent} from "./components/item-list.js";
import {FooterComponent} from "./components/footer.js";
import {ModalComponent} from "./components/modal.js";
import {Templater} from "./template.js";

const init = async () => {
    const header = document.getElementById('headerContainer');
    const content = document.getElementById('contentContainer');
    const footer = document.getElementById('footerContainer');
    const modal = document.getElementById("modalContainer");

    const templater = new Templater();

    let service = new BackendService("");
    let model = new Model(service);
    let modalComponent = new ModalComponent(modal, model, service, templater);
    modalComponent.init();
    let headerComponent = new HeaderComponent(header, model, service, templater);
    headerComponent.init();
    let itemListComponent = new ItemListComponent(content, model, service, templater);
    itemListComponent.init();
    let footerComponent = new FooterComponent(footer, service);
    footerComponent.init();

    model.refreshItems();
};

window.addEventListener('load', init);