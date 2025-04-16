export function cloneTemplate(template) {
    return template.content.cloneNode(true);
}

export function updateTemplate(template, data) {
    const bindableElements = template.querySelectorAll('[data-key]');
    bindableElements.forEach((element) => {
        const key = element.getAttribute('data-key');
        if (key in data) {
            if (element.tagName === 'INPUT') {
                element.value = data[key];
            } else if (element.tagName === 'SELECT') {
                element.value = data[key];
            } else if (element.tagName === 'TEXTAREA') {
                element.value = data[key];
            } else {
                element.textContent = data[key];
            }
        } else {
            console.warn(`No data found for key '${key}'.`);
        }
    });
}