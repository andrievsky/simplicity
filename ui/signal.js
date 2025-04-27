export function Signal(initialValue) {
    let value = initialValue;
    const subscribers = new Set();

    const notify = () => {
        subscribers.forEach(fn => fn(value));
    };

    const set = (newValue) => {
        //console.log("set", newValue);
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
    const deltaSubs = new Set();

    const notifyDelta = (type, item) => {
        deltaSubs.forEach(fn => fn({ type, item }));
    };

    const add = (item) => {
        items.push(item);
        notifyDelta('add', item);
    };

    const remove = (item) => {
        const idx = items.indexOf(item);
        if (idx !== -1) {
            items.splice(idx, 1);
            notifyDelta('remove', item);
        }
    };

    const subscribeDelta = (fn) => {
        deltaSubs.add(fn);
        return () => deltaSubs.delete(fn);
    };

    const get = () => items.slice();

    const unsubscribeAll = () => {
        deltaSubs.clear();
    };

    return { get, add, remove, subscribeDelta, unsubscribeAll };
}