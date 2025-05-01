    document.addEventListener("DOMContentLoaded", function () {
        const timestamps = document.querySelectorAll(".timestamp");
        timestamps.forEach(function (timestamp) {
            const original = timestamp.textContent.trim();
            const date = new Date(original);
            if (!isNaN(date)) {
                const options = { year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric', hour12: true };
                timestamp.textContent = date.toLocaleString('en-US', options);
            }
        });
    });
