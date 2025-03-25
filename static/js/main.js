
let lastInputString = 0.01;
// This code was generated with chatGPT
document.addEventListener("htmx:afterRequest", function(evt) {
    if (evt.detail.xhr.getResponseHeader("HX-Location")) {
        setTimeout(() => {
            window.location.reload(); // Forces full reload after redirect
        }, 100); // Small delay to ensure redirection is processed
    }
});

// Snippet got from the W3C $$
var dropdown = document.getElementsByClassName("dropdown-btn");
var i;

for (i = 0; i < dropdown.length; i++) {
    dropdown[i].addEventListener("click", function() {
        this.classList.toggle("active");
        var dropdownContent = this.nextElementSibling;
        if (dropdownContent.style.display === "flex") {
            dropdownContent.style.display = "none";
        } else {
            dropdownContent.style.display = "flex";
        }
    });
}
//$$

//Code generated with chatgpt
function validateForm(event) {
    let form = event.target.closest("form");
    if (!form.checkValidity()) {
        event.preventDefault(); // Prevent HTMX from sending request
        form.reportValidity(); // Show default validation messages
    }
}
function changeOption(input) {
    for (i = 0; i < input.length; i++) {
        if (input.options[i].selected) {
            opt = document.getElementById("remaining").options[i];
            opt.selected = true;
            document.getElementById("MaxAmount").max = opt.value;
        }
    }
}

function verifyPaymentAmount(input) {
    for (i = 0; i < input.value.length; i++) {
        if (input.value.at(i) === "." && input.value.substring(i).length > 2) {
            input.value = parseFloat(input.value).toFixed(2);
            break;
        }
    }

    if (input.value.length != 0) {
        lastInputString = input.value;
    }
    input.value = lastInputString;

    if (parseFloat(input.value) > input.max) {
        input.value = input.max;
    } else if (input.value < 0) {
        input.value = 0;
    }
}

function handleResponse(event) {
    let xhr = event.detail.xhr;
    let status = xhr.getResponseHeader("HX-Status");
    let message = xhr.getResponseHeader("HX-Message");
    if (status === "202") {
        return ;
    } else if (status === "200") {
        alert(message); // Show the success message
        document.getElementById("form-div").reset();
        setTimeout(() => {
            window.location.reload(); // Forces full reload after redirect
        }, 100); // Small delay to ensure redirection is processed
    } else if (status === "400") {
        alert(message); // Show the error message
    } else {
        alert("Something went wrong. Please try again.");
    }
}

