export function FooterComponent(container, service) {
    const init = () => {
        container.innerHTML = "© 2025 Simplicity UI - All rights reserved.";
        service.getVersion().then(response => {
            if (!response.ok) {
                console.error("Error fetching version", response.error);
                return;
            }
            container.innerHTML = "© 2025 Simplicity UI - All rights reserved. Backend Version: " + response.data.version;
        })

    }



    return {init};
}