import { useState, useEffect } from "react"

export default function Home() {
  const [goals, setGoals] = useState([])
  const [title, setTitle] = useState("")

  const fetchGoals = async () => {
    const res = await fetch("http://localhost:8080/goals")
    const data = await res.json()
    setGoals(data)
  }

  const addGoal = async () => {
    await fetch("http://localhost:8080/goals", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title })
    })
    setTitle("")
    fetchGoals()
  }

  useEffect(() => {
    fetchGoals()
  }, [])

  return (
    <div style={{ padding: 40 }}>
      <h1>Productivity Management System</h1>

      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Enter Goal"
      />
      <button onClick={addGoal}>Add Goal</button>

      <ul>
        {goals.map((goal) => (
          <li key={goal.id}>{goal.title}</li>
        ))}
      </ul>
    </div>
  )
}