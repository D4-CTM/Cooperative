//Code generated with chatgpt
function validateForm(event) {
    let form = event.target.closest('form');
    if (!form.checkValidity()) {
        event.preventDefault(); // Prevent HTMX from sending request
        form.reportValidity();  // Show default validation messages
    }
}
