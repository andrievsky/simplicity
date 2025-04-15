export async function FetchResource(url, method, body) {
    if (url === null || url === undefined) {
        throw new Error("URL is required");
    }
    if (method === null || method === undefined) {
        throw new Error("Method is required");
    }
    return fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: body
    }).then((response) => {
        if (!response.ok) {
            if (response.status === 404) {
                console.error("Resource not found", response);
                return new Result(null, response.status, "Resource not found");
            }
            return response.json().then((data) => {
                console.error("Error response", data.error, response);
                return new Result(null,response.status, data.error);
            }, (error) => {
                console.error("Error parsing error response", error, response);
                return new Result(null, response.status, `Parsing error: ${response}, error: ${error}`);
            });
        }
        return response.json().then((data) => {
            return new Result(data, response.status, null);
        }, (error) => {
            console.error("Error parsing response", error, response);
            return new Result(null, response.status, `Parsing error: ${response}, error: ${error}`);
        });
    }, (error) => {
        console.error("Fetch error", error);
        return new Result(null, new Error(`Fetch error: ${error}`));
    });
}

export function Result(data, status, error) {
    this.data = data;
    this.status = status;
    this.error = error;

    this.isSuccess = function () {
        return this.error === null;
    };

    this.isEmpty = function () {
        return this.status === 404;
    };
}

export function FromTemplate(templateId, data) {
    if (!templateId) {
        console.error('Template ID is required.');
        return null;
    }
    if (!data) {
        console.error('Data is required.');
        return null;
    }
    const template = document.getElementById(templateId);
    if (!template) {
        console.error(`Template with ID '${templateId}' not found.`);
        return null;
    }

    const clone = template.content.cloneNode(true);
    const bindableElements = clone.querySelectorAll('[data-key]');

    bindableElements.forEach((element) => {
        const key = element.getAttribute('data-key');
        if (key in data) {
            element.textContent = data[key];
        } else {
            console.warn(`No data found for key '${key}'.`);
        }
    });
    return clone;
}

export const Signal = function () {
    this.handlers = [];
    this.add = function (handler) {
        this.handlers.push(handler);
    };
    this.remove = function (handler) {
        this.handlers = this.handlers.filter((h) => h !== handler);
    };
    this.fire = function (data) {
        this.handlers.forEach((h) => h(data));
    };
}