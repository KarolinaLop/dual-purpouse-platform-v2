$(document).ready(function () {
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

