export function FormReader(element) {
    const string = (key, fallback = '') => {
        const el = element.querySelector(`[data-key="${key}"]`);
        if (!el) return fallback;

        if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
            return el.value.trim() || fallback;
        }

        return el.textContent.trim() || fallback;
    }

    const array = (key, fallback = []) => {
        const el = element.querySelector(`[data-key="${key}"]`);
        if (!el) return fallback;

        if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
            return el.value.split(',').map(tag => tag.trim()).filter(tag => tag !== '') || fallback;
        }

        return Array.from(el.querySelectorAll('li')).map(li => li.textContent.trim()) || fallback;
    }

    return {string, array};
}