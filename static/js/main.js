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

//Code generated via chatgpt
document
    .getElementById("login-form")
    .addEventListener("submit", function(event) {
        const submitButton = document.getElementById("submit-button");
        submitButton.disabled = true;
    });

// Re-enable submit button after HTMX response
document.body.addEventListener("htmx:afterSwap", function(event) {
    const submitButton = document.getElementById("submit-button");
    submitButton.disabled = false;
});

const max = document.getElementById("rCapital").max;
let lastCaptalString = 10000;
let lastPeriodString = 12;

function isValid(input) {
    if (max < 120) {
        input.checked = false;
        document.getElementById("fidu").checked = true;
        alert(
            "Your apportation account doesn't have enough funds to use this option",
        );
    }
}

function validateForm(event) {
    let form = event.target.closest("form");
    if (!form.checkValidity()) {
        event.preventDefault();
        form.reportValidity();
    }
}

function verifyPeriods(input) {
    if (input.value > 12) {
        input.value = 12;
    } else if (input.value < 1) {
        input.value = 1;
    }

    if (input.value.length != 0) {
        lastPeriodString = input.value;
        return;
    }
    input.value = lastPeriodString;
}

function verifyAmount(input) {
    let lmax = document.getElementById("fidu").checked ? 10000 : max;

    if (input.value > lmax) {
        input.value = lmax;
    } else if (input.value < 1) {
        input.value = 1;
    }

    if (input.value.length != 0) {
        lastCaptalString = input.value;
        return;
    }
    input.value = lastCaptalString;
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

let lastInputString = 0.01;
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
    console.log(event)
}

