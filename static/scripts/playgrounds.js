const uploadPhotosInputs = document.querySelectorAll("button.upload-photos");

const uploadPhotos = async (files, playgroundId) => {
    const formData = new FormData();
    for (const file of files) {
        formData.append('files[]', file);
    }
    const resp = await fetch(`/api/playground/${playgroundId}/upload`, {
        method: "POST",
        body: formData,
    });
    if (!resp.ok) {
        throw Error(`Server returned status: ${resp.status} - ${await resp.text()}`);
    }
    return await resp.json();
};

for (const upload of uploadPhotosInputs) {
    upload.onclick = async (ev) => {
        const input = document.querySelector(`input.upload-photos[data-id="${ev.target.dataset.id}"]`);
        const files = Array.from(input.files);
        const playgroundId = input.dataset.id;
        const uploadedPhotos = await uploadPhotos(files, playgroundId);
        if (uploadedPhotos && uploadedPhotos.length) {
            window.location.reload();
        }
    };
}