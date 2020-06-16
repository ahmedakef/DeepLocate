import React from 'react';
import './App.css';
import './components/mainForm/mainForm';
import MainForm from './components/mainForm/mainForm';
import Axios from 'axios';

Axios.defaults.baseURL = 'http://127.0.0.1:8080/'

function App() {
  return (
    <div className="App">
      <MainForm/>
    </div>
  );
}

export default App;
