import React, { useEffect, useState } from 'react';

function ListRecords() {
  const [records, setRecords] = useState([]);

  useEffect(() => {
    fetch('http://localhost/running-statistics')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        console.log('Fetched data:', data); // Debugging log
        setRecords(data.statistics); // Access the statistics array
      })
      .catch(error => console.error('Error fetching records:', error));
  }, []);

  return (
    <div className="container">
      <h1>Running Statistics</h1>
      <ul>
        {records.map(record => (
          <li key={record.id}>
            <strong>Date:</strong> {record.date}, <strong>Distance:</strong> {record.distance}, <strong>Time:</strong> {record.time}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ListRecords;