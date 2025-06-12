## This project serves a purpose of changing register values of a DRS Device with ModBus connection
### There are several steps before actually taken along the way: 
  First, the user is asked to choose a USB-port to be connected with. 
  After that, it is checked whether there really is a functioning DRS device connected to that port. 
  If the connection is established, the current register values are loaded and the user gets to change them. 

### The project uses the following:
  Vite Bundler and React in Frontend 
  Go (Gin) in Backend
  goburrow/modbus and go.bug.st/serial libraries for the modbus interface and serial port interface accordingly
