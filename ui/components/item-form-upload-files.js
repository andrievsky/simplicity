export function ItemFormUploadFiles(frag, service, signal) {
    const dropZone = frag.querySelector('.drop-zone');
    const input = frag.querySelector('.image-upload');

    dropZone.addEventListener("click", () => input.click());

    dropZone.addEventListener("dragover", (e) => {
        e.preventDefault();
        dropZone.classList.add("highlight");
    });
    dropZone.addEventListener("dragleave", () => {
        dropZone.classList.remove("highlight");
    });
    dropZone.addEventListener("drop", async (e) => {
        e.preventDefault();
        dropZone.classList.remove("highlight");
        await handleFiles(e.dataTransfer.files);
    });

    input.addEventListener("change", async (e) => {
        await handleFiles(e.target.files);
    });

    async function handleFiles(files) {
        for (const file of files) {
            const response = await service.uploadImage(file);
            if (response.ok()) {
                signal.add(response.data.id);
            } else {
                console.error("Error uploading image:", response.error);
            }
        }
    }
}