export function ItemEditFormInput(input, value, changeHandler) {
    input.value = value;
    input.addEventListener('input', (e) => {
        changeHandler(e.target.value);
    });

    const destroy = () => {}

    return {destroy};

}