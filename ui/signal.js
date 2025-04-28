export function Signal(initialValue) {
    let value = initialValue;
    const subscribers = new Set();

    const notify = () => {
        subscribers.forEach(fn => fn(value));
    };

    const set = (newValue) => {
        if (newValue instanceof Promise) {
            newValue.then(resolved => {
                value = resolved;
                notify();
            });
        } else {
            value = newValue;
            notify();
        }
    };

    const get = () => value;

    const subscribe = (fn) => {
        subscribers.add(fn);
        fn(value);
        return () => subscribers.delete(fn);
    };

    const unsubscribeAll = () => {
        subscribers.clear();
    };

    return { get, set, subscribe, unsubscribeAll };
}

export function CollectionSignal(initialItems = []) {
    let items = Array.from(initialItems);
    const subscribersAdd = new Set();
    const subscribersRemove = new Set();


    const notifyAdd = (value) => {
        subscribersAdd.forEach(fn => fn(value));
    };

    const add = (value) => {
        items.push(value);
        notifyAdd(value);
    };

    const notifyRemove = (value) => {
        subscribersRemove.forEach(fn => fn(value));
    };

    const remove = (value) => {
        const idx = items.indexOf(value);
        if (idx !== -1) {
            items.splice(idx, 1);
            notifyRemove(value);
        }
    };

    const subscribeAdd = (fn) => {
        subscribersAdd.add(fn);
        return () => subscribersAdd.delete(fn);
    };

    const subscribeRemove = (fn) => {
        subscribersRemove.add(fn);
        return () => subscribersRemove.delete(fn);
    }

    const get = () => items.slice();

    const unsubscribeAll = () => {
        subscribersAdd.clear();
        subscribersRemove.clear();
    };

    return { get, add, remove, subscribeAdd, subscribeRemove, unsubscribeAll };
}