//Code generated with chatgpt
function validateForm(event) {
    let form = event.target.closest("form");
    if (!form.checkValidity()) {
        event.preventDefault(); // Prevent HTMX from sending request
        form.reportValidity(); // Show default validation messages
    }
}

document.getElementById("login-form").addEventListener("submit", function(event) {
    const submitButton = document.getElementById("submit-button");
    submitButton.disabled = true;
});

// Re-enable submit button after HTMX response
document.body.addEventListener("htmx:afterSwap", function(event) {
    const submitButton = document.getElementById("submit-button");
    submitButton.disabled = false;
});
