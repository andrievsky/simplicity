import {FetchResource, Result, FromTemplate} from './utils.js';

const init = async () => {
    const header = null || document.getElementById('header_root');
    const content = null || document.getElementById('page_root');
    const footer = null || document.getElementById('footer_root');

    // Render the header and footer of the page.
    //header.innerHTML = await Navbar.render();
    //await Navbar.after_render();
    //footer.innerHTML = await Footer.render();
    //await Footer.after_render();

    let itemService = new ItemListService();
    let itemModel = new ItemListModel(itemService);
    let itemView = new ItemListView(content);
    let itemController = new ItemListController(itemView, itemModel);
    await itemController.init()
};


// Listen on page load.
window.addEventListener('load', init);

function ItemListService() {
    this.listItems = async function () {
        return FetchResource('/api/item', 'GET')
    };
}

function ItemListModel(service) {
    this.getAllItems = async function () {
        return service.listItems();
    }
}

function ItemListController(view, model) {
    this.init = async () => {
        let result = await model.getAllItems();
        if (result.isEmpty()) {
            view.showEmpty();
            return;
        }
        if (!result.isSuccess()) {
            view.showError(result.error);
            return;
        }
        view.addItems(result.data);
    }
}

function ItemListView(container) {
    this.addItems = function (items) {
        if (!items || items.length === 0) {
            this.showEmpty();
            return;
        }
        container.innerHTML = "";
        items.forEach((item) => {
            container.appendChild(FromTemplate('item-template', item));
        });
    };

    this.showEmpty = function () {
        container.innerHTML = "<p>No items found</p>";
    }

    this.showError = function (err) {
        container.innerHTML = `<p class="error">Error: ${err.message}</p>`;
    };
}
