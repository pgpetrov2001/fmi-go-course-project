const patchPhoto = async (attr, val, photoId) => {
    const resp = await fetch(`/api/photo/${photoId}`, {
        method: "PATCH",
        body: JSON.stringify({
            [attr]: val,
        }),
        headers: {
            "Content-Type": "application/json",
        },
    });
    if (!resp.ok) {
        throw Error(`Server returned status ${resp.status} - ${await resp.text()}`);
    }
};

const inputs = document.querySelectorAll('.gallery-item input[type="checkbox"]');

for (const input of inputs) {
    input.onchange = async (ev) => {
        const photo = ev.target.parentElement.parentElement.parentElement;
        const photoId = photo.dataset.id;
        const attr = ev.target.dataset.attr;
        const val = ev.target.checked;
        try {
            await patchPhoto(attr, val, photoId);
        } catch(err) {
            console.error(err.message);
            ev.preventDefault();
        }
    };
}