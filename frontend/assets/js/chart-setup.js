document.addEventListener("DOMContentLoaded", function () {
    // Hosts + ports
    document.getElementById("hostsUp").textContent = hostsUp;
    document.getElementById("openPorts").textContent = openPorts;

    // Pie chart for port states
    new Chart(document.getElementById("portStatusPie"), {
        type: 'pie',
        data: {
            labels: Object.keys(portStates),
            datasets: [{
                data: Object.values(portStates),
                backgroundColor: ['#28a745', '#dc3545', '#ffc107'],
            }]
        }
    });

    // Bar chart for port frequencies
    new Chart(document.getElementById("popularPortsBar"), {
        type: 'bar',
        data: {
            labels: Object.keys(portFreq),
            datasets: [{
                label: 'Count',
                data: Object.values(portFreq),
                backgroundColor: '#007bff'
            }]
        },
        options: {
            scales: {
                x: { title: { display: true, text: 'Port Number' } },
                y: { title: { display: true, text: 'Frequency' } }
            }
        }
    });

    // Horizontal bar for services
    new Chart(document.getElementById("servicesBar"), {
        type: 'bar',
        data: {
            labels: Object.keys(services),
            datasets: [{
                label: 'Service Count',
                data: Object.values(services),
                backgroundColor: '#17a2b8'
            }]
        },
        options: {
            indexAxis: 'y',
            scales: {
                x: { title: { display: true, text: 'Count' } },
                y: { title: { display: true, text: 'Service Name' } }
            }
        }
    });
});
