let smartDeck = null

const eventListener = {
  onKeyUp: (coordinates) => _sendCommand(JSON.stringify(coordinates)),
  onDeviceSleep: () => _sendCommand('DEVICE-SLEEPING'),
  onDeviceWake: () => _sendCommand('DEVICE-WAKE')
}

connectSocket = (inPort, inPluginUUID, inRegisterEvent, _inInfo) => {
  const websocket = new WebSocket(`ws://localhost:${inPort}`)
  smartDeck = new WebSocket('ws://localhost:1337')
  
  websocket.onopen = () => {
    const manifest = {
      event: inRegisterEvent,
      uuid: inPluginUUID
    }
    
    websocket.send(JSON.stringify(manifest))
  }
  
  websocket.onmessage = (evt) => {
    const data = JSON.parse(evt.data)
    const event = data.event
    const payload = data.payload || {}
    
    if (event === 'keyUp') {
      eventListener.onKeyUp(payload.coordinates)
    } else if (event === 'deviceDidDisconnect') {
      eventListener.onDeviceSleep()
    } else if (event === 'systemDidWakeUp') {
      eventListener.onDeviceWake()
    }
  }
  
  websocket.onclose = () => {
  }
}

_sendCommand = (dataToSend) => {
  if (smartDeck.readyState === smartDeck.OPEN) {
    smartDeck.send(dataToSend)
  } else {
    smartDeck = new WebSocket('ws://localhost:1337')
    smartDeck.onopen = () => {
      smartDeck.send(dataToSend)
    }
  }
}
