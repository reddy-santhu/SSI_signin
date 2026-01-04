import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import './Dashboard.css';

function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('sessionToken');
    if (!token) {
      navigate('/');
      return;
    }

    const fetchUserData = async () => {
      try {
        const response = await api.get('/dashboard');
        setUser(response.data.user);
      } catch (err) {
        console.error('Failed to fetch user data:', err);
        if (err.response?.status === 401) {
          localStorage.removeItem('sessionToken');
          navigate('/');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('sessionToken');
    navigate('/');
  };

  if (loading) {
    return <div className="dashboard-loading">Loading...</div>;
  }

  return (
    <div className="dashboard-page">
      <div className="dashboard-container">
        <div className="card">
          <div className="dashboard-header">
            <h1>Dashboard</h1>
            <button className="btn btn-secondary" onClick={handleLogout}>
              Logout
            </button>
          </div>

          {user && (
            <div className="user-info">
              <h2>User Information</h2>
              <div className="info-item">
                <strong>DID:</strong> {user.did || 'Not set'}
              </div>
              <div className="info-item">
                <strong>Phone:</strong> {user.phone || 'Not set'}
              </div>
              <div className="info-item">
                <strong>Member Since:</strong> {new Date(user.created_at).toLocaleDateString()}
              </div>
            </div>
          )}

          <div className="todo-section">
            <h3>Coming Soon</h3>
            <ul>
              <li>User profile management</li>
              <li>Credential viewing</li>
              <li>Settings</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;

