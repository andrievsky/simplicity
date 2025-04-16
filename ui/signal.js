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

    return { get, set, subscribe };
}