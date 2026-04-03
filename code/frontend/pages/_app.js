import { useState, useEffect } from "react"
import "../styles/main.css"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export default function Home() {
  const [goals, setGoals] = useState([])
  const [title, setTitle] = useState("")
  const [loading, setLoading] = useState(false)

  const fetchGoals = async () => {
    try {
      const res = await fetch(`${API_BASE}/goals`)
      const data = await res.json()
      setGoals(data)
    } catch (error) {
      console.error("Error fetching goals:", error)
    }
  }

  const addGoal = async () => {
    if (!title.trim()) {
      alert("Please enter a goal title")
      return
    }
    
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/goals`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title })
      })
      
      const responseText = await res.text()
      console.log("Response status:", res.status)
      console.log("Response body:", responseText)
      
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${responseText}`)
      
      setTitle("")
      await fetchGoals()
    } catch (error) {
      console.error("Error adding goal:", error)
      alert("Failed to add goal: " + error.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchGoals()
  }, [])

  return (
    <div className="eoffice-page">
      <header className="eo-header">
        <div className="eo-logo-title">
          <div className="eo-logo-circle">eO</div>
          <div>
            <h1 className="eo-title">eOffice Performance Intelligence</h1>
            <p className="eo-subtitle">Real-time insights for file movement, pendency and productivity</p>
          </div>
        </div>
        <div className="eo-header-actions">
          <select className="eo-select">
            <option>All Departments</option>
            <option>Finance</option>
            <option>HR</option>
            <option>IT</option>
            <option>Administration</option>
          </select>
          <select className="eo-select">
            <option>Last 7 days</option>
            <option>Today</option>
            <option>Last 30 days</option>
            <option>FY 2025-26</option>
          </select>
          <button className="eo-btn primary">Export Report</button>
        </div>
      </header>

      <main className="eo-main">
        <div className="eo-goal-input-row">
          <input
            className="eo-input"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && addGoal()}
            placeholder="Enter new productivity goal"
            disabled={loading}
          />
          <button 
            className="eo-btn secondary" 
            onClick={addGoal}
            disabled={loading}
          >
            {loading ? "Adding..." : "Add Goal"}
          </button>
        </div>

        <section className="eo-kpi-cards">
          <div className="eo-card kpi-card">
            <div className="kpi-header">
              <span>Average Processing Time</span>
              <span className="kpi-pill kpi-pill-amber">↔ Stable</span>
            </div>
            <div className="kpi-value">2.6 days</div>
            <div className="kpi-footer">SLA breach: 7.4%</div>
          </div>

          <div className="eo-card kpi-card">
            <div className="kpi-header">
              <span>Pending Files</span>
              <span className="kpi-pill kpi-pill-red">145 high priority</span>
            </div>
            <div className="kpi-value">732</div>
            <div className="kpi-footer">Oldest pending: 34 days</div>
          </div>

          <div className="eo-card kpi-card">
            <div className="kpi-header">
              <span>On-time Disposal Rate</span>
              <span className="kpi-pill kpi-pill-green">Target: 90%</span>
            </div>
            <div className="kpi-value">86.3%</div>
            <div className="kpi-footer">Improved by 4.2% this month</div>
          </div>
        </section>

        <section className="eo-layout-2col">
          <div className="eo-card eo-chart-card">
            <div className="eo-card-header">
              <h2>Daily File Inflow vs Disposal</h2>
              <span className="eo-chip">Last 14 days</span>
            </div>
            <div className="eo-chart-placeholder">
              <p>Chart placeholder (Inflow vs Disposal)</p>
            </div>
          </div>

          <div className="eo-card eo-chart-card">
            <div className="eo-card-header">
              <h2>Department-wise Pendency</h2>
              <span className="eo-chip warning">Action required</span>
            </div>
            <div className="eo-chart-placeholder">
              <ul className="eo-pendency-list">
                <li>
                  <span>Finance</span>
                  <span className="eo-pill pill-red">214</span>
                </li>
                <li>
                  <span>HR</span>
                  <span className="eo-pill pill-amber">126</span>
                </li>
                <li>
                  <span>IT</span>
                  <span className="eo-pill pill-green">58</span>
                </li>
                <li>
                  <span>Administration</span>
                  <span className="eo-pill pill-amber">174</span>
                </li>
              </ul>
            </div>
          </div>
        </section>

        <section className="eo-layout-2col eo-bottom-grid">
          <div className="eo-card">
            <div className="eo-card-header">
              <h2>Team Goals</h2>
              <span className="eo-chip">Productivity Management</span>
            </div>

            <div className="eo-goal-input-row">
              <input
                className="eo-input"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                onKeyPress={(e) => e.key === "Enter" && addGoal()}
                placeholder="Enter new productivity goal"
                disabled={loading}
              />
              <button className="eo-btn secondary" onClick={addGoal} disabled={loading}>
                {loading ? "Adding..." : "Add Goal"}
              </button>
            </div>

            <ul className="eo-goal-list">
              {goals.map((goal) => (
                <li key={goal.id} className="eo-goal-item">
                  <span className="eo-goal-bullet" />
                  <span>{goal.title}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="eo-card eo-table-card">
            <div className="eo-card-header">
              <h2>Critical Pending Files</h2>
              <button className="eo-link-btn">View all</button>
            </div>
            <div className="eo-table-wrapper">
              <table className="eo-table">
                <thead>
                  <tr>
                    <th>File No.</th>
                    <th>Subject</th>
                    <th>Department</th>
                    <th>Pending With</th>
                    <th>Age (days)</th>
                    <th>Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>FIN/2025/1842</td>
                    <td>Budget allocation for Q4</td>
                    <td>Finance</td>
                    <td>Under Secretary</td>
                    <td>21</td>
                    <td><span className="eo-pill pill-red">High</span></td>
                  </tr>
                  <tr>
                    <td>HR/2025/903</td>
                    <td>Promotion review batch-3</td>
                    <td>HR</td>
                    <td>Deputy Director</td>
                    <td>17</td>
                    <td><span className="eo-pill pill-amber">Medium</span></td>
                  </tr>
                  <tr>
                    <td>IT/2025/771</td>
                    <td>Data center AMC renewal</td>
                    <td>IT</td>
                    <td>Section Officer</td>
                    <td>9</td>
                    <td><span className="eo-pill pill-green">Normal</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
