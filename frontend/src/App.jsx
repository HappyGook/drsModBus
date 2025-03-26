import {useEffect, useState} from 'react'


const ModBusUI= () =>{
    const [numbers,setNumbers] = useState(Array(8).fill(""));
    const [registers, setRegisters] = useState([]);
    const [error, setError] = useState("");
    const [ports,setPorts]=useState([])
    const [selectedPort, setSelectedPort]=useState("")
    const [confirmed, setConfirmed]=useState(null)

    useEffect(()=>{
        fetch("/api/list")
            .then(response=>response.json())
            .then((data)=>{
                if(data.ports&&data.ports.length>0){
                    setPorts(data.ports)
                } else{
                    setPorts([]) //Just the empty list
                    setError("No serial ports found")
                }
            }).catch((err)=>{
            console.error("Error fetching ports",err)
            setError("Could not fetch serial ports")
            setPorts([]) //Should clear it back to nothing if error occurs
        })

    },[])

    // Fetch register values only if port is confirmed
    useEffect(() => {
        if (confirmed) {
            fetch(`/api/read?port=${confirmed}`)
                .then((response) => response.json())
                .then((data) => {
                    if (data.registers) {
                        setRegisters(data.registers);
                    } else {
                        setError("Failed to load the register values");
                    }
                })
                .catch((err) => {
                    console.error("Error fetching registers", err);
                    setError("Could not connect to backend");
                });
        }
    }, [confirmed]);

    //A change of port the user selects
    const handlePortChange = (event) => {
        setSelectedPort(event.target.value);
    };

    // Confirm selected port
    const handleConfirm = () => {
        if (selectedPort) {
            setConfirmed(selectedPort);
            setError(null); // Clear errors
        }
    };

    //A change of number in register input field
    const handleChange=(index, value)=>{
        const newNumbers=[...numbers];
        newNumbers[index]=value;
        setNumbers(newNumbers);
    };

    //Values sent to backend
    const handleSubmit = () => {
        if (!confirmed) return;
        fetch(`/api/submit?port=${confirmed}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ port: selectedPort, values: numbers.map(Number) }),
        })
            .then((res) => res.json())
            .then((data) => console.log("Sent successfully:", data))
            .catch((err) => console.error("Error sending data:", err));
    };

    return (
        <div style={styles.container}>
            <h1 style={styles.title}>ModBus</h1>

            {/* Serial Port Selection */}
            <div style={styles.portSelection}>
                <select style={styles.select} value={selectedPort} onChange={handlePortChange}>
                    <option value="">Select a Serial Port</option>
                    {ports.map((port, index) => (
                        <option key={index} value={port}>{port}</option>
                    ))}
                </select>
                <button style={styles.button} onClick={handleConfirm}>Confirm Port</button>
            </div>

            {error && <p style={styles.error}>{error}</p>}


            {/* Only show input fields if port is confirmed */}
            {confirmed && (
                <div style={styles.form}>
                    {labels.map((label, index) => (
                        <div key={index} style={styles.inputRow}>
                            <label style={styles.label}>{label} </label>
                            <input
                                type="number"
                                value={numbers[index]}
                                onChange={(e) => handleChange(index, e.target.value)}
                                placeholder={registers[index]}
                                style={styles.input}
                            />
                        </div>
                    ))}
                </div>
            )}

            {confirmed && <button style={styles.button} onClick={handleSubmit}>Send the Values</button>}
        </div>
    );
};



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
    portSelection:{
      display:"flex",
      alignItems:"center",
      gap:"10px",
      marginBottom:"20px",
    },
    confirmButton: {
        padding: "10px 20px",
        fontSize: "16px",
        backgroundColor: "#28a745",
        color: "#fff",
        border: "none",
        cursor: "pointer",
        borderRadius: "5px",
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
        flex: "1",
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
    error:{
        color:"red",
        marginTop:"10px",
    },
};

export default ModBusUI;
