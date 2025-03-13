import {use, useEffect, useState} from 'react'

function App() {
  const [numbers,setNumbers] = useState(Array(8).fill(""));
  const [registers, setRegisters] = useState(Array(8).fill(0));
  const [setError] = useState(null);
  const [ports,setPorts]=useState([])
    const [selectedPort, setSelectedPort]=useState("")


  useEffect(()=>{
      fetch("/api/read")
          .then((response)=>response.json())
          .then((data)=>{
              if(data.registers) {
                  setRegisters(data.registers)
              }else{
                  setError("Failed to load the register values")
              }
          }).catch((err)=>{
              console.error("Error fetching registers",err)
              setError("Could not connect to backend")
      })
      fetch("/api/list")
          .then(response=>response.json())
          .then((data)=>{
              if(data.ports){
                  setPorts(data.ports)
              } else{
                  setError("No serial ports found")
              }
          }).catch((err)=>{
              console.error("Error fetching ports",err)
                setError("Could not fetch serial ports")
      })

  },[])


    const handleChange=(index, value)=>{
      const newNumbers=[...numbers];
        newNumbers[index]=value;
        setNumbers(newNumbers);
    };

    const handleSubmit=async()=>{
      const response = await fetch("/api/submit",{
          method:"POST",
          headers: {"Content-Type":"application/json"},
          body: JSON.stringify({numbers}),
      });

      const data = await response.json();
      console.log("Backend's response: ", data);
    };

  return (
      <div style={styles.container}>
          <h1 style={styles.title}>ModBus</h1>

          <div style={styles.form}>
              <label style={styles.label}>Select Serial Port</label>
              <select
                  style={{ ...styles.selectContainer, cursor: "pointer"}}
                  value={selectedPort}
                  onChange={event => setSelectedPort(event.target.value)}
                  >
                  <option value="">Select a Port...</option>
                  {ports.map((port, index)=>(
                      <option key={index} value={port}>
                          {port}
                      </option>
                  ))}
              </select>
          </div>

          <div style={styles.form}>
              {labels.map((label,index)=>(
                  <div key={index} style={styles.inputRow}>
                      <label style={styles.label}>{label} </label>
                      <input
                          type="number"
                          value={numbers[index]}
                          onChange={(e)=> handleChange(index,e.target.value)}
                          placeholder={registers[index]}
                          style={styles.input}
                      />
                  </div>
              ))}
          </div>
          <button style={styles.button} onClick={handleSubmit}>Send the Values</button>
      </div>
  );
}
const labels=[
    "Output Voltage set",
    "Constant Voltage setting",
    "Floating Voltage setting",
    "CC Charge Timeout setting",
    "CV Charge Timeout setting",
    "FV Charge Timeout setting",
    "BAT_LOW Protect setting",
    "Force BAT_LOW protect setting",
]

const styles = {
    container: {
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        marginTop: "50px",
        maxWidth: "600px",
        width: "100%",
        padding: "20px",
        borderRadius:"8px",
        backgroundColor:"#f9f9f9",
        boxShadow:"0 4px 8px rgba(0, 0, 0, 0.1)",
    },
    title: {
        fontSize: "32px",
        fontWeight: "bold",
        marginBottom: "25px",
        color: "#333",
    },
    form: {
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        width: "100%",
        gap: "15px",
    },
    selectContainer: {
        width: "100%",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        marginBottom: "20px",
    },
    inputRow: {
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        width: "100%",
        maxWidth: "500px",
        gap: "10px",
    },
    label: {
        fontSize: "16px",
        fontWeight: "bold",
        flex: "1",
        textAlign: "right",
        marginRight: "10px",
        whiteSpace: "nowrap",
    },
    input: {
        flex: "2",
        padding: "10px",
        fontSize: "16px",
        border: "1px solid #ccc",
        borderRadius: "5px",
    },
    select: {
        width: "100%",
        maxWidth: "500px",
        padding: "10px",
        fontSize: "16px",
        border: "1px solid #ccc",
        borderRadius: "5px",
        cursor: "pointer",
    },
    button: {
        marginTop: "25px",
        padding: "12px 24px",
        fontSize: "18px",
        color: "#fff",
        backgroundColor: "#007bff",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
        transition: "background 0.3s",
    },
    buttonHover: {
        backgroundColor: "#0056b3",
    },
};

export default App
