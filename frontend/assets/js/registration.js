$(document).ready(function () {
    // Registration form validation
    $("form[action='/register']").on("submit", function (e) {
        let isValid = true;
        const username = $("#username").val().trim();
        const email = $("#email").val().trim();
        const password = $("#password").val().trim();
        const confirmPassword = $("#confirm-password").val().trim();

        // Clear previous error messages
        $(".form-error").remove();

        // Validate username
        if (!username) {
            isValid = false;
            $("#username").after('<div class="form-error text-danger mt-1">Username is required.</div>');
        }

        // Validate email
        if (!email) {
            isValid = false;
            $("#email").after('<div class="form-error text-danger mt-1">Email is required.</div>');
        }

        // Validate password
        if (!password) {
            isValid = false;
            $("#password").after('<div class="form-error text-danger mt-1">Password is required.</div>');
        }

        // Validate confirm password
        if (!confirmPassword) {
            isValid = false;
            $("#confirm-password").after('<div class="form-error text-danger mt-1">Please confirm your password.</div>');
        } else if (password !== confirmPassword) {
            isValid = false;
            $("#confirm-password").after('<div class="form-error text-danger mt-1">Passwords must match.</div>');
        }

        // Prevent form submission if validation fails
        if (!isValid) {
            e.preventDefault();
        }
    });
});