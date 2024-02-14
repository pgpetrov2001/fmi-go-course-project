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

const postReview = async (rating, content, playgroundId) => {
    const resp = await fetch(`/api/playground/${playgroundId}/review`, {
        method: "POST",
        body: JSON.stringify({
            stars: rating,
            content: content,
        }),
        headers: {
            "Content-Type": "application/json",
        },
    });
    if (!resp.ok) {
        throw Error(`Server returned status: ${resp.status} - ${await resp.text()}`);
    }
};

const postVote = async (up, reviewId) => {
    const resp = await fetch(`/api/review/${reviewId}/vote`, {
        method: "POST",
        body: JSON.stringify({ up }),
        headers: {
            "Content-Type": "application/json",
        },
    });
    if (!resp.ok) {
        throw Error(`Server returned status: ${resp.status} - ${await resp.text()}`);
    }
};

const uploadPhotosInputs = document.querySelectorAll("button.upload-photos");

for (const upload of uploadPhotosInputs) {
    upload.onclick = async (ev) => {
        const input = document.querySelector(`input.upload-photos[data-id="${ev.target.dataset.id}"]`);
        const files = Array.from(input.files);
        const playgroundId = input.dataset.id;
        const uploadedPhotos = await uploadPhotos(files, playgroundId);
        if (uploadedPhotos && uploadedPhotos.length) {
            // window.location.reload();
            alert('Your photos have been submitted and are awaiting approval');
        }
    };
}

const ratingStarButtons = document.querySelectorAll('.playground .clickable-star');

for (const starButton of ratingStarButtons) {
    starButton.onclick = async (ev) => {
        const ratingInput = ev.target.parentElement;
        ratingInput.classList.remove('unselected');
        let stars = 0;
        let passed = false;
        for (const star of ratingInput.children) {
            if (star.dataset.index == ev.target.dataset.index) {
                passed = true;
            }
            if (passed) {
                star.classList.add('active');
            }
            if (!passed) {
                stars++;
            }
        }
        stars = 5 - stars;

        const textarea = document.createElement("textarea");
        textarea.classList.add('review-input');
        textarea.rows = 3;
        textarea.placeholder = "Elaborate on your rating if you want";
        const button = document.createElement('button');
        button.classList.add('submit-review');
        button.innerText = "Submit Review";
        const container = ratingInput.parentElement;
        container.insertBefore(textarea, ratingInput);
        container.insertBefore(button, textarea);
        button.onclick = async () => {
            const review = await postReview(stars, textarea.value, ratingInput.dataset.id);
            container.removeChild(button);
            container.removeChild(textarea);
            alert('Thanks for submitting a review for this playground.');
            window.location.reload();
        };
    };
}

const voteButtons = document.querySelectorAll('.review-voting .vote-button');

for (const voteButton of voteButtons) {
    voteButton.onclick = async (ev) => {
        const votingEl = voteButton.parentElement;
        const reviewId = votingEl.dataset.id;
        const up = voteButton.classList.contains('like');
        if (voteButton.classList.contains('selected')) {
            return;
        }
        try {
            await postVote(up, reviewId);
            for (const currVoteButton of votingEl.children) {
                currVoteButton.classList.remove('selected');
            }
            voteButton.classList.add('selected');
        } catch (e) {
            console.error(e.message);
        }
    };
}