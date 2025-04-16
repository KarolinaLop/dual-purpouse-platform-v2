

$(document).ready(function () {
    // register click event for the logout button
    $('#logout').on('click', function (e) {
        e.preventDefault(); // Prevent the link from navigating

        $.ajax({
            url: '/logout',
            type: 'DELETE',
            success: function () {
                // Redirect to login or homepage after logout
                window.location.href = '/login';
            },
            error: function (xhr, status, error) {
                console.error('Logout failed:', error);
            }
        });
    });
});