export function BackendService(host) {
    const listItems = async function () {
        return fetchResource(`${host}/api/item`, 'GET');
    }

    return {listItems};
}
const TIMEOUT_MS = 5000;

const fetchResource = async function (url, method = 'GET', body = null) {
    if (!url) throw new Error("URL is required");
    if (!method) throw new Error("Method is required");

    async function parseJsonSafe(response) {
        try {
            const data = await response.json();
            return { data, error: null };
        } catch (err) {
            console.error("Error parsing JSON", err);
            return { data: null, error: err };
        }
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), TIMEOUT_MS);

    try {
        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: body ? JSON.stringify(body) : null,
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