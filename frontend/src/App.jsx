import {useEffect, useState} from 'react'

function App() {
  const [numbers,setNumbers] = useState(Array(8).fill(""));
  const [registers, setRegisters] = useState(Array(8).fill(0));
  const [error, setError] = useState(null);


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
          <button onClick={handleSubmit}>Send the Values</button>
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
        maxWidth: "800px",
        width: "100%",
    },
    title: {
        fontSize: "30px",
        marginBottom: "20px",
    },
    form: {
        display: "flex",
        flexWrap: "wrap",
        justifyContent: "space-between",
        gap: "10px", // Space between inputs
        width: "100%",
    },
    inputRow: {
        display: "flex",
        alignItems: "center",
        width: "calc(50%-10px)",
        minWidth: "300 px",
    },
    label: {
        fontSize: "16px",
        flex: "1",
        textAlign: "right",
        marginRight: "10px",
        marginLeft: "20px",
        whiteSpace: "nowrap",
    },
    input: {
        flex: "1",
        padding: "8px",
        fontSize: "16px",
    },
    button: {
        marginTop: "20px",
        padding: "10px 20px",
        fontSize: "16px",
        cursor: "pointer",
    },
};

export default App
