export function ItemFormInput(input, signal) {
    input.value = signal.get();
    input.addEventListener('input', (e) => {
        signal.set(e.target.value);
    });

    const destroy = () => {}

    return {destroy};

}