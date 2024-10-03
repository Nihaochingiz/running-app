import React, { useState } from 'react';

function CreateRecord() {
  const [date, setDate] = useState('');
  const [distance, setDistance] = useState('');
  const [time, setTime] = useState('');

  const handleSubmit = (event) => {
    event.preventDefault();
    const newRecord = { date, distance, time };

    fetch('http://localhost/running-statistics', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newRecord),
    })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
        setDate('');
        setDistance('');
        setTime('');
      })
      .catch(error => console.error('Error:', error));
  };

  return (
    <div className="container">
      <h1>Create a New Running Statistic</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Date:</label>
          <input type="text" value={date} onChange={(e) => setDate(e.target.value)} />
        </div>
        <div className="form-group">
          <label>Distance:</label>
          <input type="text" value={distance} onChange={(e) => setDistance(e.target.value)} />
        </div>
        <div className="form-group">
          <label>Time:</label>
          <input type="text" value={time} onChange={(e) => setTime(e.target.value)} />
        </div>
        <button type="submit">Create</button>
      </form>
    </div>
  );
}

export default CreateRecord;