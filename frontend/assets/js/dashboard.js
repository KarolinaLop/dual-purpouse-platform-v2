$(document).ready(function () {
    // register click event for the new scan button
    $('#new-scan').on('click', function (e) {
        // disable the button to prevent multiple clicks
        $(this).prop('disabled', true);

        $('span#hint').removeClass('d-none');
        // show user feedback
        $(this).html('Starting Scan...');
        $.ajax({
            url: '/scans',
            type: 'POST',
            success: function () {
                // Redirect to scans list
                window.location.href = '/scans';
            },
            error: function (xhr, status, error) {
                console.error('Could not start scan:', error);
            }
        });
    });


    const ctx = document.getElementById('myChart');

    const data = [
        document.getElementById('open-ports').value,
        document.getElementById('closed-ports').value,
        document.getElementById('filtered-ports').value,
    ];
    new Chart(ctx, {
        type: 'pie',
        data: {
          labels: ['Open', 'Closed', 'Filtered'],
          datasets: [{
            borderColor: '#e6e9ef', // Change border color here,
            label: 'Ports',
            data: data,
            backgroundColor: ['#40a02b', '#1e66f5', '#209fb5'],
          }]
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: 'top',
            }
          }
        }
      });
});

