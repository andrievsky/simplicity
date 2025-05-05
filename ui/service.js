export function BackendService(host) {
    const listItems = async function () {
        return fetchResource(`${host}/api/item/`, 'GET');
    }

    const getItem = async function (id) {
        if (!id) return ErrorResult("ID is required");
        return fetchResource(`${host}/api/item/${id}`, 'GET');
    }

    const createItem = async function (item) {
        if (!item) return ErrorResult("Item is required");
        return fetchResource(`${host}/api/item/`, 'POST', item);
    }

    const updateItem = async function (id, item) {
        if (!item) return ErrorResult("Item is required");
        return fetchResource(`${host}/api/item/${id}`, 'PUT', item);
    }

    const uploadImage = async function (file) {
        if (!file) return ErrorResult("File is required");
        const formData = new FormData();
        formData.append('file', file);
        const headers = new Headers();
        headers.append('Accept', 'application/json');
        return fetchResource(`${host}/api/image/upload`, 'POST', formData, headers);
    }

    const deleteImage = async function (id) {
        if (!id) return ErrorResult("ID is required");
        return fetchResource(`${host}/api/image/files/${id}`, 'DELETE');
    }


    const getVersion = async function () {
        return fetchResource(`${host}/api/version`, 'GET');
    }

    return {listItems, getItem, createItem, updateItem, uploadImage, deleteImage, getVersion};
}
const TIMEOUT_MS = 15000;

const fetchResource = async function (url, method = 'GET', body = null, headers = null) {
    if (!url) throw new Error("URL is required");
    if (!method) throw new Error("Method is required");

    async function parseJsonSafe(response) {
        try {
            const text = await response.text();
            const data = text ? JSON.parse(text) : null;
            return { data, error: null };
        } catch (err) {
            console.error("Error parsing JSON", err);
            return { data: null, error: err };
        }
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), TIMEOUT_MS);

    const isFormData = body instanceof FormData;

    // Set default headers if not provided
    if (!headers) {
        headers = new Headers();
        if (!isFormData) {
            headers.append('Content-Type', 'application/json');
        }
        headers.append('Accept', 'application/json');
    }

    try {
        const response = await fetch(url, {
            method,
            headers,
            body: body
                ? isFormData
                    ? body // send FormData directly
                    : JSON.stringify(body) // JSON for other types
                : null,
            signal: controller.signal
        });

        clearTimeout(timeoutId);

        const { data, error: parseError } = await parseJsonSafe(response);

        if (!response.ok) {
            const errorMessage = data?.error || `HTTP error ${response.status}`;
            console.error("Fetch error response:", errorMessage, response);
            return new Result(null, response.status, errorMessage);
        }

        if (parseError) {
            return new Result(null, response.status, `Parsing error: ${parseError.message}`);
        }

        return new Result(data, response.status, null);

    } catch (error) {
        clearTimeout(timeoutId);
        const isAbort = error.name === 'AbortError';
        console.error("Fetch error:", error);
        return new Result(null, 0, isAbort ? 'Request timed out' : `Fetch error: ${error.message}`);
    }
}



function Result(data, status, error) {
    this.data = data;
    this.status = status;
    this.error = error;

    this.ok = function () {
        return this.error === null;
    };
}

function ErrorResult(message) {
    return new Result(null, 400, message);
}