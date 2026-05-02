import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [devices, setDevices] = useState([])

  useEffect(() => { 
    fetch('/api/devices')
      .then(res => {
        console.log('Status:', res.status)
        return res.json()
      })
      .then(data => {
        console.log('Data:', data)
        setDevices(data.devices)
      })
      .catch(err => console.error('Error:', err))
  }, [])

  
const toggleDevice = (id, currentPower) => {
  const isOn = currentPower === true || currentPower === 'on'
  const action = isOn ? 'turn_off' : 'turn_on'

  fetch(`/api/devices/${id}/command`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ action }),
  })
    .then(res => res.json())
    .then(() => {
      setDevices(prev =>
        prev.map(device =>
          device.id === id
            ? { ...device, state: { ...device.state, power: !isOn } }
            : device
        )
      )
    })
    .catch(err => console.error('Error:', err))
}


  return (
    <div>
      {devices.length === 0 ? (
        <p>Loading devices.</p>
      ) : (
      <ul>
        {devices.map(device => (
          <li key={device.id}>
            {device.state.name}
            <button
              onClick={() => toggleDevice(device.id, device.state.power)}
              style={{
                backgroundColor: device.state.power === true || device.state.power === 'on' ? 'green' : 'grey',
                color: 'white',
                border: 'none',
                borderRadius: '20px',
                padding: '5px 15px',
                cursor: 'pointer',
              }}
            >
              {device.state.power === true || device.state.power === 'on' ? 'On' : 'Off'}
            </button>
          </li>
        ))}
      </ul>
      )}
    </div>
  )

}

export default App
