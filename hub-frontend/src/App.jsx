import { useState, useEffect } from 'react'
import { Slider } from '@mui/material'
import './App.css'

function DeviceCard({ device }) {
  const [brightness, setBrightness] = useState(device.state.brightness)
  const [isOn, setIsOn] = useState(device.state.power === true || device.state.power === 'on')

  const toggle = () => {
    fetch(`/api/devices/${device.id}/command`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: isOn ? 'turn_off' : 'turn_on' }),
    })
      .then(res => res.json())
      .then(() => setIsOn(prev => !prev))
      .catch(err => console.error("ERROR:", err))
  }

  const handleBrightness = (e, value) => {
    setBrightness(value)
  }

  const handleBrightnessCommitted = (e, value) => {
    setBrightness(value)
    fetch(`/api/devices/${device.id}/command`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'set_brightness', params: { value } }),
    })
      .then(res => res.json())
      .catch(err => console.error('Error:', err))
  }
   return (
    <li>
      {device.state.name}
      <button
        onClick={toggle}
        style={{
          marginLeft: '10px',
          backgroundColor: isOn ? 'green' : 'grey',
          color: 'white',
          border: 'none',
          borderRadius: '20px',
          padding: '5px 15px',
          cursor: 'pointer',
        }}
      >
        {isOn ? 'On' : 'Off'}
      </button>
      <Slider
        value={brightness}
        onChange={handleBrightness}
        onChangeCommitted={handleBrightnessCommitted}
        min={0}
        max={100}
      />
    </li>
  )
}


function App() {
  const [devices, setDevices] = useState([])

  useEffect(() => {
    fetch('/api/devices')
      .then(res => res.json())
      .then(data => {
        const sorted = data.devices.sort((a, b) => a.state.name.localeCompare(b.state.name))
        setDevices(sorted)
      })
      .catch(err => console.error('Error:', err))
  }, [])

  return (
    <div>
      <ul>
        {devices.map(device => (
          <DeviceCard key={device.id} device={device} />
        ))}
      </ul>
    </div>
  )
}


export default App
