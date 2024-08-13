import { useEffect } from 'react';
import './App.css';
import utils from './services/utils';

function App() {
  useEffect(() => {
    console.log(utils.getAccessToken());
  })

  return (
    <>
      <h1>VGLibrary</h1>
    </>
  );
}

export default App;
