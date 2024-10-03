import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import ListRecords from './components/ListRecords';
import CreateRecord from './components/CreateRecord';
import './App.css';

function App() {
  return (
    <Router>
      <nav>
        <ul>
          <li><Link to="/">Home</Link></li>
          <li><Link to="/create">Create Record</Link></li>
        </ul>
      </nav>
      <div className="container">
        <Routes>
          <Route path="/" element={<ListRecords />} />
          <Route path="/create" element={<CreateRecord />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;