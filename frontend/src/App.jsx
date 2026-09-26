import { useEffect, useState } from "react";
import axios from "axios";

function App() {
  const [question, setQuestion] = useState("");
  const [option1, setOption1] = useState("");
  const [option2, setOption2] = useState("");
  const [polls, setPolls] = useState([]);

  const fetchPolls = async () => {
    try {
      const res = await axios.get("http://localhost:8080/polls");
      setPolls(res.data || []);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchPolls();
  }, []);

  const createPoll = async () => {
    if (!question || !option1 || !option2) {
      alert("Please fill all fields");
      return;
    }

    const pollData = {
      id: Date.now().toString(),
      question: question,
      options: [option1, option2],
    };

    try {
      await axios.post(
        "http://localhost:8080/create-poll",
        pollData
      );

      alert("Poll Created Successfully");

      setQuestion("");
      setOption1("");
      setOption2("");

      fetchPolls();
    } catch (err) {
      console.error(err);
      alert("Failed to create poll");
    }
  };

  return (
    <div style={{ padding: "20px", fontFamily: "Arial" }}>
      <h1>Live Poll App</h1>

      <h2>Create Poll</h2>

      <input
        type="text"
        placeholder="Poll Question"
        value={question}
        onChange={(e) => setQuestion(e.target.value)}
      />

      <br />
      <br />

      <input
        type="text"
        placeholder="Option 1"
        value={option1}
        onChange={(e) => setOption1(e.target.value)}
      />

      <br />
      <br />

      <input
        type="text"
        placeholder="Option 2"
        value={option2}
        onChange={(e) => setOption2(e.target.value)}
      />

      <br />
      <br />

      <button onClick={createPoll}>
        Create Poll
      </button>

      <hr />

      <h2>All Polls</h2>

      {polls.length === 0 ? (
        <p>No polls available</p>
      ) : (
        polls.map((poll) => (
          <div
            key={poll.id}
            style={{
              border: "1px solid #ccc",
              padding: "10px",
              marginBottom: "10px",
            }}
          >
            <h3>{poll.question}</h3>

            <ul>
              {poll.options.map((option, index) => (
                <li key={index}>{option}</li>
              ))}
            </ul>
          </div>
        ))
      )}
    </div>
  );
}

export default App;